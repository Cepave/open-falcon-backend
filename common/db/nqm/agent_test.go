package nqm

import (
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

// Tests the list of agents with various conditions
func (suite *TestAgentSuite) TestListAgents(c *C) {
	testCases := []struct {
		query *nqmModel.AgentQuery
		pageSize int32
		pagePosition int32
		expectedCountOfCurrentPage int
		expectedCountOfAll int32
	} {
		{ // All data
			&nqmModel.AgentQuery {},
			10, 1, 3, 3,
		},
		{ // 2nd page
			&nqmModel.AgentQuery {},
			2, 2, 1, 3,
		},
		{ // Match nothing for futher page
			&nqmModel.AgentQuery {},
			10, 10, 0, 3,
		},
		{ // Match 1 row by all of the conditions
			&nqmModel.AgentQuery {
				Name: "ag-name-1",
				ConnectionId: "ag-list-1",
				Hostname: "hn-list-1",
				HasIspId: true,
				IspId: 3,
				IpAddress: "123.52",
				HasStatusCondition: true,
				Status: true,
			}, 10, 1, 1, 1,
		},
		{ // Match 1 row(by special IP address)
			&nqmModel.AgentQuery {
				IpAddress: "12.37",
			}, 10, 1, 1, 1,
		},
		{ // Match nothing
			&nqmModel.AgentQuery {
				ConnectionId: "ag-list-1",
				Hostname: "hn-list-2",
			}, 10, 1, 0, 0,
		},
	}

	for _, testCase := range testCases {
		paging := commonModel.Paging{
			Size: testCase.pageSize,
			Position: testCase.pagePosition,
			OrderBy: []*commonModel.OrderByEntity {
				&commonModel.OrderByEntity{ "status", commonModel.Ascending },
				&commonModel.OrderByEntity{ "name", commonModel.Ascending },
				&commonModel.OrderByEntity{ "connection_id", commonModel.Ascending },
				&commonModel.OrderByEntity{ "comment", commonModel.Ascending },
				&commonModel.OrderByEntity{ "province", commonModel.Ascending },
				&commonModel.OrderByEntity{ "city", commonModel.Ascending },
				&commonModel.OrderByEntity{ "last_heartbeat_time", commonModel.Ascending },
				&commonModel.OrderByEntity{ "name_tag", commonModel.Ascending },
				&commonModel.OrderByEntity{ "group_tag", commonModel.Descending },
			},
		}

		testedResult, newPaging := ListAgents(
			testCase.query, paging,
		)

		c.Logf("[List] Query condition: %v. Number of agents: %d", testCase.query, len(testedResult))

		for _, agent := range testedResult {
			c.Logf("[List] Agent: %v.", agent)
			//c.Assert(agent.IspId, Equals, agent.Isp.Id)
		}
		c.Assert(testedResult, HasLen, testCase.expectedCountOfCurrentPage)
		c.Assert(newPaging.TotalCount, Equals, testCase.expectedCountOfAll)
	}
}

func (s *TestAgentSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestAgentSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}

func (s *TestAgentSuite) SetUpTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentSuite.TestListAgents":
		executeInTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(4990, 'CISCO 機房')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(12001, '上海光速群'),(12002, '湖南SSD群')
			`,
			`
			-- IP: 123.52.14.21
			INSERT INTO nqm_agent(
				ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status,
				ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
			)
			VALUES(7061, 'ag-name-1', 'ag-list-1', 'hn-list-1', x'7B340E15', 1, 3, 3, 5, 4990)
			`,
			`
			INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
			VALUES(7061, 12001),(7061, 12002)
			`,
			`
			-- IP: 12.5.104.121
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_status)
			VALUES(7062, 'ag-list-2', 'hn-list-2', x'0C056879', 4, 0)
			`,
			`
			-- IP: 12.37.22.48
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_status)
			VALUES(7063, 'ag-list-3', 'hn-list-3', x'0C251630', 3, 1)
			`,
		)
	}
}
func (s *TestAgentSuite) TearDownTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentSuite.TestListAgents":
		executeInTx(
			`
			DELETE FROM nqm_agent
			WHERE ag_id >= 7061 AND ag_id <= 7063
			`,
			`
			DELETE FROM owl_name_tag
			WHERE nt_id = 4990
			`,
			`
			DELETE FROM owl_group_tag
			WHERE gt_id >= 12001 AND gt_id <= 12002
			`,
		)
	}
}
