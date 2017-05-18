package service

import (
	"net/http"
	"sync"
	"time"

	oqueue "github.com/Cepave/open-falcon-backend/common/queue"
	osling "github.com/Cepave/open-falcon-backend/common/sling"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/dghubble/sling"
)

var agentHeartbeatService *AgentHeartbeatService

type AgentHeartbeatService struct {
	sync.WaitGroup
	safeQ            *oqueue.Queue
	started          bool
	slingInit        *sling.Sling
	rowsAffectedCnt  int64
	agentsDroppedCnt int64
}

func NewAgentHeartbeatService() *AgentHeartbeatService {
	return &AgentHeartbeatService{
		slingInit: NewSlingBase().Post("api/v1/agent/heartbeat"),
	}
}

func (s *AgentHeartbeatService) Start() {
	if s.started {
		return
	}
	s.started = true
	s.safeQ = oqueue.New()

	s.Add(1)
	go func() {
		defer s.Done()

		for {
			if !s.started {
				break
			}

			s.consumeHeartbeatQueue(100*time.Millisecond, false)

			if !s.started {
				break
			}

			time.Sleep(5 * time.Second)
		}

		s.consumeHeartbeatQueue(0, true)
	}()
}

func (s *AgentHeartbeatService) consumeHeartbeatQueue(waitForQueue time.Duration, logFlag bool) {
	for {
		/*
		 * ToDo
		 * Configuration
		 * ToReview
		 */
		c := oqueue.Config{}
		var elementType *model.AgentHeartbeat
		absArray := s.safeQ.DrainNWithDurationByType(&c, elementType)
		agents := absArray.([]*model.AgentHeartbeat)
		agentsNum := len(agents)
		if agentsNum == 0 {
			break
		}

		s.heartbeat(agents)
		if logFlag {
			logger.Infof("Flushing [%d] agents", agentsNum)
		}
	}
}

func (s *AgentHeartbeatService) Stop() {
	if !s.started {
		return
	}

	s.started = false
	logger.Infof("Stopping AgentHeartbeatService. Size of queue: [%d]", s.CurrentSize())

	/**
	 * Waiting for queue to be processed
	 */
	s.Wait()
	s.safeQ = nil
}

func (s *AgentHeartbeatService) Put(agent *model.AgentHeartbeat) {
	if !s.started {
		return
	}
	s.safeQ.Enqueue(agent)
}

func (s *AgentHeartbeatService) CurrentSize() int {
	return s.safeQ.Len()
}

func (s *AgentHeartbeatService) CumulativeAgentsDropped() int64 {
	return s.agentsDroppedCnt
}

func (s *AgentHeartbeatService) CumulativeRowsAffected() int64 {
	return s.rowsAffectedCnt
}

func (s *AgentHeartbeatService) heartbeat(agents []*model.AgentHeartbeat) {
	param := struct {
		UpdateOnly bool `json:"update_only"`
	}{updateOnlyFlag}
	s.slingInit = s.slingInit.BodyJSON(agents).QueryStruct(&param)

	res := model.AgentHeartbeatResult{}
	err := osling.ToSlintExt(s.slingInit).DoReceive(http.StatusOK, &res)
	if err != nil {
		s.agentsDroppedCnt += int64(len(agents))
		logger.Errorln("Heartbeat:", err)
		return
	}

	s.rowsAffectedCnt += res.RowsAffected
}
