package restful

import (
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	testingDb "github.com/Cepave/open-falcon-backend/modules/nqm-conf-mng/testing"

	rdb "github.com/Cepave/open-falcon-backend/modules/nqm-conf-mng/rdb"

	"github.com/dghubble/sling"

	. "gopkg.in/check.v1"
)

type TestAgentItSuite struct{}

var _ = Suite(&TestAgentItSuite{})

func (suite *TestAgentItSuite) TestListAgents(c *C) {
	client := sling.New().Get(httpClientConfig.String()).
		Path("/api/v1/nqm/agents")

	slintChecker := testingHttp.NewCheckSlint(c, client)

	slintChecker.AssertHasPaging()
	message := slintChecker.GetJsonBody(200)

	prettyJson, err := message.EncodePretty()

	c.Assert(err, IsNil)
	c.Logf("JSON Result: %s", prettyJson)
	c.Assert(len(message.MustArray()), Equals, 3)
}

func (s *TestAgentItSuite) SetUpSuite(c *C) {
	testingDb.InitRdb(c)
}
func (s *TestAgentItSuite) TearDownSuite(c *C) {
	testingDb.ReleaseRdb(c)
}

func (s *TestAgentItSuite) SetUpTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentItSuite.TestListAgents":
		inTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(22091, 'agent-it-01', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id)
			VALUES(4321, 22091, 'agent-it-01', 'agent-01@28.71.19.22', 'agent-01.fb.com', x'1C471316', 7),
				(4322, 22091, 'agent-it-02', 'agent-02@28.71.19.23', 'agent-02.fb.com', x'1C471317', 7),
				(4323, 22091, 'agent-it-03', 'agent-03@28.71.19.23', 'agent-03.fb.com', x'1C471318', 7)
			`,
		)
	}
}
func (s *TestAgentItSuite) TearDownTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentItSuite.TestListAgents":
		inTx(
			"DELETE FROM nqm_agent WHERE ag_id >= 4321 AND ag_id <= 4323",
			"DELETE FROM host WHERE id = 22091",
		)
	}
}
