package restful

import (
	"net"

	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/service"
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
	var r *nqmModel.Agent
	if rdb.NotNewNqmAgent(req.ConnectionId) {
		service.NqmQueue.Put(req)
		r = overwrittenNqmAgent(req)
	} else {
		r = rdb.InsertNqmAgentByHeartbeat(req)
	}
	return mvc.JsonOutputBody(r)
}

// overwrittenNqmAgent overwrites the result with the values from the heartbeat
// request. The values in the database should be identical in the end.
func overwrittenNqmAgent(req *model.NqmAgentHeartbeatRequest) *nqmModel.Agent {
	r := rdb.SelectNqmAgentByConnId(req.ConnectionId)
	r.ConnectionId = req.ConnectionId
	r.Hostname = req.Hostname
	r.IpAddress = net.ParseIP(req.IpAddress.String())
	r.LastHeartBeat = req.Timestamp
	return r
}
