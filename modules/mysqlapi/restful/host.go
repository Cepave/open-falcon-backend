package restful

import (
	mvc "github.com/Cepave/open-falcon-backend/common/gin/mvc"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
)

func listHosts(
	paging *struct {
		Page *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[id#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	agents, resultPaging := rdb.ListHosts(*paging.Page)

	return resultPaging, mvc.JsonOutputBody(agents)
}
