package restful

import (
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/service"
)

func health() mvc.OutputBody {
	diagRdb := rdb.DiagnoseRdb(
		rdb.DbConfig.Dsn,
		rdb.DbFacade.SqlDb,
	)
	resp := &model.HealthView{
		Rdb: diagRdb,
		Http: &model.Http{
			Listening: GinConfig.GetAddress(),
		},
		Nqm: &model.Nqm{
			Heartbeat: &model.Heartbeat{
				Count: service.NqmQueue.ConsumedCount(),
			},
		},
	}

	return mvc.JsonOutputBody(resp)
}
