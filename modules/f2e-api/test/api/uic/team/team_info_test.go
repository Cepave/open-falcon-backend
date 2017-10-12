package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
*	authapi_team.GET("/team", Teams)
*	authapi_team.GET("/team/:team_id", GetTeam)
 */

func TestTeamGetInfo(t *testing.T) {
	routes := SetUpGin()
	Convey("get team info & list", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("get team list", func() {
			w, r = NewTestContextWithDefaultSession("GET", "/api/v1/team", nil)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			checkR := gjson.Get(respBody, "#.Team")
			So(len(checkR.Array()), ShouldBeGreaterThanOrEqualTo, 3)
			So(w.Code, ShouldEqual, 200)
		})
		Convey("get team info by id", func() {
			w, r = NewTestContextWithDefaultSession("GET", "/api/v1/team/1", nil)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			checkR := gjson.Get(respBody, "users.#.id")
			So(respBody, ShouldContainSubstring, "\"name\"")
			So(len(checkR.Array()), ShouldBeGreaterThanOrEqualTo, 3)
			So(w.Code, ShouldEqual, 200)
		})
	})
}
