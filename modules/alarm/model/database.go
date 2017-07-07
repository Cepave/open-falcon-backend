package model

import (
	"github.com/Cepave/open-falcon-backend/modules/alarm/g"
	"github.com/Cepave/open-falcon-backend/modules/alarm/model/event"
	"github.com/Cepave/open-falcon-backend/modules/alarm/model/uic"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Uic.Addr, config.Uic.Idle, config.Uic.Max)
	orm.RegisterDataBase("falcon_portal", "mysql", config.FalconPortal.Addr, config.FalconPortal.Idle, config.FalconPortal.Max)
	orm.RegisterDataBase("boss", "mysql", config.BossConfig.Addr, config.BossConfig.Idle, config.BossConfig.Max)
	// register model
	orm.RegisterModel(new(uic.User), new(uic.Session), new(event.Events), new(event.EventCases), new(event.AlarmType))
	if config.Debug {
		orm.Debug = true
	}
}
