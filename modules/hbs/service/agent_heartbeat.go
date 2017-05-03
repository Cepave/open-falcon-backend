package service

import (
	"net/http"

	osling "github.com/Cepave/open-falcon-backend/common/sling"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/dghubble/sling"
)

type AgentHeartbeatService struct {
	started         bool
	slingInit       *sling.Sling
	rowsAffectedCnt int64
}

func NewAgentHeartbeatService(httpClient *http.Client) *AgentHeartbeatService {
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
}

func (s *AgentHeartbeatService) Stop() {
	if !s.started {
		return
	}
	/*
	 * ToDo
	 * Stop queue
	 * Flush the data
	 */
}

func (s *AgentHeartbeatService) CurrentSize() int {
	/*
	 * ToDo
	 * Return the size of queue
	 */
	return 0
}

func (s *AgentHeartbeatService) CumulativeRowsAffected() int64 {
	return s.rowsAffectedCnt
}

func (s *AgentHeartbeatService) Heartbeat(agents []*model.AgentHeartbeat) {
	param := struct {
		UpdateOnly bool `json:"update_only"`
	}{}
	if g.Config().Hosts != "" {
		param.UpdateOnly = true
	}
	s.slingInit = s.slingInit.BodyJSON(agents).QueryStruct(&param)

	res := model.AgentHeartbeatResult{}
	err := osling.ToSlintExt(s.slingInit).DoReceive(http.StatusOK, &res)
	if err != nil {
		logger.Errorln("Heartbeat:", err)
	}

	s.rowsAffectedCnt += res.RowsAffected
}
