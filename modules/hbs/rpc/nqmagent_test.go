package rpc

import (
	"sort"
	"net/rpc"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	testJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"
	. "gopkg.in/check.v1"
)

type TestRpcNqmAgentSuite struct{}

var _ = Suite(&TestRpcNqmAgentSuite{})

/**
 * Tests the data validation for ping task
 */
func (suite *TestRpcNqmAgentSuite) TestValidatePingTask(c *C) {
	var testeCases = []struct {
		connectionId string
		hostname     string
		ipAddress    string
		checker      Checker
	}{
		{"120.49.58.19", "localhost.localdomain", "120.49.58.19", IsNil},
		{"", "localhost.localdomain", "120.49.58.19", NotNil},
		{"120.49.58.19", "", "120.49.58.19", NotNil},
		{"120.49.58.19", "localhost.localdomain", "", NotNil},
	}

	for _, v := range testeCases {
		err := validatePingTask(
			&model.NqmTaskRequest{
				ConnectionId: v.connectionId,
				Hostname:     v.hostname,
				IpAddress:    v.ipAddress,
			},
		)

		c.Assert(err, v.checker)
	}
}

/**
 * Tests the data from Task()
 */
type byID []model.NqmTarget

func (targets byID) Len() int           { return len(targets) }
func (targets byID) Swap(i, j int)      { targets[i], targets[j] = targets[j], targets[i] }
func (targets byID) Less(i, j int) bool { return targets[i].Id < targets[j].Id }
func (suite *TestRpcNqmAgentSuite) TestTask(c *C) {
	var req = model.NqmTaskRequest{
		ConnectionId: "ag-rpc-1",
		Hostname:     "rpc-1.org",
		IpAddress:    "45.65.0.1",
	}
	var resp model.NqmTaskResponse

	testJsonRpc.OpenClient(c, func(jsonRpcClient *rpc.Client) {
		err := jsonRpcClient.Call(
			"NqmAgent.Task", req, &resp,
		)

		/**
		 * Asserts the agent
		 */
		c.Assert(err, IsNil)
		c.Logf("Got response: %v", &resp)
		c.Logf("Agent : %v", resp.Agent)

		c.Assert(resp.NeedPing, Equals, true)
		c.Assert(resp.Agent.Id, Equals, 405001)
		c.Assert(resp.Agent.Name, Equals, "ag-name-1")
		c.Assert(resp.Agent.IspId, Equals, int16(3))
		c.Assert(resp.Agent.IspName, Equals, "移动")
		c.Assert(resp.Agent.ProvinceId, Equals, int16(2))
		c.Assert(resp.Agent.ProvinceName, Equals, "山西")
		c.Assert(resp.Agent.CityId, Equals, model.UNDEFINED_CITY_ID)
		c.Assert(resp.Agent.CityName, Equals, model.UNDEFINED_STRING)

		c.Assert(len(resp.Targets), Equals, 3)
		c.Assert(resp.Measurements["fping"].Command[0], Equals, "fping")
		// :~)

		/**
		 * Asserts the 1st target
		 */
		for _, v := range resp.Targets {
			c.Logf("Target: %v", &v)
		}

		sort.Sort(byID(resp.Targets))

		c.Assert(
			resp.Targets[0], DeepEquals,
			model.NqmTarget{
				Id: 630001, Host: "1.2.3.4",
				IspId: 1, IspName: "北京三信时代",
				ProvinceId: 4, ProvinceName: "北京",
				CityId: model.UNDEFINED_CITY_ID, CityName: model.UNDEFINED_STRING,
				NameTagId: model.UNDEFINED_NAME_TAG_ID, NameTag: model.UNDEFINED_STRING,
			},
		)
		// :~)
	})
}

func (s *TestRpcNqmAgentSuite) SetUpSuite(c *C) {
	db.DbInit(dbTest.GetDbConfig(c))
}
func (s *TestRpcNqmAgentSuite) TearDownSuite(c *C) {
	db.Release()
}

func (s *TestRpcNqmAgentSuite) SetUpTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestRpcNqmAgentSuite.TestTask":
		executeInTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES (9901, 'tag-1')
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(54091, 'rpc-1.org', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_pv_id, ag_ct_id)
			VALUES (405001, 54091, 'ag-name-1', 'ag-rpc-1', 'rpc-1.org', 0x12345672, 3, 2, -1)
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_probed_by_all, tg_nt_id, tg_available, tg_status
			)
			VALUES
				(630001, 'tgn-1', '1.2.3.4', 1, 4, -1, true, -1, true, true),
				(630002, 'tgn-2', '1.2.3.5', 2, 4, -1, true, 9901, true, true),
				(630003, 'tgn-3', '1.2.3.6', 3, 4, -1, true, -1, true, true)
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_period)
			VALUES(32001, 10)
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES(405001, 32001)
			`,
		)
	}
}

func (s *TestRpcNqmAgentSuite) TearDownTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestRpcNqmAgentSuite.TestTask":
		executeInTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id = 405001",
			"DELETE FROM nqm_ping_task WHERE pt_id = 32001",
			"DELETE FROM nqm_agent WHERE ag_id = 405001",
			"DELETE FROM host WHERE id = 54091",
			"DELETE FROM nqm_target WHERE tg_id >= 630001 AND tg_id <= 630003",
			"DELETE FROM owl_name_tag WHERE nt_id = 9901",
		)
	}
}

