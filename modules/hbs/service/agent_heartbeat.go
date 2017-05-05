package service

import (
	"net/http"
	"time"

	osling "github.com/Cepave/open-falcon-backend/common/sling"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/dghubble/sling"
)

var agentHeartbeatService *AgentHeartbeatService

type AgentHeartbeatService struct {
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
	/*
	 * ToDo
	 * Initial & start queue
	 */
	s.started = true

	go func() {
		for {
			if s.CurrentSize() > 0 {
				/*
				 * ToDo
				 * Get agents from queue
				 */
				toDoGetAgents := make([]*model.AgentHeartbeat, 10)
				s.Heartbeat(toDoGetAgents)
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}()
}

func (s *AgentHeartbeatService) Stop() {
	if !s.started {
		return
	}

	/*
	 * ToDo
	 * Close/Stop queue
	 */
	s.started = false
	queueSize := s.CurrentSize()
	logger.Infof("Stopping AgentHeartbeatService. Size of queue: [%d]", queueSize)

	/**
	 * Waiting for queue to be processed
	 */
	maxTimes := 15
	for queueSize > 0 && maxTimes > 0 {
		time.Sleep(2 * time.Second)
		maxTimes--
		queueSize = s.CurrentSize()
		logger.Infof("Sleep for 2 seconds to wait queue to be processed... Current size: [%d]", queueSize)
	}
}

func (s *AgentHeartbeatService) CurrentSize() int {
	/*
	 * ToDo
	 * Return the size of queue
	 */
	return 0
}

func (s *AgentHeartbeatService) CumulativeAgentsDropped() int64 {
	return s.agentsDroppedCnt
}

func (s *AgentHeartbeatService) CumulativeRowsAffected() int64 {
	return s.rowsAffectedCnt
}

func (s *AgentHeartbeatService) Heartbeat(agents []*model.AgentHeartbeat) {
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
