package PortalTest

import (
	"fmt"
	"testing"

	"github.com/Jeffail/gabs"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPortalEventCase(t *testing.T) {
	Convey("UIC for account testing", t, func() {
		Convey("get login session", func() {
			resp := DoPost("/api/v1/auth/login", fmt.Sprintf(`name=%s;password=%s`, name, password))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			session := jsParsed.Search("data", "sig").Data().(string)
			So(session, ShouldNotBeEmpty)
		})
	})
	const caseId = "s_118_b66b973ef551e4e503fad475dfc9e418"
	GetAuthSessoion()
	Convey("Portal - Test alarm Cases", t, func() {
		Convey("test alerts with Timestamp range set up", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d", name, session, 1465920000, 1465981200))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldBeGreaterThan, 0)
			So(ACaseIds[0].Search("id").Data().(string), ShouldEqual, "s_118_b66b973ef551e4e503fad475dfc9e418")
		})
		Convey("test alerts with status found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d;status=%s", name, session, 1465920000, 1465981200, "PROBLEM"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldBeGreaterThan, 0)
			So(ACaseIds[0].Search("id").Data().(string), ShouldEqual, "s_118_b66b973ef551e4e503fad475dfc9e418")
		})
		Convey("test alerts with status not found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d;status=%s", name, session, 1465920000, 1465981200, "OK"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldEqual, 0)
		})
		Convey("test alerts with process status found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d;process_status=%s", name, session, 1465920000, 1465981200, "ignored,in progress,unresolved,resolved"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldBeGreaterThan, 0)
			So(ACaseIds[0].Search("id").Data().(string), ShouldEqual, "s_118_b66b973ef551e4e503fad475dfc9e418")
		})
		Convey("test alerts with process status not found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d;process_status=%s", name, session, 1465920000, 1465981200, "test"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldEqual, 0)
		})
		Convey("test alerts with metrics found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d;metrics=%s", name, session, 1465920000, 1465981200, "cpu.idle,df.statistics.total,net.if.in.bits/iface=eth0"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldBeGreaterThan, 0)
			So(ACaseIds[0].Search("id").Data().(string), ShouldEqual, "s_118_b66b973ef551e4e503fad475dfc9e418")
		})
		Convey("test alerts with metrics not found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d;metrics=%s", name, session, 1465920000, 1465981200, "test"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldEqual, 0)
		})
		Convey("test alerts with caseId found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;caseId=%s", name, session, caseId))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldEqual, 1)
			So(ACaseIds[0].Search("id").Data().(string), ShouldEqual, "s_118_b66b973ef551e4e503fad475dfc9e418")
		})
		Convey("test alerts with caseId not found", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;caseId=%s", name, session, "testid"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldEqual, 0)
		})
		Convey("test alerts with limit", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;limit=%d", name, session, 2))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(len(ACaseIds), ShouldEqual, 2)
		})
		Convey("test alerts with elimit", func() {
			resp := DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;limit=%d;elimit=%d", name, session, 1, 2))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			AEvent, _ := ACaseIds[0].Search("evevnts").Children()
			So(len(AEvent), ShouldEqual, 2)
		})
	})

}
