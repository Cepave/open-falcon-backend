package service

import (
	"net/http"

	commonSling "github.com/Cepave/open-falcon-backend/common/sling"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
)

func agentHeartbeatCall(agents []*model.FalconAgentHeartbeat) (rowsAffectedCnt int64, agentsDroppedCnt int64) {
	param := struct {
		UpdateOnly bool `json:"update_only"`
	}{updateOnlyFlag}
	req := NewSlingBase().Post("api/v1/agent/heartbeat").BodyJSON(agents).QueryStruct(&param)

	res := model.FalconAgentHeartbeatResult{}
	err := commonSling.ToSlintExt(req).DoReceive(http.StatusOK, &res)
	if err != nil {
		logger.Errorln("[AgentHeartbeat]", err)
		return 0, int64(len(agents))
	}

	return res.RowsAffected, 0
}
