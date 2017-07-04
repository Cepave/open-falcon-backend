package service

import (
	"reflect"
	"sync"
	"sync/atomic"

	cModel "github.com/Cepave/open-falcon-backend/common/model"
	cQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/hbs/cache"
)

var (
	elementType = reflect.TypeOf(new(cModel.FalconAgentHeartbeat))
)

type AgentHeartbeatService struct {
	wg               *sync.WaitGroup
	safeQ            *cQueue.Queue
	qConfig          *cQueue.Config
	running          bool
	agentsPutCnt     int64
	heartbeatCall    func([]*cModel.FalconAgentHeartbeat) (int64, int64)
	rowsAffectedCnt  int64
	agentsDroppedCnt int64
}

func NewAgentHeartbeatService(config *cQueue.Config) *AgentHeartbeatService {
	return &AgentHeartbeatService{
		wg:            &sync.WaitGroup{},
		safeQ:         cQueue.New(),
		qConfig:       config,
		heartbeatCall: agentHeartbeatCall,
	}
}

func (s *AgentHeartbeatService) Start() {
	if s.running {
		logger.Infoln("[AgentHeartbeat][Skipped] Service is already running.")
		return
	}
	s.running = true
	logger.Infoln("[AgentHeartbeat] Service is starting.")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			if !s.running {
				break
			}

			s.consumeHeartbeatQueue(false)
		}

		s.consumeHeartbeatQueue(true)
	}()
}

func (s *AgentHeartbeatService) consumeHeartbeatQueue(flushing bool) {

	agents := s.safeQ.DrainNWithDurationByReflectType(s.qConfig, elementType).([]*cModel.FalconAgentHeartbeat)

	if len(agents) == 0 {
		return
	}

	r, d := s.heartbeatCall(agents)
	s.rowsAffectedCnt += r
	s.agentsDroppedCnt += d

	if flushing {
		logger.Infof("[AgentHeartbeat] Service is flushing. Number of agents: %d ", len(agents))
		s.consumeHeartbeatQueue(flushing)
	}
}

func (s *AgentHeartbeatService) Stop() {
	if !s.running {
		logger.Infoln("[AgentHeartbeat][Skipped] Service is already stopped.")
		return
	}

	s.running = false
	logger.Infof("[AgentHeartbeat] Service is stopping. Size of queue: %d", s.CurrentSize())

	/**
	 * Waiting for queue to be processed
	 */
	s.wg.Wait()
}

func (s *AgentHeartbeatService) Put(req *cModel.AgentReportRequest, updateTime int64) {
	if !s.running {
		logger.Infoln("[AgentHeartbeat][Skipped] Put when stopped.")
		return
	}

	cache.Agents.Put(req, updateTime)
	hb := requestToHeartbeat(req, updateTime)
	s.safeQ.Enqueue(hb)
	atomic.AddInt64(&(s.agentsPutCnt), 1)
}

func (s *AgentHeartbeatService) CurrentSize() int {
	return s.safeQ.Len()
}

func (s *AgentHeartbeatService) CumulativeAgentsDropped() int64 {
	return s.agentsDroppedCnt
}

func (s *AgentHeartbeatService) CumulativeAgentsPut() int64 {
	return s.agentsPutCnt
}

func (s *AgentHeartbeatService) CumulativeRowsAffected() int64 {
	return s.rowsAffectedCnt
}

func requestToHeartbeat(req *cModel.AgentReportRequest, updateTime int64) *cModel.FalconAgentHeartbeat {
	return &cModel.FalconAgentHeartbeat{
		Hostname:      req.Hostname,
		IP:            req.IP,
		AgentVersion:  req.AgentVersion,
		PluginVersion: req.PluginVersion,
		UpdateTime:    updateTime,
	}
}
