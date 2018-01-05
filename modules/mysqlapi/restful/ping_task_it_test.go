package restful

import (
	"net/http"

	json "github.com/Cepave/open-falcon-backend/common/json"
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	testingDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/testing"

	rdb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"

	. "gopkg.in/check.v1"
)

type TestPintTaskItSuite struct{}

var _ = Suite(&TestPintTaskItSuite{})

// Tests the listing of agents
func (suite *TestPintTaskItSuite) TestListAgents(c *C) {
	client := httpClientConfig.NewClient().Get("api/v1/nqm/pingtask/9801/agents").
		Set("order-by", "applied#desc")

	slintChecker := testingHttp.NewCheckSlint(c, client)

	slintChecker.AssertHasPaging()
	message := slintChecker.GetJsonBody(http.StatusOK)

	c.Logf("[List Agents] JSON Result: %s", json.MarshalPrettyJSON(message))
	c.Assert(len(message.MustArray()), Equals, 3)

	// Ordered by applied status of agent
	c.Assert(message.GetIndex(0).Get("id").MustInt(), Equals, 5621)

	/**
	 * Asserts the applied status
	 */
	c.Assert(message.GetIndex(0).Get("applying_ping_task").MustBool(), Equals, true)
	c.Assert(message.GetIndex(1).Get("applying_ping_task").MustBool(), Equals, false)
	c.Assert(message.GetIndex(2).Get("applying_ping_task").MustBool(), Equals, false)
	// :~)
}

func (s *TestPintTaskItSuite) SetUpSuite(c *C) {
	itSkipForGocheck(c)
	testingDb.InitRdb(c)
}
func (s *TestPintTaskItSuite) TearDownSuite(c *C) {
	testingDb.ReleaseRdb(c)
}

func (s *TestPintTaskItSuite) SetUpTest(c *C) {
	inPortalTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestPintTaskItSuite.TestListAgents":
		inPortalTx(
			`
			INSERT INTO nqm_ping_task(pt_id, pt_name, pt_period)
			VALUES(9801, 'ag-in-city', 40)
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(19531, 'agent-it-01', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id)
			VALUES(5621, 19531, 'agent-it-01', 'agent-01@28.71.19.22', 'agent-01.fb.com', x'1C471316', 7),
				(5622, 19531, 'agent-it-02', 'agent-02@28.71.19.23', 'agent-02.fb.com', x'1C471317', 7),
				(5623, 19531, 'agent-it-03', 'agent-03@28.71.19.23', 'agent-03.fb.com', x'1C471318', 7)
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES(5621, 9801)
			`,
		)
	}
}
func (s *TestPintTaskItSuite) TearDownTest(c *C) {
	inPortalTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestPintTaskItSuite.TestListAgents":
		inPortalTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_pt_id = 9801",
			"DELETE FROM nqm_agent WHERE ag_id >= 5621 AND ag_id <= 5623",
			"DELETE FROM nqm_ping_task WHERE pt_id = 9801",
			"DELETE FROM host WHERE id = 19531",
		)
	}
}
