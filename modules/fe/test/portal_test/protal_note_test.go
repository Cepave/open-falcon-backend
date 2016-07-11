package PortalTest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPortalNote(t *testing.T) {
	const caseId = "s_118_b66b973ef551e4e503fad475dfc9e418"
	GetAuthSessoion()
	Convey("Portal - Test Note", t, func() {
		note := fmt.Sprintf("this is test note - %s", time.Now().Format(time.RFC822))
		processStatus := "in progress"
		Convey("add note", func() {
			postData := fmt.Sprintf("cName=%s;cSig=%s;id=%s;note=%s;status=%s", name, session, caseId, note, processStatus)
			resp := DoPost("/api/v1/portal/eventcases/addnote", postData)
			So(strings.Index(resp.Body, "success"), ShouldBeGreaterThan, 0)
			//check event case should update after note added
			resp = DoPost("/api/v1/portal/eventcases/get", fmt.Sprintf("cName=%s;cSig=%s;caseId=%s", name, session, caseId))
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			ACaseIds, _ := jsParsed.Search("data", "eventCases").Children()
			So(ACaseIds[0].Search("process_status").Data().(string), ShouldEqual, processStatus)
		})

		Convey("get note with id", func() {
			postData := fmt.Sprintf("cName=%s;cSig=%s;id=%s", name, session, caseId)
			resp := DoPost("/api/v1/portal/eventcases/notes", postData)
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			notes, _ := jsParsed.Search("data", "notes").Children()
			So(len(notes), ShouldBeGreaterThan, 0)
			So(strings.Index(resp.Body, "success"), ShouldBeGreaterThan, 0)
			So(notes[0].Search("note").Data().(string), ShouldEqual, note)
			So(notes[0].Search("status").Data().(string), ShouldEqual, processStatus)
		})

		Convey("get note with id and timerange and filterIgnored", func() {
			postData := fmt.Sprintf("cName=%s;cSig=%s;id=%s;startTime=%d;endTime=%d;filterIgnored=true", name, session, caseId, 1466784000, 1467302400)
			resp := DoPost("/api/v1/portal/eventcases/notes", postData)
			jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
			notes, _ := jsParsed.Search("data", "notes").Children()
			So(len(notes), ShouldEqual, 1)
			So(strings.Index(resp.Body, "success"), ShouldBeGreaterThan, 0)
		})
	})
}
