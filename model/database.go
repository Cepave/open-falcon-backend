package model

import (
	"github.com/Cepave/fe/g"
	. "github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Uic.Addr, config.Uic.Idle, config.Uic.Max)

	// register model
	orm.RegisterModel(new(User), new(Team), new(Session), new(RelTeamUser))

	if config.Log == "debug" {
		orm.Debug = true
	}
}
