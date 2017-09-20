package http

import (
	"testing"

	json "github.com/Cepave/open-falcon-backend/common/json"
	th "github.com/Cepave/open-falcon-backend/common/testing/http"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestCommonSuite struct{}

var _ = Suite(&TestCommonSuite{})

var httpClientConfig = th.SlingClientConf{itConfig}

// Tests the health
func (suite *TestCommonSuite) TestHealth(c *C) {
	slint := httpClientConfig.NewClient().Get("health")

	slintChecker := th.NewCheckSlint(c, slint)

	message := slintChecker.GetStringBody(200)

	c.Assert(message, Equals, "ok")
}

// Tests the version
func (suite *TestCommonSuite) TestVersion(c *C) {
	slint := httpClientConfig.NewClient().Get("version")

	slintChecker := th.NewCheckSlint(c, slint)

	message := slintChecker.GetStringBody(200)

	c.Logf("Test version[/version]: %s", message)
}

// Tests the workdir
func (suite *TestCommonSuite) TestWorkdir(c *C) {
	slint := httpClientConfig.NewClient().Get("workdir")

	slintChecker := th.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Test workdir[/workdir]:\n%s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

// Tests the reload of configuration
func (suite *TestCommonSuite) TestReloadConfig(c *C) {
	slint := httpClientConfig.NewClient().Get("config/reload")

	slintChecker := th.NewCheckSlint(c, slint)

	jsonMessage := slintChecker.GetJsonBody(200)

	c.Logf("Test reloading of configuration[/config/reload]:\n%s", json.MarshalPrettyJSON(jsonMessage))
	c.Assert(jsonMessage.Get("msg").MustString(), Equals, "success")
}

func (s *TestCommonSuite) SetUpSuite(c *C) {
	if !testFlags.HasHttpClient() {
		c.Skip("Skipping testng because properties HTTP client is missing")
		return
	}
	c.Logf("Testing service of common: %s", httpClientConfig)
}
