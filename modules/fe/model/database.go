package model

import (
	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/dashboard"
	event "github.com/Cepave/open-falcon-backend/modules/fe/model/falcon_portal"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {

	// register model
	orm.RegisterModel(new(uic.User), new(uic.Team), new(uic.Session),
		new(uic.RelTeamUser), new(dashboard.Endpoint), new(dashboard.EndpointCounter),
		new(dashboard.HostGroup), new(dashboard.Hosts), new(event.EventCases),
		new(event.Events), new(event.Tpl), new(event.EventNote))
	// databases_name, sql_derver, db_addr, idle_time, mix_connection
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Uic.Addr, config.Uic.Idle, config.Uic.Max)
	orm.RegisterDataBase("graph", "mysql", config.GraphDB.Addr, config.GraphDB.Idle, config.GraphDB.Max)
	orm.RegisterDataBase("falcon_portal", "mysql", config.FalconPortal.Addr, config.FalconPortal.Idle, config.FalconPortal.Max)
	orm.RegisterDataBase("boss", "mysql", config.BossDB.Addr, config.BossDB.Idle, config.BossDB.Max)

	if config.Log == "debug" {
		orm.Debug = true
	}
}
