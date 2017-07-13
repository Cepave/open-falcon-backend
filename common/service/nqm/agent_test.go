package nqm

import (
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"time"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

var testedAgentService = NewAgentService(
	cache.DataCacheConfig{
		MaxSize:  10,
		Duration: time.Minute * 5,
	},
)

// Tests the getting of simple agent by id
func (suite *TestAgentSuite) TestGetSimpleAgent1ById(c *C) {
	testCases := []*struct {
		sampleId int32
		hasFound bool
	}{
		{40571, true},
		{40572, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		/**
		 * Asserts the found data
		 */
		testedResult := testedAgentService.GetSimpleAgent1ById(testCase.sampleId)
		c.Assert(testedResult, ocheck.ViableValue, testCase.hasFound, comment)
		// :~)

		/**
		 * Asserts the cache
		 */
		testedCache := testedAgentService.cache.Get(getKeyByAgentId(testCase.sampleId))
		c.Assert(testedCache, ocheck.ViableValue, testCase.hasFound, comment)
		// :~)
	}

}

// Tests the loading of SimpleAgent1 by filter
func (suite *TestAgentSuite) TestGetSimpleAgent1sByFilter(c *C) {
	testCases := []*struct {
		sampleFilter   *nqmModel.AgentFilter
		expectedCache  []int32
		expectedNumber int
	}{
		{
			&nqmModel.AgentFilter{
				Name: []string{"no-such-1"},
			},
			[]int32{},
			0,
		},
		{
			&nqmModel.AgentFilter{
				Name: []string{"ag-tg-1", "ag-tg-2"},
			},
			[]int32{75061, 75062},
			2,
		},
		{
			&nqmModel.AgentFilter{},
			[]int32{75061, 75062, 75063},
			3,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedAgentService.GetSimpleAgent1sByFilter(testCase.sampleFilter)
		c.Assert(testedResult, HasLen, testCase.expectedNumber, comment)

		for _, id := range testCase.expectedCache {
			comment = Commentf("%s. Needs cache by agent id: [%d]", comment.CheckCommentString(), id)

			testedCache := testedAgentService.cache.Get(
				getKeyByAgentId(id),
			)

			c.Assert(testedCache, ocheck.ViableValue, true, comment)
		}

		testedAgentService.cache.Clear()
	}
}

func (s *TestAgentSuite) SetUpTest(c *C) {
	var inTx = nqmDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentSuite.TestGetSimpleAgent1sByFilter":
		inTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(6071, 'load-tk-1', '', ''),
				(6072, 'load-tk-2', '', ''),
				(6073, 'load-tk-3', '', '')
			`,
			`
			INSERT INTO nqm_agent(
				ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status,
				ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
			)
			VALUES(75061, 6071, 'ag-tg-1-C01', 'ag-yk-1@201.3.116.1', 'ag-yk-1-C01', x'C9037401', 1, -1, -1, -1, -1),
				(75062, 6072, 'ag-tg-2-C01', 'ag-yk-2@201.3.116.2', 'ag-yk-2-C01', x'C9037402', 1, -1, -1, -1, -1),
				(75063, 6073, 'ag-tg-3', 'ag-yk-3@201.4.23.3', 'ag-yk-3', x'C9041703', 1, -1, -1, -1, -1)
			`,
		)
	case "TestAgentSuite.TestGetSimpleAgent1ById":
		inTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(44091, 'simple-test-1', '', '')
			`,
			`
			INSERT INTO nqm_agent(
				ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status,
				ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
			)
			VALUES(40571, 44091, 'ag-name-1', 'simple-get-1@187.93.16.55', 'ag-get-1.nohh.com', x'375A1637', 1, 3, 3, 5, -1)
			`,
		)
	}
}
func (s *TestAgentSuite) TearDownTest(c *C) {
	var inTx = nqmDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentSuite.TestGetSimpleAgent1sByFilter":
		inTx(
			`DELETE FROM nqm_agent WHERE ag_id >= 75061 AND ag_id <= 75063`,
			`DELETE FROM host WHERE id >= 6071 AND id <= 6073`,
		)
	case "TestAgentSuite.TestGetSimpleAgent1ById":
		inTx(
			"DELETE FROM nqm_agent WHERE ag_id = 40571",
			"DELETE FROM host WHERE id = 44091",
		)
	}
}

func (s *TestAgentSuite) SetUpSuite(c *C) {
	nqmDb.DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestAgentSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, nqmDb.DbFacade)
}
