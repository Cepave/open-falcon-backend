package model

import (
	"github.com/Cepave/alarm/g"
	"github.com/Cepave/alarm/model/event"
	"github.com/Cepave/alarm/model/uic"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Uic.Addr, config.Uic.Idle, config.Uic.Max)
	orm.RegisterDataBase("falcon_portal", "mysql", config.FalconPortal.Addr, config.FalconPortal.Idle, config.FalconPortal.Max)
	// register model
	orm.RegisterModel(new(uic.User), new(uic.Session), new(event.Events), new(event.EventCases))
	if config.Debug {
		orm.Debug = true
	}
}
