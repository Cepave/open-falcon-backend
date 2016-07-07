package PortalTest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Jeffail/gabs"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPortalEvents(t *testing.T) {
	const caseId = "s_118_b66b973ef551e4e503fad475dfc9e418"
	GetAuthSessoion()
	Convey("Portal - Test Events", t, func() {
		Convey("test events with caseId", func() {
			resp := DoPost("/api/v1/portal/events/get", fmt.Sprintf("cName=%s;cSig=%s;caseId=%s", name, session, caseId))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			events, _ := jsParsed.Search("data", "events").Children()
			So(len(events), ShouldBeGreaterThan, 0)
			So(strings.Index(resp.Body, "success"), ShouldBeGreaterThan, 0)
			So(events[0].Search("event_caseId").Data().(string), ShouldEqual, caseId)
		})
		Convey("test events with timeRange", func() {
			resp := DoPost("/api/v1/portal/events/get", fmt.Sprintf("cName=%s;cSig=%s;startTime=%d;endTime=%d", name, session, 1467216000, 1467302400))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			events, _ := jsParsed.Search("data", "events").Children()
			So(len(events), ShouldBeGreaterThan, 0)
			So(strings.Index(resp.Body, "success"), ShouldBeGreaterThan, 0)
			So(events[0].Search("event_caseId").Data().(string), ShouldEqual, "s_128_d5049645592b60591502c626cbb125bf")
		})
		Convey("test events with status 'OK' ", func() {
			resp := DoPost("/api/v1/portal/events/get", fmt.Sprintf("cName=%s;cSig=%s;caseId=%s;status=%s", name, session, "s_128_d5049645592b60591502c626cbb125bf", "OK"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			events, _ := jsParsed.Search("data", "events").Children()
			So(len(events), ShouldEqual, 1)
			So(strings.Index(resp.Body, "success"), ShouldBeGreaterThan, 0)
		})
		Convey("test events with status 'PROBLEM' ", func() {
			resp := DoPost("/api/v1/portal/events/get", fmt.Sprintf("cName=%s;cSig=%s;caseId=%s;status=%s", name, session, "s_128_d5049645592b60591502c626cbb125bf", "PROBLEM"))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			events, _ := jsParsed.Search("data", "events").Children()
			So(len(events), ShouldEqual, 1)
			So(strings.Index(resp.Body, "success"), ShouldBeGreaterThan, 0)
		})
		Convey("test events with limit ", func() {
			resp := DoPost("/api/v1/portal/events/get", fmt.Sprintf("cName=%s;cSig=%s;limit=%d", name, session, 3))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			events, _ := jsParsed.Search("data", "events").Children()
			So(len(events), ShouldEqual, 3)
		})
	})
}
