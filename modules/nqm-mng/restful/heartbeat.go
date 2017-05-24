package restful

import (
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	"github.com/gin-gonic/gin"
)

type heartbeatOfAgents []*model.AgentHeartbeat

func (agents *heartbeatOfAgents) Bind(context *gin.Context) {
	ogin.BindJson(context, agents)
}

func agentHeartbeat(
	agents *heartbeatOfAgents,
	q *struct {
		UpdateOnly bool `mvc:"query[update_only]"`
	},
) mvc.OutputBody {
	retBody := rdb.AgentHeartbeat(*agents, q.UpdateOnly)
	return mvc.JsonOutputBody(retBody)
}
