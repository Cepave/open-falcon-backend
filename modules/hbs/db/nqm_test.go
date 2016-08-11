package db

import (
	"database/sql"
	"net"
	"sort"
	"time"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/hbs/model"
	hbstesting "github.com/Cepave/open-falcon-backend/modules/hbs/testing"
	. "gopkg.in/check.v1"
)

type TestDbNqmSuite struct{}

var _ = Suite(&TestDbNqmSuite{})

func (s *TestDbNqmSuite) SetUpSuite(c *C) {
	(&TestDbSuite{}).SetUpSuite(c)
}

func (s *TestDbNqmSuite) TearDownSuite(c *C) {
	(&TestDbSuite{}).TearDownSuite(c)
}

/**
 * Tests the insertion and refresh for a agent
 */
type refreshAgentTestCase struct {
	connectionId string
	hostName     string
	ipAddress    string
}

func (suite *TestDbNqmSuite) TestRefreshAgentInfo(c *C) {
	var testedCases = []refreshAgentTestCase{
		{"refresh-1", "refresh1.com", "100.20.44.12"}, // First time creation of data
		{"refresh-1", "refresh2.com", "100.20.44.13"}, // Refresh of data
	}

	for _, v := range testedCases {
		testRefreshAgentInfo(c, v)
	}
}

func testRefreshAgentInfo(c *C, args refreshAgentTestCase) {
	var testedAgent = model.NewNqmAgent(
		&commonModel.NqmTaskRequest{
			ConnectionId: args.connectionId,
			Hostname:     args.hostName,
			IpAddress:    args.ipAddress,
		},
	)

	err := RefreshAgentInfo(testedAgent)

	/**
	 * Asserts the new id
	 */
	c.Assert(err, IsNil)
	c.Logf("Got agent id: %d", testedAgent.Id)
	c.Assert(testedAgent.Id > 0, Equals, true)
	// :~)

	var testedHostName string
	var testedConnectionId string
	var testedIpAddress net.IP
	var testedLenOfIpAddress int

	hbstesting.QueryForRow(
		func(row *sql.Row) {
			row.Scan(&testedConnectionId, &testedHostName, &testedIpAddress, &testedLenOfIpAddress)
		},
		"SELECT ag_connection_id, ag_hostname, ag_ip_address, BIT_LENGTH(ag_ip_address) AS len_of_ip_address FROM nqm_agent WHERE ag_id = ?",
		testedAgent.Id,
	)

	c.Logf("Ip Address: \"%s\". Length(bits): [%d]", testedIpAddress, testedLenOfIpAddress)

	/**
	 * Asserts the data on database
	 */
	c.Assert(testedConnectionId, Equals, testedAgent.ConnectionId())
	c.Assert(testedHostName, Equals, testedAgent.Hostname())
	c.Assert(testedIpAddress.Equal(testedAgent.IpAddress), Equals, true)
	c.Assert(testedLenOfIpAddress, Equals, 32)
	// :~)
}

/**
 * Tests getting targets by filter
 */
type byId []commonModel.NqmTarget

func (targets byId) Len() int           { return len(targets) }
func (targets byId) Swap(i, j int)      { targets[i], targets[j] = targets[j], targets[i] }
func (targets byId) Less(i, j int) bool { return targets[i].Id < targets[j].Id }

func (suite *TestDbNqmSuite) TestGetTargetsByAgentForRpc(c *C) {
	testedCases := []struct {
		agentId             int
		expectedIdOfTargets []int
	}{
		{230001, []int{ 402001, 402002, 402003 }}, // All of the targets
		{230002, []int{ 402001, 402002 }}, // Targets are matched by ISP(other matchings are tested on vw_enabled_targets_by_ping_task)
		{230003, []int{ 402001 }}, // Nothing matched except probed by all
	}

	for _, v := range testedCases {
		testedTargets, err := GetTargetsByAgentForRpc(v.agentId)

		c.Assert(err, IsNil)

		c.Assert(len(testedTargets), Equals, len(v.expectedIdOfTargets))

		sort.Sort(byId(testedTargets))

		/**
		 * Asserts the matching for concise id of targets
		 */
		for i, target := range testedTargets {
			c.Assert(target.Id, Equals, v.expectedIdOfTargets[i])
		}
		// :~)
	}
}

/**
 * Tests getting data of agent for RPC
 */
type getAndRefreshNeedPingAgentTestCase struct {
	agentId   int
	checkTimeAsString string
	checker   Checker

	checkTimeAsTime time.Time
	testedAgent *commonModel.NqmAgent
	testedErr error
}

func (suite *TestDbNqmSuite) TestGetAndRefreshNeedPingAgentForRpc(c *C) {
	testedCases := []getAndRefreshNeedPingAgentTestCase{
		{agentId: 130001, checkTimeAsString: "2115-08-08T00:00:00+00:00", checker: IsNil},  // No ping task setting
		{agentId: 130002, checkTimeAsString: "2010-05-05T10:59:00+08:00", checker: IsNil},  // The period is not elapsed yet
		{agentId: 130003, checkTimeAsString: "2013-10-01T00:00:00+08:00", checker: NotNil}, // Never executed
		{agentId: 130004, checkTimeAsString: "2012-06-10T09:00:00+08:00", checker: NotNil}, // The period is elapsed
		{agentId: 130005, checkTimeAsString: "2012-06-10T09:00:00+08:00", checker: IsNil},  // Disabled agent
	}

	for _, testCase := range testedCases {
		testCase.checkTimeAsTime, _ = time.Parse(time.RFC3339, testCase.checkTimeAsString)

		testCase.testedAgent, testCase.testedErr = GetAndRefreshNeedPingAgentForRpc(
			testCase.agentId, testCase.checkTimeAsTime,
		)

		assertNeedPingAgent(c, testCase)
	}
}

func assertNeedPingAgent(
	c *C, testCase getAndRefreshNeedPingAgentTestCase,
) {
	c.Logf("Current tested agent Id: %d", testCase.agentId)

	testedAgent := testCase.testedAgent

	/**
	 * Asserts the result data
	 */
	c.Assert(testCase.testedErr, IsNil)
	c.Assert(testedAgent, testCase.checker)
	// :~)

	if testCase.checker == IsNil {
		return
	}

	/**
	 * Asserts the content of returned agent
	 */
	c.Assert(testedAgent.Id, Equals, testCase.agentId)
	c.Assert(testedAgent.IspId, Equals, int16(3))
	c.Assert(testedAgent.ProvinceId, Equals, commonModel.UNDEFINED_PROVINCE_ID)
	c.Assert(testedAgent.CityId, Equals, commonModel.UNDEFINED_CITY_ID)
	c.Assert(testedAgent.NameTagId, Equals, commonModel.UNDEFINED_NAME_TAG_ID)
	// :~)

	/**
	 * Asserts the updated time of ping task
	 */
	//var updatedTime int64
	var unixTime int64
	hbstesting.QueryForRow(
		func(row *sql.Row) {
			c.Assert(row.Scan(&unixTime), IsNil)
		},
		`
		SELECT UNIX_TIMESTAMP(apt_time_last_execute)
		FROM nqm_agent_ping_task
		WHERE apt_ag_id = ?
		`,
		testCase.agentId,
	)

	c.Assert(unixTime, Equals, testCase.checkTimeAsTime.Unix())
	// :~)
}

/**
 * Tests the state of ping task
 */
func (suite *TestDbNqmSuite) TestGetPingTaskState(c *C) {
	testedCases := []struct {
		agentId        int
		expectedStatus int
	} {
		{2001, NO_PING_TASK}, // The agent has no ping task
		{2002, NO_PING_TASK}, // The agent has ping task, which are disabled
		{2003, HAS_PING_TASK}, // The agent has ping task(enabled, with filters)
		{2004, HAS_PING_TASK_MATCH_ANY_TARGET}, // The agent has ping task(enabled, without filters)
	}

	for _, v := range testedCases {
		testedResult, err := getPingTaskState(v.agentId)

		c.Assert(err, IsNil)
		c.Assert(testedResult, Equals, v.expectedStatus)
	}
}

/**
 * Tests the triggers for filters of PING TASK
 */
func (suite *TestDbNqmSuite) TestTriggersOfFiltersForPingTask(c *C) {
	testedCases := []struct {
		sqls []string
		expectedNumberOfNameTagFilters int
		expectedNumberOfIspFilters int
		expectedNumberOfProvinceFilters int
		expectedNumberOfCityFilters int
	} {
		{ // Tests the trigger of insertion for filters
			[]string {
				`INSERT INTO nqm_pt_target_filter_name_tag(tfnt_pt_id, tfnt_nt_id) VALUES(9201, 3071), (9201, 3072)`,
				`INSERT INTO nqm_pt_target_filter_isp(tfisp_pt_id, tfisp_isp_id) VALUES(9201, 2), (9201, 3)`,
				`INSERT INTO nqm_pt_target_filter_province(tfpv_pt_id, tfpv_pv_id) VALUES(9201, 6), (9201, 7)`,
				`INSERT INTO nqm_pt_target_filter_city(tfct_pt_id, tfct_ct_id) VALUES(9201, 16), (9201, 17)`,
			},
			2, 2, 2, 2,
		},
		{ // Tests the trigger of deletion for filters
			[]string {
				`DELETE FROM nqm_pt_target_filter_name_tag WHERE tfnt_pt_id = 9201`,
				`DELETE FROM nqm_pt_target_filter_isp WHERE tfisp_pt_id = 9201`,
				`DELETE FROM nqm_pt_target_filter_province WHERE tfpv_pt_id = 9201`,
				`DELETE FROM nqm_pt_target_filter_city WHERE tfct_pt_id = 9201`,
			},
			0, 0, 0, 0,
		},
	}

	for _, testCase := range testedCases {
		/**
		 * Executes INSERT/DELETE statements
		 */
		hbstesting.ExecuteQueriesOrFailInTx(
			testCase.sqls...,
		)
		// :~)

		numberOfRows := 0
		hbstesting.QueryForRow(
			func(row *sql.Row) {
				numberOfRows++

				var numberOfNameTagFilters int
				var numberOfIspFilters int
				var numberOfProvinceFilters int
				var numberOfCityFilters int

				row.Scan(
					&numberOfNameTagFilters,
					&numberOfIspFilters,
					&numberOfProvinceFilters,
					&numberOfCityFilters,
				)

				/**
				 * Asserts the cached value for number of filters
				 */
				c.Assert(numberOfNameTagFilters, Equals, testCase.expectedNumberOfNameTagFilters);
				c.Assert(numberOfIspFilters, Equals, testCase.expectedNumberOfIspFilters);
				c.Assert(numberOfProvinceFilters, Equals, testCase.expectedNumberOfProvinceFilters);
				c.Assert(numberOfCityFilters, Equals, testCase.expectedNumberOfCityFilters);
				// :~)
			},
			`
			SELECT pt_number_of_name_tag_filters,
				pt_number_of_isp_filters,
				pt_number_of_province_filters,
				pt_number_of_city_filters
			FROM nqm_ping_task WHERE pt_id = 9201
			`,
		)

		// Ensures that the row is effective
		c.Assert(numberOfRows, Equals, 1)
	}
}

func (suite *TestDbNqmSuite) Test_vw_enabled_targets_by_ping_task(c *C) {
	testCases := []struct {
		pingTaskId int
		expectedNumberOfData int
	} {
		{ 47301, 4 },
		{ 47302, 0 },
	}

	for _, testCase := range testCases {
		var numberOfRows int = 0
		hbstesting.QueryForRows(
			func (row *sql.Rows) {
				numberOfRows++
			},
			`
			SELECT * FROM vw_enabled_targets_by_ping_task
			WHERE tg_pt_id = ?
			`,
			testCase.pingTaskId,
		)

		c.Assert(numberOfRows, Equals, testCase.expectedNumberOfData)
	}
}

func (s *TestDbNqmSuite) SetUpTest(c *C) {
	if !hbstesting.HasDbEnvForMysqlOrSkip(c) {
		return
	}

	switch c.TestName() {
	case "TestDbNqmSuite.Test_vw_enabled_targets_by_ping_task":
		hbstesting.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES (4071, 'vw-tag-1'), (4072, 'vw-tag-2')
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id,
				tg_status, tg_available
			)
			VALUES
				(72001, 'tgn-e-1', '105.12.3.1', 3, -1, -1, -1, true, true), # Matched by ISP
				(72002, 'tgn-e-2', '105.12.3.2', -1, 6, -1, -1, true, true), # Matched by province
				(72003, 'tgn-e-3', '105.12.3.3', -1, -1, 12, -1, true, true), # Matched by city
				(72004, 'tgn-e-4', '105.12.3.4', -1, -1, -1, 4071, true, true), # Matched by name tag
				(72013, 'tgn-d-1', '106.12.3.1', 4, 7, 13, 4072, false, false), # Matched, but disabled
				(72014, 'tgn-d-2', '106.12.3.2', 4, 7, 13, 4072, false, false) # Matched, but disabled
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_id, pt_period, pt_enable
			)
			VALUES (47301, 20, true), (47302, 20, false)
			`,
			`
			INSERT INTO nqm_pt_target_filter_isp(
				tfisp_pt_id, tfisp_isp_id
			)
			VALUES (47301, 3), (47302, 4)
			`,
			`
			INSERT INTO nqm_pt_target_filter_province(
				tfpv_pt_id, tfpv_pv_id
			)
			VALUES (47301, 6), (47302, 7)
			`,
			`
			INSERT INTO nqm_pt_target_filter_city(
				tfct_pt_id, tfct_ct_id
			)
			VALUES (47301, 12), (47302, 13)
			`,
			`
			INSERT INTO nqm_pt_target_filter_name_tag(
				tfnt_pt_id, tfnt_nt_id
			)
			VALUES (47301, 4071), (47302, 4072)
			`,
		)
	case "TestDbNqmSuite.TestTriggersOfFiltersForPingTask":
		hbstesting.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES (3071, 'tri-tag-1'), (3072, 'tri-tag-2')
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_period)
			VALUES (9201, 30)
			`,
		)
	case "TestDbNqmSuite.TestGetAndRefreshNeedPingAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			`SET time_zone = '+08:00'`,
			`
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_status)
			VALUES
				(130001, 'gc-1', 'tt1.org', 0x12345678, 3, TRUE), # No Ping task setting
				(130002, 'gc-2', 'tt2.org', 0x13345678, 3, TRUE), # The period is not elapse yet
				(130003, 'gc-3', 'tt3.org', 0x14345678, 3, TRUE), # Never executed
				(130004, 'gc-4', 'tt4.org', 0x15345678, 3, TRUE), # The period is elapsed
				(130005, 'gc-5', 'tt5.org', 0x15345678, 3, FALSE) # The agent is disabled
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_period)
			VALUES
				(9402, 60),
				(9403, 60),
				(9404, 60),
				(9405, 60)
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id, apt_time_last_execute)
			VALUES
				(130002, 9402, '2010-05-05 10:00:00'), # The period is not elapse yet
				(130003, 9403, NULL),
				(130004, 9404, '2012-06-10 08:00:00'), # The period is elapsed
				(130005, 9405, NULL)
			`,
		)
	case "TestDbNqmSuite.TestGetPingTaskState":
		hbstesting.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(2001, 'pt-01', 'aaa1.ccc', 0x12345678), # The agent has no ping task
				(2002, 'pt-02', 'aaa2.ccc', 0x13345678), # The agent has ping task, which are disabled
				(2003, 'pt-03', 'aaa3.ccc', 0x14345678), # The agent has ping task with filters
				(2004, 'pt-04', 'aaa4.ccc', 0x14345678) # The agent has ping task without filter
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_id, pt_period, pt_enable
			)
			VALUES
				(7001, 20, false),
				(7002, 20, false),
				(7003, 20, true),
				(7004, 20, true)
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES
				(2002, 7001),
				(2002, 7002),
				(2003, 7003),
				(2004, 7004)
			`,
			`
			INSERT INTO nqm_pt_target_filter_province(tfpv_pt_id, tfpv_pv_id)
			VALUES(7003, 2)
			`,
		)
	case "TestDbNqmSuite.TestGetTargetsByAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(230001, 'tl-01', 'ccb1.ccc', 0x12345678),
				(230002, 'tl-02', 'ccb2.ccc', 0x22345678),
				(230003, 'tl-03', 'ccb3.ccc', 0x32345678)
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_probed_by_all, tg_nt_id,
				tg_status, tg_available
			)
			VALUES
				(402001, 'tgn-1', '1.2.3.4', -1, -1, -1, true, -1, true, true), # Probed by all
				(402002, 'tgn-2', '1.2.3.5', 5, -1, -1, false, -1, true, true),
				(402003, 'tgn-3', '1.2.3.6', -1, -1, -1, false, -1, true, true)
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_id, pt_period
			)
			VALUES
				(34021, 20), # All of the targets
				(34022, 20), # Has ISP filter
				(34023, 20) # Match none except probed by all
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES (230001, 34021), (230002, 34022), (230003, 34023)
			`,
			`
			INSERT INTO nqm_pt_target_filter_isp(
				tfisp_pt_id, tfisp_isp_id
			)
			VALUES (34022, 5), (34023, 6)
			`,
		)
	}
}

func (s *TestDbNqmSuite) TearDownTest(c *C) {
	switch c.TestName() {
	case "TestDbNqmSuite.Test_vw_enabled_targets_by_ping_task":
		hbstesting.ExecuteQueriesOrFailInTx(
			`DELETE FROM nqm_ping_task WHERE pt_id >= 47301 AND pt_id <= 47302`,
			`DELETE FROM nqm_target WHERE tg_id >= 72001 AND tg_id <= 72014`,
			`DELETE FROM owl_name_tag WHERE nt_id >= 4071 AND nt_id <= 4072`,
		)
	case "TestDbNqmSuite.TestTriggersOfFiltersForPingTask":
		hbstesting.ExecuteQueriesOrFailInTx(
			`DELETE FROM nqm_ping_task WHERE pt_id = 9201`,
			`DELETE FROM owl_name_tag WHERE nt_id >= 3071 AND nt_id <= 3072`,
		)
	case "TestDbNqmSuite.TestRefreshAgentInfo":
		hbstesting.ExecuteOrFail(
			"DELETE FROM nqm_agent WHERE ag_connection_id = 'refresh-1'",
		)
	case "TestDbNqmSuite.TestGetAndRefreshNeedPingAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 130001 AND apt_ag_id <= 130005",
			"DELETE FROM nqm_ping_task WHERE pt_id >= 9402 AND pt_id <= 9405",
			"DELETE FROM nqm_agent WHERE ag_id >= 130001 AND ag_id <= 130005",
		)
	case "TestDbNqmSuite.TestGetPingTaskState":
		hbstesting.ExecuteQueriesOrFailInTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 2001 AND apt_ag_id <= 2004",
			"DELETE FROM nqm_ping_task WHERE pt_id >= 7001 AND pt_id <= 7004",
			"DELETE FROM nqm_agent WHERE ag_id >= 2001 AND ag_id <= 2004",
		)
	case "TestDbNqmSuite.TestGetTargetsByAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 230001 AND apt_ag_id <= 230003",
			"DELETE FROM nqm_ping_task WHERE pt_id >= 34021 AND pt_id <= 34023",
			"DELETE FROM nqm_agent WHERE ag_id >= 230001 AND ag_id <= 230003",
			"DELETE FROM nqm_target WHERE tg_id >= 402001 AND tg_id <= 402003",
		)
	}
}
