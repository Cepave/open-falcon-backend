package rpc

import (
	"github.com/Cepave/common/model"
	hbstesting "github.com/Cepave/hbs/testing"
	. "gopkg.in/check.v1"
	"sort"
)

type TestRpcNqmAgentSuite struct{}

var _ = Suite(&TestRpcNqmAgentSuite{})

/**
 * Tests the data validation for ping task
 */
func (suite *TestRpcNqmAgentSuite) TestValidatePingTask(c *C) {
	var testeCases = [][]interface{} {
		{ "120.49.58.19", "localhost.localdomain", "120.49.58.19", IsNil },
		{ "", "localhost.localdomain", "120.49.58.19", NotNil },
		{ "120.49.58.19", "", "120.49.58.19", NotNil },
		{ "120.49.58.19", "localhost.localdomain", "", NotNil },
	}

	for _, v := range testeCases {
		err := validatePingTask(
			&model.NqmPingTaskRequest{
				ConnectionId: v[0].(string),
				Hostname: v[1].(string),
				IpAddress: v[2].(string),
			},
		)

		c.Assert(err, v[3].(Checker))
	}
}

/**
 * Tests the data content of ping task
 */
type byId []model.NqmTarget
func (targets byId) Len() int           { return len(targets) }
func (targets byId) Swap(i, j int)      { targets[i], targets[j] = targets[j], targets[i] }
func (targets byId) Less(i, j int) bool { return targets[i].Id < targets[j].Id }
func (suite *TestRpcNqmAgentSuite) TestPingTask(c *C) {
	var req = model.NqmPingTaskRequest {
		ConnectionId: "ag-rpc-1",
		Hostname: "rpc-1.org",
		IpAddress: "45.65.0.1",
	}
	var resp model.NqmPingTaskResponse

	hbstesting.DefaultListenAndExecute(
		new(NqmAgent),
		func(rpcTestEnvInstance *hbstesting.RpcTestEnv) {
			err := rpcTestEnvInstance.RpcClient.Call(
				"NqmAgent.PingTask", req, &resp,
			)

			/**
			 * Asserts the agent
			 */
			c.Assert(err, IsNil)
			c.Assert(resp.NeedPing, Equals, true)
			c.Assert(resp.Agent.Id, Equals, 4051)
			c.Assert(resp.Agent.IspId, Equals, int16(3))
			c.Assert(resp.Agent.ProvinceId, Equals, int16(2))
			c.Assert(resp.Agent.CityId, Equals, model.UNDEFINED_CITY_ID)

			c.Assert(len(resp.Targets), Equals, 3)
			c.Assert(resp.Command[0], Equals, "fping")
			// :~)

			/**
			 * Asserts the targets
			 */
			sort.Sort(byId(resp.Targets))
			var expectedTargets = []model.NqmTarget {
				{ Id: 6301, Host: "1.2.3.4", IspId: 1, ProvinceId: 4, CityId: -1, NameTag: model.UNDEFINED_STRING, },
				{ Id: 6302, Host: "1.2.3.5", IspId: 2, ProvinceId: 4, CityId: -1, NameTag: "tag-1", },
				{ Id: 6303, Host: "1.2.3.6", IspId: 3, ProvinceId: 4, CityId: -1, NameTag: model.UNDEFINED_STRING, },
			}

			c.Assert(resp.Targets, DeepEquals, expectedTargets)
			// :~)
		},
	)
}

func (s *TestRpcNqmAgentSuite) SetUpSuite(c *C) {
	(&TestRpcSuite{}).SetUpSuite(c)
}
func (s *TestRpcNqmAgentSuite) TearDownSuite(c *C) {
	(&TestRpcSuite{}).TearDownSuite(c)
}

func (s *TestRpcNqmAgentSuite) SetUpTest(c *C) {
	switch c.TestName() {
	case "TestRpcNqmAgentSuite.TestPingTask":
		if !hbstesting.HasDbEnvForMysqlOrSkip(c) {
			return
		}

		hbstesting.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_pv_id, ag_ct_id)
			VALUES (4051, 'ag-rpc-1', 'rpc-1.org', 0x12345672, 3, 2, -1)
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_probed_by_all, tg_name_tag
			)
			VALUES
				(6301, 'tgn-1', '1.2.3.4', 1, 4, -1, true, null),
				(6302, 'tgn-2', '1.2.3.5', 2, 4, -1, true, 'tag-1'),
				(6303, 'tgn-3', '1.2.3.6', 3, 4, -1, true, null)
			`,
			`
			INSERT INTO nqm_ping_task(pt_ag_id, pt_period)
			VALUES(4051, 10)
			`,
		)
	}
}

func (s *TestRpcNqmAgentSuite) TearDownTest(c *C) {
	switch c.TestName() {
	case "TestRpcNqmAgentSuite.TestPingTask":
		hbstesting.ExecuteQueriesOrFailInTx(
			"DELETE FROM nqm_ping_task",
			"DELETE FROM nqm_target",
			"DELETE FROM nqm_agent",
		)
	}
}
