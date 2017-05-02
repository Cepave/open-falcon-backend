package service

import (
	"net/http"

	"fmt"

	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/dghubble/sling"
)

type AgentHeartbeatService struct {
	slingInit *sling.Sling
}

func NewAgentHeartbeatService(httpClient *http.Client) *AgentHeartbeatService {
	/*
	 * ToDo
	 * HttpConfig
	 */
	return &AgentHeartbeatService{
		slingInit: sling.New().Client(httpClient).Base("ToDoBase").Post("api/v1/agent/heartbeat"),
	}
}

func (s *AgentHeartbeatService) Start() {
	/*
	 * ToDo
	 * Initial & start queue
	 */
}

func (s *AgentHeartbeatService) Stop() {
	/*
	 * ToDo
	 * Stop queue
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
	/*
	 * ToDo
	 * Return the cumulative value of rows affected
	 */
	return 0
}

func (s *AgentHeartbeatService) Heartbeat(agents []*model.AgentHeartbeat) {
	param := struct {
		UpdateOnly bool `json:"update_only"`
	}{}
	if g.Config().Hosts != "" {
		param.UpdateOnly = true
	}
	res := model.AgentHeartbeatResult{}
	s.slingInit = s.slingInit.BodyJSON(agents).QueryStruct(&param)
	resp, err := s.slingInit.ReceiveSuccess(&res)
	fmt.Printf("Resp: %v, Err: %v\n", resp, err)
}
