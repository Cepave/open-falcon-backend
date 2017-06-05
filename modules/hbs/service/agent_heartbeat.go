package service

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	commonSling "github.com/Cepave/open-falcon-backend/common/sling"
	"github.com/Cepave/open-falcon-backend/modules/hbs/cache"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
)

var (
	timeWaitForQueue = 100 * time.Millisecond
	timeWaitForInput = 5 * time.Second
)

type AgentHeartbeatService struct {
	wg               *sync.WaitGroup
	safeQ            *commonQueue.Queue
	qConfig          *commonQueue.Config
	started          bool
	agentsPutCnt     int64
	heartbeatCall    func([]*model.AgentHeartbeat) (int64, int64)
	rowsAffectedCnt  int64
	agentsDroppedCnt int64
}

func NewAgentHeartbeatService(config *commonQueue.Config) *AgentHeartbeatService {
	return &AgentHeartbeatService{
		wg:            &sync.WaitGroup{},
		safeQ:         commonQueue.New(),
		qConfig:       config,
		heartbeatCall: buildHeartbeatCall(),
	}
}

func (s *AgentHeartbeatService) Start() {
	if s.started {
		logger.Infoln("[AgentHeartbeat][Skipped] Service is already started.")
		return
	}
	s.started = true
	logger.Infoln("[AgentHeartbeat] Service is starting.")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			if !s.started {
				break
			}

			s.consumeHeartbeatQueue(timeWaitForQueue, false)

			if !s.started {
				break
			}

			time.Sleep(timeWaitForInput)
		}

		s.consumeHeartbeatQueue(0, true)
	}()
}

func (s *AgentHeartbeatService) consumeHeartbeatQueue(waitForQueue time.Duration, logFlag bool) {
	for {
		var elementType *model.AgentHeartbeat
		absArray := s.safeQ.DrainNWithDurationByType(s.qConfig, elementType)
		agents := absArray.([]*model.AgentHeartbeat)
		agentsNum := len(agents)
		if agentsNum == 0 {
			break
		}

		r, d := s.heartbeatCall(agents)
		s.rowsAffectedCnt += r
		s.agentsDroppedCnt += d
		if logFlag {
			logger.Infof("[AgentHeartbeat] Service is flushing. Number of agents: %d ", agentsNum)
		}
		time.Sleep(waitForQueue)
	}
}

func (s *AgentHeartbeatService) Stop() {
	if !s.started {
		logger.Infoln("[AgentHeartbeat][Skipped] Service is already stopped.")
		return
	}

	s.started = false
	logger.Infof("[AgentHeartbeat] Service is stopping. Size of queue: %d", s.CurrentSize())

	/**
	 * Waiting for queue to be processed
	 */
	s.wg.Wait()
}

func (s *AgentHeartbeatService) Put(req *commonModel.AgentReportRequest) {
	if !s.started {
		logger.Infoln("[AgentHeartbeat][Skipped] Put when stopped.")
		return
	}
	now := time.Now().Unix()
	cache.Agents.Put(req, now)
	agent := &model.AgentHeartbeat{
		Hostname:      req.Hostname,
		IP:            req.IP,
		AgentVersion:  req.AgentVersion,
		PluginVersion: req.PluginVersion,
		UpdateTime:    now,
	}
	s.safeQ.Enqueue(agent)
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

func buildHeartbeatCall() func([]*model.AgentHeartbeat) (int64, int64) {

	return func(agents []*model.AgentHeartbeat) (rowsAffectedCnt int64, agentsDroppedCnt int64) {
		param := struct {
			UpdateOnly bool `json:"update_only"`
		}{updateOnlyFlag}
		req := NewSlingBase().Post("api/v1/agent/heartbeat").BodyJSON(agents).QueryStruct(&param)

		res := model.AgentHeartbeatResult{}
		err := commonSling.ToSlintExt(req).DoReceive(http.StatusOK, &res)
		if err != nil {
			logger.Errorln("[AgentHeartbeat]", err)
			return 0, int64(len(agents))
		}

		return res.RowsAffected, 0
	}
}
