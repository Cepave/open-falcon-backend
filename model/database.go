package model

import (
	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/model/dashboard"
	"github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Uic.Addr, config.Uic.Idle, config.Uic.Max)
	orm.RegisterDataBase("graph", "mysql", config.Graph.Addr, config.Graph.Idle, config.Graph.Max)

	// register model
	orm.RegisterModel(new(uic.User), new(uic.Team), new(uic.Session), new(uic.RelTeamUser), new(dashboard.Endpoint), new(dashboard.EndpointCounter))

	if config.Log == "debug" {
		orm.Debug = true
	}
}
