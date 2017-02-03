package rpc

import (
	"net/rpc"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"

	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	testJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"

	. "gopkg.in/check.v1"
)

type TestNqmAgentSuite struct{}

var _ = Suite(&TestNqmAgentSuite{})

// Tests the refreshing and retrieving list of targets for NQM agent
func (suite *TestNqmAgentSuite) TestTask(c *C) {
	testCases := []*struct {
		req model.NqmTaskRequest
		needPing bool
		expectedTargetIds []int32
	} {
		{
			model.NqmTaskRequest{
				ConnectionId: "ag-rpc-1",
				Hostname:     "rpc-1.org",
				IpAddress:    "45.65.0.1",
			},
			true, []int32{ 630001, 630002, 630003 },
		},
		{
			model.NqmTaskRequest{
				ConnectionId: "ag-rpc-2",
				Hostname:     "rpc-2.org",
				IpAddress:    "45.65.0.2",
			},
			false, []int32{},
		},
		{ // The period is not elapsed yet
			model.NqmTaskRequest{
				ConnectionId: "ag-rpc-1",
				Hostname:     "rpc-1.org",
				IpAddress:    "45.65.0.1",
			},
			false, []int32{},
		},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		var resp model.NqmTaskResponse
		testJsonRpc.OpenClient(c, func(jsonRpcClient *rpc.Client) {
			err := jsonRpcClient.Call(
				"NqmAgent.Task", testCase.req, &resp,
			)

			c.Assert(err, IsNil)
		})

		connectionId := testCase.req.ConnectionId
		c.Logf("[%s] Response Need Ping: %v", connectionId, resp.NeedPing)
		c.Logf("[%s] Response Agent: %#v", connectionId, resp.Agent)
		c.Logf("[%s] Response Measurements: %#v", connectionId, resp.Measurements)

		/**
		 * The case with no needed of PING
		 */
		if !testCase.needPing {
			c.Assert(resp.NeedPing, Equals, false, comment)
			c.Assert(resp.Agent, IsNil)
			continue
		}
		// :~)

		c.Assert(resp.NeedPing, Equals, true)
		c.Assert(
			resp.Agent, DeepEquals,
			&model.NqmAgent {
				Id:405001,
				Name:"ag-name-1",
				IspId:3, IspName:"移动",
				ProvinceId:2, ProvinceName:"山西",
				CityId:-1, CityName:"<UNDEFINED>",
				NameTagId:-1,
				GroupTagIds:[]int32{9081, 9082},
			},
		)

		c.Assert(len(resp.Targets), Equals, 3)
		c.Assert(resp.Measurements["fping"].Command[0], Equals, "fping")

		c.Assert(resp.Targets, HasLen, len(testCase.expectedTargetIds), comment)
		for i, targetId := range testCase.expectedTargetIds {
			c.Logf("\tTarget: %#v", resp.Targets[i])
			c.Assert(int32(resp.Targets[i].Id), Equals, targetId, comment)
		}
	}
}

func (s *TestNqmAgentSuite) SetUpSuite(c *C) {
	db.DbInit(dbTest.GetDbConfig(c))
}
func (s *TestNqmAgentSuite) TearDownSuite(c *C) {
	db.Release()
}

func (s *TestNqmAgentSuite) SetUpTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNqmAgentSuite.TestTask":
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
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(9081, 'rpc-gt-1'), (9082, 'rpc-gt-2')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_pv_id, ag_ct_id, ag_status)
			VALUES (405001, 54091, 'ag-name-1', 'ag-rpc-1', 'rpc-1.org', 0x12345672, 3, 2, -1, TRUE),
				(405002, 54091, 'ag-name-2', 'ag-rpc-2', 'rpc-2.org', 0x12345673, 3, 2, -1, FALSE)
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id, tg_available, tg_status
			)
			VALUES
				(630001, 'tgn-1', '1.2.3.4', 1, 4, 1, -1, true, true),
				(630002, 'tgn-2', '1.2.3.5', 2, 5, 280, 9901, true, true),
				(630003, 'tgn-3', '1.2.3.6', 3, 5, 285, -1, true, true)
			`,
			`
			INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
			VALUES(405001, 9081), (405001, 9082)
			`,
			`
			INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
			VALUES(630001, 9081), (630001, 9082),
				(630002, 9081), (630003, 9082)
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_period)
			VALUES(32001, 10)
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES(405001, 32001), (405002, 32001)
			`,
		)
	}
}

func (s *TestNqmAgentSuite) TearDownTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNqmAgentSuite.TestTask":
		executeInTx(
			"DELETE FROM nqm_cache_agent_ping_list_log WHERE apll_ag_id >= 405001 AND apll_ag_id <= 405002",
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 405001 AND apt_ag_id <= 405002",
			"DELETE FROM nqm_ping_task WHERE pt_id = 32001",
			"DELETE FROM nqm_agent WHERE ag_id >= 405001 AND ag_id <= 405002",
			"DELETE FROM host WHERE id = 54091",
			"DELETE FROM nqm_target WHERE tg_id >= 630001 AND tg_id <= 630003",
			"DELETE FROM owl_name_tag WHERE nt_id = 9901",
			"DELETE FROM owl_group_tag WHERE gt_id >= 9081 AND gt_id <= 9082",
		)
	}
}

