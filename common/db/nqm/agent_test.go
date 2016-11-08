package nqm

import (
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"net"
	"reflect"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

// Tests the getting of agent by id
func (suite *TestAgentSuite) TestGetAgentById(c *C) {
	testCases := []struct {
		sampleIdOfAgent int32
		hasFound bool
	} {
		{ 88971, true },
		{ 88972, false },
	}

	for _, testCase := range testCases {
		result := GetAgentById(testCase.sampleIdOfAgent)

		if testCase.hasFound {
			c.Logf("Found agent by id: %v", result)
			c.Assert(result, NotNil)
		} else {
			c.Assert(result, IsNil)
		}
	}
}

// Tests the adding of new agent
func (suite *TestAgentSuite) TestAddAgent(c *C) {
	/**
	 * sample agents
	 */
	defaultAgent_1 := nqmModel.NewAgentForAdding()
	defaultAgent_1.ConnectionId = "def-agent-1"
	defaultAgent_1.Hostname = "hs-def-agent-1"
	defaultAgent_1.IpAddress = net.ParseIP("0.0.0.0")

	defaultAgent_2 := nqmModel.NewAgentForAdding()
	defaultAgent_2.ConnectionId = "def-agent-2"
	defaultAgent_2.Hostname = "hs-def-agent-2"
	defaultAgent_2.IpAddress = net.ParseIP("33.29.111.10")
	defaultAgent_2.Status = false
	defaultAgent_2.Name = "sample-agent"
	defaultAgent_2.Comment = "This is sample agent"
	defaultAgent_2.IspId = 3
	defaultAgent_2.ProvinceId = 20
	defaultAgent_2.CityId = 6
	defaultAgent_2.NameTagValue = "CISCO-617"
	defaultAgent_2.GroupTags = []string {
		"TPE-03", "TPE-04", "TPE-05",
	}

	defaultAgent_3 := *defaultAgent_2
	defaultAgent_3.CityId = 50
	// :~)

	testCases := []struct {
		addedAgent *nqmModel.AgentForAdding
		hasError bool
		errorType reflect.Type
	} {
		{ defaultAgent_1, false, nil }, // Use the minimum value
		{ defaultAgent_1, true, reflect.TypeOf(ErrDuplicatedNqmAgent{}) },
		{ defaultAgent_2, false, nil }, // Use every properties
		{ &defaultAgent_3, true, reflect.TypeOf(owlDb.ErrNotInSameHierarchy{}) }, // Duplicated connection id
	}

	for _, testCase := range testCases {
		currentAddedAgent := testCase.addedAgent
		newAgent, err := AddAgent(currentAddedAgent)

		/**
		 * Asserts the occuring error
		 */
		if testCase.hasError {
			c.Assert(newAgent, IsNil)
			c.Assert(err, NotNil)
			c.Logf("Has error: %v", err)

			c.Assert(reflect.TypeOf(err), Equals, testCase.errorType)
			continue
		}
		// :~)

		c.Assert(err, IsNil)
		c.Logf("New Agent: %v", newAgent)
		c.Logf("New Agent[Group Tags]: %v", newAgent.GroupTags)

		c.Assert(newAgent.Name, Equals, currentAddedAgent.Name)
		c.Assert(newAgent.ConnectionId, Equals, currentAddedAgent.ConnectionId)
		c.Assert(newAgent.Hostname, Equals, currentAddedAgent.Hostname)
		c.Assert(newAgent.IpAddress.String(), Equals, currentAddedAgent.IpAddress.String())
		c.Assert(newAgent.IspId, Equals, currentAddedAgent.IspId)
		c.Assert(newAgent.ProvinceId, Equals, currentAddedAgent.ProvinceId)
		c.Assert(newAgent.CityId, Equals, currentAddedAgent.CityId)
		c.Assert(newAgent.NameTagId, Equals, currentAddedAgent.NameTagId)
		c.Assert(newAgent.GroupTags, HasLen, len(currentAddedAgent.GroupTags))
	}
}

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
	owlDb.DbFacade = DbFacade
}

func (s *TestAgentSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
	owlDb.DbFacade = nil
}

func (s *TestAgentSuite) SetUpTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentSuite.TestGetAgentById":
		executeInTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(12571, 'hn-get-1', '', '')
			`,
			`
			-- IP: 87.90.6.55
			INSERT INTO nqm_agent(
				ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status,
				ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
			)
			VALUES(88971, 12571, 'ag-name-1', 'ag-get-1@87.90.6.55', 'ag-get-1.nohh.com', x'575A0637', 1, 3, 3, 5, -1)
			`,
		)
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
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(67001, 'hn-list-1', '', '')
			`,
			`
			-- IP: 123.52.14.21
			INSERT INTO nqm_agent(
				ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status,
				ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
			)
			VALUES(7061, 67001, 'ag-name-1', 'ag-list-1', 'hn-list-1', x'7B340E15', 1, 3, 3, 5, 4990)
			`,
			`
			INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
			VALUES(7061, 12001),(7061, 12002)
			`,
			`
			-- IP: 12.5.104.121
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_status)
			VALUES(7062, 67001, 'ag-list-2', 'hn-list-2', x'0C056879', 4, 0)
			`,
			`
			-- IP: 12.37.22.48
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_status)
			VALUES(7063, 67001, 'ag-list-3', 'hn-list-3', x'0C251630', 3, 1)
			`,
		)
	}
}
func (s *TestAgentSuite) TearDownTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentSuite.TestGetAgentById":
		executeInTx(
			`
			DELETE FROM nqm_agent
			WHERE ag_id = 88971
			`,
			`
			DELETE FROM host
			WHERE id = 12571
			`,
		)
	case "TestAgentSuite.TestListAgents":
		executeInTx(
			`
			DELETE FROM nqm_agent
			WHERE ag_id >= 7061 AND ag_id <= 7063
			`,
			`
			DELETE FROM host WHERE id = 67001
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
	case "TestAgentSuite.TestAddAgent":
		executeInTx(
			`
			DELETE FROM nqm_agent
			WHERE ag_connection_id LIKE 'def-agent-%'
			`,
			`
			DELETE FROM host
			WHERE hostname LIKE 'hs-def-agent-%'
			`,
			`
			DELETE FROM owl_name_tag
			WHERE nt_value = 'CISCO-617'
			`,
			`
			DELETE FROM owl_group_tag
			WHERE gt_name LIKE 'TPE-%'
			`,
		)
	}
}
