package nqm

import (
	"database/sql"
	"net"
	"sort"
	"time"

	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"

	. "gopkg.in/check.v1"
)

type TestDbNqmSuite struct{}

var _ = Suite(&TestDbNqmSuite{})

func (s *TestDbNqmSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestDbNqmSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
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
	var testedAgent = nqmModel.NewNqmAgent(
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

	DbFacade.SqlDbCtrl.QueryForRow(
		commonDb.RowCallbackFunc(func(row *sql.Row) {
			row.Scan(&testedConnectionId, &testedHostName, &testedIpAddress, &testedLenOfIpAddress)
		}),
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

	for _, testCase := range testedCases {
		testedTargets, err := GetTargetsByAgentForRpc(testCase.agentId)

		c.Assert(err, IsNil)
		c.Assert(len(testedTargets), Equals, len(testCase.expectedIdOfTargets))

		sort.Sort(byId(testedTargets))

		/**
		 * Asserts the matching for concise id of targets
		 */
		for i, target := range testedTargets {
			c.Assert(target.Id, Equals, testCase.expectedIdOfTargets[i])

			switch target.Id {
			case 402001:
				c.Assert(target.GroupTagIds, IsNil)
			case 402002:
				c.Assert(target.GroupTagIds, DeepEquals, []int32 { 12021, 12022, 12023 })
			case 402003:
				c.Assert(target.GroupTagIds, DeepEquals, []int32 { 12023, 12024 })
			default:
				c.Fatalf("Unknown id of target: [%v]", target.Id)
			}
		}
		// :~)
	}
}

/**
 * Tests getting data of agent for RPC
 */
type getAndRefreshNeedPingAgentTestCase struct {
	agentId int
	checkTimeAsString string
	expectedUpdatedPingTask int

	testedAgent *commonModel.NqmAgent
	checkTimeAsTime time.Time
	testedErr error
}

func (suite *TestDbNqmSuite) TestGetAndRefreshNeedPingAgentForRpc(c *C) {
	testedCases := []getAndRefreshNeedPingAgentTestCase{
		{
			agentId: 130001, checkTimeAsString: "2010-05-05T11:00:00+08:00",
			expectedUpdatedPingTask: 3,
		},
		{
			agentId: 130002, checkTimeAsString: "2010-05-05T11:00:00+08:00",
			expectedUpdatedPingTask: 0,
		},
	}

	for _, testCase := range testedCases {
		testCase.checkTimeAsTime, _ = time.Parse(time.RFC3339, testCase.checkTimeAsString)

		testCase.testedAgent, testCase.testedErr = GetAndRefreshNeedPingAgentForRpc(
			testCase.agentId, testCase.checkTimeAsTime,
		)

		assertRefreshedPingTask(c, &testCase);
	}
}
func assertRefreshedPingTask(c *C, testCase *getAndRefreshNeedPingAgentTestCase) {
	c.Assert(testCase.testedErr, IsNil)

	/**
	 * Asserts the number of modified time of last executed
	 */
	var numberOfModified int = -1

	DbFacade.SqlDbCtrl.QueryForRow(
		commonDb.RowCallbackFunc(func(row *sql.Row) {
			row.Scan(&numberOfModified)
		}),
		`
		SELECT COUNT(*)
		FROM nqm_agent_ping_task
		WHERE apt_ag_id = ?
			AND apt_time_last_execute = FROM_UNIXTIME(?)
		`,
		testCase.agentId, testCase.checkTimeAsTime.Unix(),
	)

	c.Assert(numberOfModified, Equals, testCase.expectedUpdatedPingTask)
	// :~)

	/**
	 * Asserts the result data of agent
	 */
	if testCase.expectedUpdatedPingTask > 0 {
		agentData := testCase.testedAgent

		c.Assert(agentData.IspId, Equals, int16(3))
		c.Assert(agentData.ProvinceId, Equals, commonModel.UNDEFINED_PROVINCE_ID)
		c.Assert(agentData.CityId, Equals, commonModel.UNDEFINED_CITY_ID)
		c.Assert(agentData.ProvinceId, Equals, commonModel.UNDEFINED_CITY_ID)
		c.Assert(agentData.NameTagId, Equals, commonModel.UNDEFINED_NAME_TAG_ID)
		c.Assert(agentData.GroupTagIds, DeepEquals, []int32 { 9931, 9932, 9933 })
	}
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
		{2003, HAS_PING_TASK}, // The agent has ping task(enabled, with ISP filter)
		{2004, HAS_PING_TASK}, // The agent has ping task(enabled, with province filter)
		{2005, HAS_PING_TASK}, // The agent has ping task(enabled, with city filter)
		{2006, HAS_PING_TASK}, // The agent has ping task(enabled, with name tag filter)
		{2007, HAS_PING_TASK}, // The agent has ping task(enabled, with group tag filter)
		{2010, HAS_PING_TASK_MATCH_ANY_TARGET}, // The agent has ping task(enabled, without filters)
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
		expectedNumberOfIspFilters int
		expectedNumberOfProvinceFilters int
		expectedNumberOfCityFilters int
		expectedNumberOfNameTagFilters int
		expectedNumberOfGroupTagFilters int
	} {
		{ // Tests the trigger of insertion for filters
			[]string {
				`INSERT INTO nqm_pt_target_filter_name_tag(tfnt_pt_id, tfnt_nt_id) VALUES(9201, 3071), (9201, 3072)`,
				`INSERT INTO nqm_pt_target_filter_isp(tfisp_pt_id, tfisp_isp_id) VALUES(9201, 2), (9201, 3)`,
				`INSERT INTO nqm_pt_target_filter_province(tfpv_pt_id, tfpv_pv_id) VALUES(9201, 6), (9201, 7)`,
				`INSERT INTO nqm_pt_target_filter_city(tfct_pt_id, tfct_ct_id) VALUES(9201, 16), (9201, 17)`,
				`INSERT INTO nqm_pt_target_filter_group_tag(tfgt_pt_id, tfgt_gt_id) VALUES(9201, 70021), (9201, 70022)`,
			},
			2, 2, 2, 2, 2,
		},
		{ // Tests the trigger of deletion for filters
			[]string {
				`DELETE FROM nqm_pt_target_filter_name_tag WHERE tfnt_pt_id = 9201`,
				`DELETE FROM nqm_pt_target_filter_isp WHERE tfisp_pt_id = 9201`,
				`DELETE FROM nqm_pt_target_filter_province WHERE tfpv_pt_id = 9201`,
				`DELETE FROM nqm_pt_target_filter_city WHERE tfct_pt_id = 9201`,
				`DELETE FROM nqm_pt_target_filter_group_tag WHERE tfgt_pt_id = 9201`,
			},
			0, 0, 0, 0, 0,
		},
	}

	for _, testCase := range testedCases {
		/**
		 * Executes INSERT/DELETE statements
		 */
		DbFacade.SqlDbCtrl.InTx(commonDb.BuildTxForSqls(testCase.sqls...))
		// :~)

		numberOfRows := 0
		DbFacade.SqlDbCtrl.QueryForRow(
			commonDb.RowCallbackFunc(func(row *sql.Row) {
				numberOfRows++

				var numberOfIspFilters int
				var numberOfProvinceFilters int
				var numberOfCityFilters int
				var numberOfNameTagFilters int
				var numberOfGroupTagFilters int

				row.Scan(
					&numberOfIspFilters,
					&numberOfProvinceFilters,
					&numberOfCityFilters,
					&numberOfNameTagFilters,
					&numberOfGroupTagFilters,
				)

				/**
				 * Asserts the cached value for number of filters
				 */
				c.Assert(numberOfIspFilters, Equals, testCase.expectedNumberOfIspFilters);
				c.Assert(numberOfProvinceFilters, Equals, testCase.expectedNumberOfProvinceFilters);
				c.Assert(numberOfCityFilters, Equals, testCase.expectedNumberOfCityFilters);
				c.Assert(numberOfNameTagFilters, Equals, testCase.expectedNumberOfNameTagFilters);
				c.Assert(numberOfGroupTagFilters, Equals, testCase.expectedNumberOfGroupTagFilters);
				// :~)
			}),
			`
			SELECT
				pt_number_of_isp_filters,
				pt_number_of_province_filters,
				pt_number_of_city_filters,
				pt_number_of_name_tag_filters,
				pt_number_of_group_tag_filters
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
		{ 47301, 1 },
		{ 47302, 1 },
		{ 47303, 1 },
		{ 47304, 1 },
		{ 47305, 1 },
		{ 47311, 1 },
	}

	for _, testCase := range testCases {
		c.Logf("Current tested id of ping task: [%d]", testCase.pingTaskId)

		var numberOfRows int = 0
		DbFacade.SqlDbCtrl.QueryForRows(
			commonDb.RowsCallbackFunc(func (row *sql.Rows) commonDb.IterateControl {
				numberOfRows++

				var targetId int32

				row.Scan(&targetId)
				c.Logf("Current target: [%v]", targetId)

				return commonDb.IterateContinue
			}),
			`
			SELECT tg_id FROM vw_enabled_targets_by_ping_task
			WHERE tg_pt_id = ?
			`,
			testCase.pingTaskId,
		)

		c.Logf("Number of matched targets: [%d]", numberOfRows)
		c.Assert(numberOfRows, Equals, testCase.expectedNumberOfData)
	}
}

func (s *TestDbNqmSuite) SetUpTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestDbNqmSuite.Test_vw_enabled_targets_by_ping_task":
		executeInTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES (4071, 'vw-tag-1'), (4072, 'vw-tag-2')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES (23201, 'group-tag-1'),
				(23202, 'group-tag-2'),
				(23203, 'group-tag-3')
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id,
				tg_status, tg_available
			)
			VALUES
				(72001, 'tgn-e-1', '105.12.3.1', 3, -1, -1, -1, TRUE, TRUE), # Matched by ISP
				(72002, 'tgn-e-2', '105.12.3.2', -1, 6, -1, -1, TRUE, TRUE), # Matched by province
				(72003, 'tgn-e-3', '105.12.3.3', -1, -1, 12, -1, TRUE, TRUE), # Matched by city
				(72004, 'tgn-e-4', '105.12.3.4', -1, -1, -1, 4071, TRUE, TRUE), # Matched by name tag
				(72005, 'tgn-e-5', '105.12.3.5', -1, -1, -1, -1, TRUE, TRUE), # Matched by group tag
				(72011, 'tgn-e-11', '105.12.3.11', 4, 7, 13, 4072, TRUE, TRUE), # Matched with all of the filters
				(72013, 'tgn-d-1', '106.12.3.1', 4, 7, 13, 4072, TRUE, FALSE), # Matched, but disabled(status)
				(72014, 'tgn-d-2', '106.12.3.2', 4, 7, 13, 4072, FALSE, TRUE) # Matched, but disabled(available)
			`,
			`
			INSERT INTO nqm_target_group_tag(
				tgt_tg_id, tgt_gt_id
			)
			VALUES(72005, 23201),
				(72013, 23202), (72014, 23202), -- Disabled
				(72011, 23202), (72011, 23203) -- Matched with all of the filters
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_id, pt_period, pt_enable
			)
			VALUES (47301, 20, true), -- Just ISP filter
				(47302, 20, true), -- Just province filter
				(47303, 20, true), -- Just city filter
				(47304, 20, true), -- Just filter of name tag
				(47305, 20, true), -- Just filter of group tag
				(47306, 20, false), -- disabled filter
				(47311, 20, true) -- Has all of the filters
			`,
			`
			INSERT INTO nqm_pt_target_filter_isp(
				tfisp_pt_id, tfisp_isp_id
			)
			VALUES (47301, 3),
				(47311, 3), (47311, 4),
				(47306, 3)
			`,
			`
			INSERT INTO nqm_pt_target_filter_province(
				tfpv_pt_id, tfpv_pv_id
			)
			VALUES (47302, 6),
				(47311, 6), (47311, 7)
			`,
			`
			INSERT INTO nqm_pt_target_filter_city(
				tfct_pt_id, tfct_ct_id
			)
			VALUES (47303, 12),
				(47311, 12), (47311, 13)
			`,
			`
			INSERT INTO nqm_pt_target_filter_name_tag(
				tfnt_pt_id, tfnt_nt_id
			)
			VALUES (47304, 4071),
				(47311, 4071), (47311, 4072)
			`,
			`
			INSERT INTO nqm_pt_target_filter_group_tag(
				tfgt_pt_id, tfgt_gt_id
			)
			VALUES (47305, 23201),
				(47311, 23201),
				(47311, 23202), (47311, 23203)
			`,
		)
	case "TestDbNqmSuite.TestTriggersOfFiltersForPingTask":
		executeInTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES (3071, 'tri-tag-1'), (3072, 'tri-tag-2')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES (70021, 'gt-01'), (70022, 'gt-02')
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_period)
			VALUES (9201, 30)
			`,
		)
	case "TestDbNqmSuite.TestGetAndRefreshNeedPingAgentForRpc":
		executeInTx(
			`SET time_zone = '+08:00'`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(9931, 'blue-1'), (9932, 'blue-2'), (9933, 'blue-3')
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(40051, 'tt1.org', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_status)
			VALUES
				(130001, 40051, 'gc-1', 'tt1.org', 0x12345678, 3, TRUE), # Enabled agent(with complex situation)
				(130002, 40051, 'gc-5', 'tt5.org', 0x15345678, 3, FALSE) # The agent is disabled
			`,
			`
			INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
			VALUES(130001, 9931), (130001, 9932), (130001, 9933)
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_period, pt_enable)
			VALUES
				(9401, 60, false), # Disabled ping task
				(9402, 60, true), # The period is not elapsed
				(9403, 60, true), # The period is not elapsed
				(9404, 60, true), # Never executed
				(9405, 60, true), # Never executed
				(9406, 60, true) # The period is elapsed
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id, apt_time_last_execute)
			VALUES
				/**
				 * Enabled agent
				 */
				(130001, 9401, '2010-05-05 08:00:00'),
				(130001, 9402, '2010-05-05 10:01:00'),
				(130001, 9403, '2010-05-05 10:13:00'),
				(130001, 9404, NULL),
				(130001, 9405, NULL),
				(130001, 9406, '2010-05-05 09:58:00'),
				# :~)
				/**
				 * 1. The agent is disabled
				 * 2. Two of the ping task should be executed if the agent is enabled
				 */
				(130002, 9404, NULL),
				(130002, 9405, '2012-05-05 09:58:00')
				# :~)
			`,
		)
	case "TestDbNqmSuite.TestGetPingTaskState":
		executeInTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(9031, 'nt-1')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(20051, 'gt-1')
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(32021, 'aaa1.ccc', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(2001, 32021, 'pt-01', 'aaa1.ccc', 0x12345671), # The agent has no ping task
				(2002, 32021, 'pt-02', 'aaa2.ccc', 0x12345672), # The agent has ping task, which are disabled
				(2003, 32021, 'pt-03', 'aaa3.ccc', 0x12345673), # The agent has ping task with filter(isp)
				(2004, 32021, 'pt-04', 'aaa4.ccc', 0x12345674), # The agent has ping task with filter(province)
				(2005, 32021, 'pt-05', 'aaa5.ccc', 0x12345675), # The agent has ping task with filter(city)
				(2006, 32021, 'pt-06', 'aaa6.ccc', 0x12345676), # The agent has ping task with filter(name tag)
				(2007, 32021, 'pt-07', 'aaa7.ccc', 0x12345677), # The agent has ping task with filter(group tag)
				(2010, 32021, 'pt-10', 'aaa10.ccc', 0x14345679) # The agent has ping task without filter
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_id, pt_period, pt_enable
			)
			VALUES
				(7001, 20, false),
				(7002, 20, false),
				(7003, 20, true), # With ISP filter
				(7004, 20, true), # With province filter
				(7005, 20, true), # With city filter
				(7006, 20, true), # With name tag filter
				(7007, 20, true), # With group tag filter
				(7010, 20, true)
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES
				(2002, 7001),
				(2002, 7002),
				(2003, 7003),
				(2004, 7004),
				(2005, 7005),
				(2006, 7006),
				(2007, 7007),
				(2010, 7010)
			`,
			`
			INSERT INTO nqm_pt_target_filter_isp(tfisp_pt_id, tfisp_isp_id)
			VALUES(7003, 3)
			`,
			`
			INSERT INTO nqm_pt_target_filter_province(tfpv_pt_id, tfpv_pv_id)
			VALUES(7004, 2)
			`,
			`
			INSERT INTO nqm_pt_target_filter_city(tfct_pt_id, tfct_ct_id)
			VALUES(7005, 12)
			`,
			`
			INSERT INTO nqm_pt_target_filter_name_tag(tfnt_pt_id, tfnt_nt_id)
			VALUES(7006, 9031)
			`,
			`
			INSERT INTO nqm_pt_target_filter_group_tag(tfgt_pt_id, tfgt_gt_id)
			VALUES(7007, 20051)
			`,
		)
	case "TestDbNqmSuite.TestGetTargetsByAgentForRpc":
		executeInTx(
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(12021, 'bmw-1'), (12022, 'bmw-2'), (12023, 'bmw-3'), (12024, 'bmw-4')
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(40091, 'tl-01', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(230001, 40091, 'tl-01', 'ccb1.ccc', 0x12345678),
				(230002, 40091, 'tl-02', 'ccb2.ccc', 0x22345678),
				(230003, 40091, 'tl-03', 'ccb3.ccc', 0x32345678)
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_probed_by_all, tg_nt_id,
				tg_status, tg_available
			)
			VALUES
				# group tags: <none>
				(402001, 'tgn-1', '1.2.3.4', -1, -1, -1, true, -1, true, true), # Probed by all
				# group tags: 12021, 12022, 12023
				(402002, 'tgn-2', '1.2.3.5', 5, -1, -1, false, -1, true, true),
				# group tags: 12023, 12024
				(402003, 'tgn-3', '1.2.3.6', -1, -1, -1, false, -1, true, true),
				/**
				 * Disabled target
				 */
				(402005, 'tgn-4', '1.2.3.11', 5, -1, -1, true, -1, false, true),
				(402006, 'tgn-5', '1.2.3.12', 5, -1, -1, true, -1, true, false)
				# :~)
			`,
			`
			INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
			VALUES(402002, 12021), (402002, 12022), (402002, 12023),
				(402003, 12023), (402003, 12024)
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
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestDbNqmSuite.Test_vw_enabled_targets_by_ping_task":
		executeInTx(
			`DELETE FROM nqm_ping_task WHERE pt_id >= 47301 AND pt_id <= 47311`,
			`DELETE FROM nqm_target WHERE tg_id >= 72001 AND tg_id <= 72014`,
			`DELETE FROM owl_name_tag WHERE nt_id >= 4071 AND nt_id <= 4072`,
			`DELETE FROM nqm_target_group_tag WHERE tgt_tg_id >= 72001 AND tgt_tg_id <= 72014`,
			`DELETE FROM owl_group_tag WHERE gt_id >= 23201 AND gt_id <= 23203`,
		)
	case "TestDbNqmSuite.TestTriggersOfFiltersForPingTask":
		executeInTx(
			`DELETE FROM nqm_ping_task WHERE pt_id = 9201`,
			`DELETE FROM owl_name_tag WHERE nt_id >= 3071 AND nt_id <= 3072`,
			`DELETE FROM owl_group_tag WHERE gt_id >= 70021 AND gt_id <= 70022`,
		)
	case "TestDbNqmSuite.TestRefreshAgentInfo":
		executeInTx(
			"DELETE FROM nqm_agent WHERE ag_connection_id = 'refresh-1'",
			"DELETE FROM host WHERE hostname = 'refresh1.com'",
		)
	case "TestDbNqmSuite.TestGetAndRefreshNeedPingAgentForRpc":
		executeInTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 130001 AND apt_ag_id <= 130005",
			"DELETE FROM nqm_ping_task WHERE pt_id >= 9401 AND pt_id <= 9410",
			"DELETE FROM nqm_agent_group_tag WHERE agt_ag_id >= 130001 AND agt_ag_id <= 130005",
			"DELETE FROM nqm_agent WHERE ag_id >= 130001 AND ag_id <= 130005",
			"DELETE FROM host WHERE id = 40051",
			"DELETE FROM owl_group_tag WHERE gt_id >= 9931 AND gt_id <= 9933",
		)
	case "TestDbNqmSuite.TestGetPingTaskState":
		executeInTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 2001 AND apt_ag_id <= 2010",
			"DELETE FROM nqm_ping_task WHERE pt_id >= 7001 AND pt_id <= 7010",
			"DELETE FROM nqm_agent WHERE ag_id >= 2001 AND ag_id <= 2010",
			"DELETE FROM host WHERE id = 32021",
			"DELETE FROM owl_name_tag WHERE nt_id = 9031",
			"DELETE FROM owl_group_tag WHERE gt_id = 20051",
		)
	case "TestDbNqmSuite.TestGetTargetsByAgentForRpc":
		executeInTx(
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 230001 AND apt_ag_id <= 230003",
			"DELETE FROM nqm_ping_task WHERE pt_id >= 34021 AND pt_id <= 34023",
			"DELETE FROM nqm_target_group_tag WHERE tgt_tg_id >= 402001 AND tgt_tg_id <= 402010",
			"DELETE FROM nqm_agent WHERE ag_id >= 230001 AND ag_id <= 230003",
			"DELETE FROM host WHERE id = 40091",
			"DELETE FROM nqm_target WHERE tg_id >= 402001 AND tg_id <= 402010",
			"DELETE FROM owl_group_tag WHERe gt_id >= 12021 AND gt_id <= 12024",
		)
	}
}
