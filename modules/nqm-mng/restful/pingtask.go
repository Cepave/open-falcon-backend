package restful

import (
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonNqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"

	mvc "github.com/Cepave/open-falcon-backend/common/gin/mvc"
)

func listAgentsByPingTask(
	query *commonNqmModel.AgentQueryWithPingTask,
	paging *struct {
		Page *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[status#desc:connection_id#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	agents, resultPaging := commonNqmDb.ListAgentsWithPingTask(query, *paging.Page)

	return resultPaging, mvc.JsonOutputBody(agents)
}
