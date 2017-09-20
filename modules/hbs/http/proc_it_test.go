package http

import (
	"net"
	"net/http"
	"net/http/httptest"

	json "github.com/Cepave/open-falcon-backend/common/json"
	httptesting "github.com/Cepave/open-falcon-backend/common/testing/http"
	. "gopkg.in/check.v1"
)

type TestProcSuite struct{}

var _ = Suite(&TestProcSuite{})

var MOCK_URL = "localhost:5566"

// Tests the expressions service
func (suite *TestProcSuite) TestExpressions(c *C) {
	tsResp := `
	[
		{
			"id":3,
			"metric":"ss.close.wait",
			"tags":{
				"endpoint":"oth-bj-119-090-062-121"
			},
			"func":"all(#1)",
			"operator":"!=",
			"right_value":0,
			"max_step":1,
			"priority":4,
			"note":"boss oth-bj-119-090-062-121 连接数大于10",
			"action_id":91
		}
	]
	`
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tsResp))
		}))
	defer ts.Close()
	l, err := net.Listen("tcp", MOCK_URL)
	c.Assert(err, IsNil)
	ts.Listener = l
	ts.Start()

	slint := httpClientConfig.NewClient().Get("expressions")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Expressions: %s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

// Tests the plugins service
func (suite *TestProcSuite) TestPlugins(c *C) {
	tsResp := `
	["basic/chk","basic/cpu","basic/dev","basic/gd","basic/sys","chk","cpu","file","net"]
	`
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tsResp))
		}))
	defer ts.Close()
	l, err := net.Listen("tcp", MOCK_URL)
	c.Assert(err, IsNil)
	ts.Listener = l
	ts.Start()

	slint := httpClientConfig.NewClient().Get("plugins/test-host")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Plugins: %s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

func (s *TestProcSuite) SetUpSuite(c *C) {
	if !testFlags.HasHttpClient() {
		c.Skip("Skipping testng because properties HTTP client is missing")
		return
	}

	c.Logf("Testing service of proc: %s", httpClientConfig)
}
