package restful

import (
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/service/queue"
	"github.com/gin-gonic/gin"
)

type heartbeatOfFalconAgents []*model.FalconAgentHeartbeat

func (agents *heartbeatOfFalconAgents) Bind(context *gin.Context) {
	ogin.BindJson(context, agents)
}

func falconAgentHeartbeat(
	agents *heartbeatOfFalconAgents,
	q *struct {
		UpdateOnly bool `mvc:"query[update_only]"`
	},
) mvc.OutputBody {
	retBody := rdb.FalconAgentHeartbeat(*agents, q.UpdateOnly)
	return mvc.JsonOutputBody(retBody)
}

func nqmAgentHeartbeat(
	req *model.NqmAgentHeartbeatRequest,
) mvc.OutputBody {
	if rdb.NotNewNqmAgent(req.ConnectionId) {
		queue.NqmQueue.Put(req)
	} else {
		rdb.InsertNqmAgent(req)
	}
	r := rdb.SelectNqmAgentByConnId(req.ConnectionId)
	return mvc.JsonOutputBody(r)
}
