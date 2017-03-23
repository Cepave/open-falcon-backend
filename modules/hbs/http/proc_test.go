package http

import (
	httptesting "github.com/Cepave/open-falcon-backend/common/testing/http"
	json "github.com/Cepave/open-falcon-backend/common/json"
	. "gopkg.in/check.v1"
)

type TestProcSuite struct{}

var _ = Suite(&TestProcSuite{})

// Tests the expressions service
func (suite *TestProcSuite) TestExpressions(c *C) {
	slint := httpClientConfig.NewSlingByBase().Get("expressions")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Expressions: %s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

// Tests the plugins service
func (suite *TestProcSuite) TestPlugins(c *C) {
	slint := httpClientConfig.NewSlingByBase().Get("plugins/")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Plugins: %s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

func (s *TestProcSuite) SetUpSuite(c *C) {
	c.Logf("Testing service of proc: %s", httpClientConfig)
}
