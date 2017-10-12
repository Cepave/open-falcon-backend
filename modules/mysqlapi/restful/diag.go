package restful

import (
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	apiModel "github.com/Cepave/open-falcon-backend/common/model/mysqlapi"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/service"
)

func health() mvc.OutputBody {
	diagRdb := rdb.DiagnoseRdb(
		rdb.DbConfig.Dsn,
		rdb.DbFacade.SqlDb,
	)
	resp := &apiModel.HealthView{
		Rdb: diagRdb,
		Http: &apiModel.Http{
			Listening: GinConfig.GetAddress(),
		},
		Nqm: &apiModel.Nqm{
			Heartbeat: &apiModel.Heartbeat{
				Count: service.NqmQueue.ConsumedCount(),
			},
		},
	}

	return mvc.JsonOutputBody(resp)
}
