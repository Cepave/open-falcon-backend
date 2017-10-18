package test

import (
	"flag"
	"testing"

	coommonModel "github.com/Cepave/open-falcon-backend/common/model"
	eventOpt "github.com/Cepave/open-falcon-backend/modules/alarm/model/event"
	"github.com/astaxie/beego/orm"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	err error
	q   orm.Ormer
	// events []eventOpt.EventCases
	event              eventOpt.EventCases
	runIntegrationTest bool
)

func init() {
	flag.BoolVar(&runIntegrationTest, "integration", false, "run with integration testing. (mysql docker image required)")
}

func TestAlarmBase(t *testing.T) {
	if !runIntegrationTest {
		return
	}
	initTest()
	q = orm.NewOrm()
	q.Using("falcon_portal")
	strategyTemplate := coommonModel.Strategy{
		Id:         1,
		Metric:     "test.m1",
		Tags:       map[string]string{},
		Func:       "all(#2)",
		Operator:   "<",
		RightValue: float64(1),
		MaxStep:    3,
		Priority:   1,
		Note:       "this is a test alarm",
		Tpl: &coommonModel.Template{
			Id:       9,
			Name:     "test template 999",
			ParentId: 0,
			ActionId: 1,
			Creator:  "root",
		},
	}
	// expressTemplate := coommonModel.Expression{
	// 	Id:         1,
	// 	Metric:     "test.m2",
	// 	Tags:       map[string]string{},
	// 	Func:       "all(#2)",
	// 	Operator:   "!=",
	// 	RightValue: float64(-1),
	// 	MaxStep:    3,
	// 	Priority:   1,
	// 	Note:       "this is a test alarm2",
	// 	ActionId:   1,
	// }
	Convey("insert a new alarm", t, func() {
		tnow := int64(1508036708)
		// test insert first alarm case
		eve := coommonModel.Event{
			Id:          "00001",
			Strategy:    &strategyTemplate,
			Expression:  nil,
			Status:      "PROBLEM",
			Endpoint:    "host-001",
			LeftValue:   float64(0),
			CurrentStep: 1,
			EventTime:   tnow,
			PushedTags:  map[string]string{},
		}
		err = eventOpt.InsertEvent(&eve, "owl")
		if err != nil {
			log.Debugf("insert event got error with: %v", err.Error())
		}
		So(err, ShouldBeNil)
		q.Raw("select * from event_cases where id = '00001'").QueryRow(&event)
		So(event.Status, ShouldEqual, "PROBLEM")
		So(event.CurrentStep, ShouldEqual, 1)
		So(event.Timestamp, ShouldEqual, event.UpdateAt)
		So(event.AlarmTypeId, ShouldEqual, 1)

		// test insert second alarm case
		eve.EventTime = tnow + 60
		eve.CurrentStep = 2
		err = eventOpt.InsertEvent(&eve, "owl")
		if err != nil {
			log.Debugf("insert event got error with: %v", err.Error())
		}
		So(err, ShouldBeNil)
		q.Raw("select * from event_cases where id = '00001'").QueryRow(&event)
		So(event.Status, ShouldEqual, "PROBLEM")
		So(event.CurrentStep, ShouldEqual, 2)
		So(event.Timestamp, ShouldNotEqual, event.UpdateAt)

		// test recovered alarm
		eve.EventTime = tnow + 120
		eve.CurrentStep = 1
		eve.Status = "OK"
		err = eventOpt.InsertEvent(&eve, "owl")
		if err != nil {
			log.Debugf("insert event got error with: %v", err.Error())
		}
		So(err, ShouldBeNil)
		q.Raw("select * from event_cases where id = '00001'").QueryRow(&event)
		So(event.Status, ShouldEqual, "OK")
		So(event.CurrentStep, ShouldEqual, 1)
		So(event.Timestamp, ShouldNotEqual, event.UpdateAt)
	})
}
