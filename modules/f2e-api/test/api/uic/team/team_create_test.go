package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
 *	authapi_team.POST("/team", CreateTeam)
 */

func TestTeamCreate(t *testing.T) {
	routes := SetUpGin()
	Convey("create a new team", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("create a new team ok", func() {
			/* insert team 1 */
			postb := map[string]interface{}{
				"team_name": "team_X",
				"resume":    "this is resumeA",
				"users":     []int{1, 2, 3},
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "team created")
			So(w.Code, ShouldEqual, 200)
			///* insert team 2 */
			//postb = map[string]interface{}{
			//	"team_name": "team_B",
			//	"resume":    "this is resumeB",
			//	"users":     []int{1},
			//}
			//b, _ = json.Marshal(postb)
			//w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			//routes.ServeHTTP(w, r)
			//respBody = w.Body.String()
			//So(respBody, ShouldContainSubstring, "team created")
			//So(w.Code, ShouldEqual, 200)
			///* insert team 3 */
			//postb = map[string]interface{}{
			//	"team_name": "team_C",
			//	"resume":    "this is resumeC",
			//	"users":     []int{2, 3},
			//}
			//b, _ = json.Marshal(postb)
			//w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			//routes.ServeHTTP(w, r)
			//respBody = w.Body.String()
			//So(respBody, ShouldContainSubstring, "team created")
			//So(w.Code, ShouldEqual, 200)
			///* insert team 4 */
			//postb = map[string]interface{}{
			//	"team_name": "team_D",
			//	"resume":    "this is resumeD",
			//	"users":     []int{},
			//}
			//b, _ = json.Marshal(postb)
			//w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			//routes.ServeHTTP(w, r)
			//respBody = w.Body.String()
			//So(respBody, ShouldContainSubstring, "team created")
			//So(w.Code, ShouldEqual, 200)
			///* insert team 5 */
			//postb = map[string]interface{}{
			//	"team_name": "team_D1",
			//	"resume":    "this is resumeD1",
			//	"users":     []int{1, 2},
			//}
			//b, _ = json.Marshal(postb)
			//w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			//routes.ServeHTTP(w, r)
			//respBody = w.Body.String()
			//So(respBody, ShouldContainSubstring, "team created")
			//So(w.Code, ShouldEqual, 200)
		})
		Convey("create a new team faild", func() {
			postb := map[string]interface{}{
				"resume":  "this is resume3",
				"userIDs": []int{2, 3},
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 400)
		})
	})
}
