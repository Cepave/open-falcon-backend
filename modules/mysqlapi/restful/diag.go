package restful

import (
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/service"
)

func health() mvc.OutputBody {
	portalRdbDiag := rdb.GlobalDbHolder.Diagnose(rdb.DB_PORTAL)
	graphRdbDiag := rdb.GlobalDbHolder.Diagnose(rdb.DB_GRAPH)

	health := &model.HealthView{
		Rdb: map[string]interface{}{
			"dsn":              portalRdbDiag.Dsn,
			"open_connections": portalRdbDiag.OpenConnections,
			"ping_result":      portalRdbDiag.PingResult,
			"ping_message":     portalRdbDiag.PingMessage,
			"portal":           portalRdbDiag,
			"graph":            graphRdbDiag,
		},
		Http: &model.Http{
			Listening: GinConfig.GetAddress(),
		},
		Nqm: &model.Nqm{
			Heartbeat: &model.Heartbeat{
				Count: service.NqmQueue.ConsumedCount(),
			},
		},
	}

	return mvc.JsonOutputBody(health)
}
