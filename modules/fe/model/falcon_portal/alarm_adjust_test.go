package falconPortal

import (
	"log"
	"testing"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func RecoveryALL(q orm.Ormer) {
	q.Raw("UPDATE event_cases SET status = 'PROBLEM' WHERE status != 'OK'").Exec()
}

func TestPortalEventCase(t *testing.T) {
	g.ParseConfig("../../cfg.json")
	config := g.Config()
	log.Printf("%v", config.Uic.Addr)
	orm.RegisterDataBase("default", "mysql", config.Uic.Addr, config.Uic.Idle, config.Uic.Max)
	orm.RegisterDataBase("graph", "mysql", config.GraphDB.Addr, config.GraphDB.Idle, config.GraphDB.Max)
	orm.RegisterDataBase("falcon_portal", "mysql", config.FalconPortal.Addr, config.FalconPortal.Idle, config.FalconPortal.Max)
	q := orm.NewOrm()
	q.Using("falcon_portal")
	Convey("test - WhenStrategyUpdated", t, func() {
		err, affectedRows := WhenStrategyUpdated(41)
		So(err, ShouldEqual, nil)
		So(affectedRows, ShouldEqual, 1)
		RecoveryALL(q)
	})
	Convey("test - WhenStrategyDeleted", t, func() {
		err, affectedRows := WhenStrategyDeleted(41)
		So(err, ShouldEqual, nil)
		So(affectedRows, ShouldEqual, 1)
		RecoveryALL(q)
	})
	Convey("test - WhenTempleteDeleted", t, func() {
		err, affectedRows := WhenTempleteDeleted(3)
		So(err, ShouldEqual, nil)
		So(affectedRows, ShouldEqual, 78)
		RecoveryALL(q)
	})
	Convey("test - WhenTempleteUnbind", t, func() {
		err, affectedRows := WhenTempleteUnbind(7, 6)
		So(err, ShouldEqual, nil)
		So(affectedRows, ShouldEqual, 1)
		RecoveryALL(q)
	})
	Convey("test - WhenEndpointUnbind", t, func() {
		err, affectedRows := WhenEndpointUnbind(1, 6)
		So(err, ShouldEqual, nil)
		So(affectedRows, ShouldEqual, 1)
		RecoveryALL(q)
	})
}
