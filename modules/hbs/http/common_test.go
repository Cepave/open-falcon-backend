package http

import (
	"testing"
	httptesting "github.com/Cepave/open-falcon-backend/common/testing/http"
	json "github.com/Cepave/open-falcon-backend/common/json"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestCommonSuite struct{}

var _ = Suite(&TestCommonSuite{})

var httpClientConfig = httptesting.NewHttpClientConfigByFlag()

// Tests the health
func (suite *TestCommonSuite) TestHealth(c *C) {
	slint := httpClientConfig.NewSlingByBase().Get("health")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	message := slintChecker.GetStringBody(200)

	c.Assert(message, Equals, "ok")
}

// Tests the version
func (suite *TestCommonSuite) TestVersion(c *C) {
	slint := httpClientConfig.NewSlingByBase().Get("version")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	message := slintChecker.GetStringBody(200)

	c.Logf("Test version[/version]: %s", message)
}

// Tests the workdir
func (suite *TestCommonSuite) TestWorkdir(c *C) {
	slint := httpClientConfig.NewSlingByBase().Get("workdir")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Test workdir[/workdir]:\n%s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

// Tests the reload of configuration
func (suite *TestCommonSuite) TestReloadConfig(c *C) {
	slint := httpClientConfig.NewSlingByBase().Get("config/reload")

	slintChecker := httptesting.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Test reloading of configuration[/config/reload]:\n%s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

func (s *TestCommonSuite) SetUpSuite(c *C) {
	c.Logf("Testing service of common: %s", httpClientConfig)
}
