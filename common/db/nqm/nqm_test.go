package nqm

import (
	"database/sql"
	"net"
	"time"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	otest "github.com/Cepave/open-falcon-backend/common/testing"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"

	. "gopkg.in/check.v1"
)

type TestDbNqmSuite struct{}

var _ = Suite(&TestDbNqmSuite{})

func (suite *TestDbNqmSuite) TestRefreshAgentInfo(c *C) {
	testCases := []*struct {
		connectionId        string
		hostname            string
		ipAddress           string
		expectedAgentDetail Checker
	}{
		{"agrh33-1@19.20", "rhs1.01.hostname", "20.98.1.31", NotNil},  // New NQM agent
		{"agrh33-2@19.20", "rhs1.02.hostname", "20.98.1.32", NotNil},  // refreshed
		{"agrh33-3@19.33", "rhs1.99.hostname", "20.98.12.101", IsNil}, // Old NQM agent(disabled)
	}

	type agentData struct {
		Hostname      string    `db:"ag_hostname"`
		IpAddress     net.IP    `db:"ag_ip_address"`
		HeartbeatTime time.Time `db:"ag_last_heartbeat"`
	}
	sampleTime := otest.ParseTime(c, "2016-07-11T10:44:01Z")
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		sampleAgent := nqmModel.NewNqmAgent(
			&commonModel.NqmTaskRequest{
				ConnectionId: testCase.connectionId,
				Hostname:     testCase.hostname,
				IpAddress:    testCase.ipAddress,
			},
		)
		agentDetail := RefreshAgentInfo(sampleAgent, sampleTime)

		// Asserts the enable/disabled status
		c.Logf("Agent Detail: %#v", agentDetail)
		c.Assert(agentDetail, testCase.expectedAgentDetail, comment)

		/**
		 * Asserts the refreshed columns
		 */
		testedAgentData := &agentData{}
		DbFacade.SqlxDbCtrl.Get(
			testedAgentData,
			`
			SELECT ag_hostname, ag_ip_address, ag_last_heartbeat
			FROM nqm_agent
			WHERE ag_connection_id = ?
			`,
			testCase.connectionId,
		)

		c.Logf("Refreshed Agent: %#v. Id: [%d]", testedAgentData, sampleAgent.Id)
		c.Assert(sampleAgent.Id, ocheck.LargerThan, 0, comment)
		c.Assert(testedAgentData.Hostname, Equals, testCase.hostname, comment)
		c.Assert(testedAgentData.IpAddress.String(), Equals, testCase.ipAddress, comment)
		c.Assert(testedAgentData.HeartbeatTime, ocheck.TimeEquals, sampleTime, comment)
		// :~)
	}
}

// Tests the getting of log of cached list of ping
func (suite *TestDbNqmSuite) TestGetCacheLogOfPingList(c *C) {
	testCases := []*struct {
		agentId int32
		checker Checker
	}{
		{80981, NotNil},
		{80982, IsNil},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		c.Assert(getCacheLogOfPingList(testCase.agentId), testCase.checker, comment)
	}
}

// Tests the building of cache on ping list
func (suite *TestDbNqmSuite) TestBuildCacheOfPingList(c *C) {
	testCases := []*struct {
		agentId                      int32
		checkedTime                  string
		expectedTimeOfAccessOnTarget int64
		expectedPeriod               []int16
	}{
		/**
		 * 1st build
		 */
		{50761, "2014-05-05T10:20:30Z", 0, []int16{-1, 20, 20, 20}},
		{50762, "2014-05-05T10:20:30Z", 1399382640, []int16{-1, 30, 30}},
		{50763, "2014-05-05T10:20:30Z", 0, []int16{-1}},
		// :~)
		/**
		 * 2nd build
		 */
		{50761, "2014-05-06T13:30:17Z", 0, []int16{-1, 20, 20, 20}},
		{50762, "2014-05-06T13:30:17Z", 1399382640, []int16{-1, 30, 30}},
		{50763, "2014-05-06T13:30:17Z", 0, []int16{-1}},
		// :~)
	}

	type targetPeriod struct {
		TargetId   int32     `db:"apl_tg_id"`
		Period     int16     `db:"apl_min_period"`
		AccessTime time.Time `db:"apl_time_access"`
	}
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		checkedTime := otest.ParseTime(c, testCase.checkedTime)
		logInDb := BuildCacheOfPingList(
			testCase.agentId, checkedTime,
		)

		expectedLen := len(testCase.expectedPeriod)

		/**
		 * Asserts the log object
		 */
		c.Logf("Log: %#v", logInDb)
		c.Assert(logInDb.RefreshTime, ocheck.TimeEquals, checkedTime, comment)
		c.Assert(logInDb.NumberOfTargets, Equals, int32(expectedLen), comment)
		// :~)

		/**
		 * Asserts the period for each target
		 */
		testedList := []*targetPeriod{}
		DbFacade.SqlxDbCtrl.Select(
			&testedList,
			`
			SELECT apl_tg_id, apl_min_period, apl_time_access
			FROM nqm_cache_agent_ping_list
			WHERE apl_apll_ag_id = ?
			ORDER BY apl_tg_id ASC
			`,
			testCase.agentId,
		)

		c.Assert(testedList, HasLen, expectedLen, comment)

		expectedTimeOfAccessOnTarget := time.Unix(testCase.expectedTimeOfAccessOnTarget, 0)
		for i, t := range testedList {
			c.Logf("Target: %#v. Access Time: %s", t, t.AccessTime)
			c.Assert(testCase.expectedPeriod[i], Equals, t.Period, comment)
			c.Assert(t.AccessTime, ocheck.TimeEquals, expectedTimeOfAccessOnTarget, comment)
		}
		// :~)
	}
}

// Tests the getting for list of ping
func (suite *TestDbNqmSuite) TestGetPingList(c *C) {
	testCases := []*struct {
		agentId        int32
		expectedResult []int32
	}{
		{70071, []int32{99021, 99022, 99023}}, // Access time is elapsed
		{70072, []int32{}},                    // Access time is not elapsed
	}

	sampleTime := otest.ParseTime(c, "2015-06-21T20:33:43Z")
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := getPingList(
			&nqmModel.NqmAgent{
				Id:        int(testCase.agentId),
				IpAddress: net.ParseIP("20.91.83.101"),
			},
			sampleTime,
		)

		/**
		 * Asserts the result data
		 */
		c.Assert(testedResult, HasLen, len(testCase.expectedResult), comment)

		for i, target := range testedResult {
			c.Logf("Target %#v", target)
			c.Assert(target.Id, Equals, int(testCase.expectedResult[i]), comment)
		}
		// :~)
	}
}

// Tests the update time of access on cache
func (suite *TestDbNqmSuite) TestUpdateAccessTime(c *C) {
	accessTime := time.Now()

	updateAccessTime(33404, accessTime)

	var accessTimeOfLog time.Time

	/**
	 * Asserts the access time of log
	 */
	DbFacade.SqlxDbCtrl.Get(
		&accessTimeOfLog,
		`
		SELECT apll_time_access
		FROM nqm_cache_agent_ping_list_log
		WHERE apll_ag_id = 33404
		`,
	)
	c.Logf("Access time of list log: %s", accessTimeOfLog)
	c.Assert(accessTimeOfLog, ocheck.TimeEquals, accessTime)
	// :~)

	/**
	 * Asserts the access time of target
	 */
	var numberOfMatchedTargets int
	DbFacade.SqlxDbCtrl.Get(
		&numberOfMatchedTargets,
		`
		SELECT COUNT(apl_tg_id)
		FROM nqm_cache_agent_ping_list
		WHERE apl_apll_ag_id = 33404
			AND apl_time_access = FROM_UNIXTIME(?)
		`,
		accessTime.Unix(),
	)
	c.Assert(numberOfMatchedTargets, Equals, 3)
	// :~)
}

// Tests the triggers for filters of PING TASK
func (suite *TestDbNqmSuite) TestTriggersOfFiltersForPingTask(c *C) {
	testedCases := []*struct {
		sqls                            []string
		expectedNumberOfIspFilters      int
		expectedNumberOfProvinceFilters int
		expectedNumberOfCityFilters     int
		expectedNumberOfNameTagFilters  int
		expectedNumberOfGroupTagFilters int
	}{
		{ // Tests the trigger of insertion for filters
			[]string{
				`INSERT INTO nqm_pt_target_filter_name_tag(tfnt_pt_id, tfnt_nt_id) VALUES(9201, 3071), (9201, 3072)`,
				`INSERT INTO nqm_pt_target_filter_isp(tfisp_pt_id, tfisp_isp_id) VALUES(9201, 2), (9201, 3)`,
				`INSERT INTO nqm_pt_target_filter_province(tfpv_pt_id, tfpv_pv_id) VALUES(9201, 6), (9201, 7)`,
				`INSERT INTO nqm_pt_target_filter_city(tfct_pt_id, tfct_ct_id) VALUES(9201, 16), (9201, 17)`,
				`INSERT INTO nqm_pt_target_filter_group_tag(tfgt_pt_id, tfgt_gt_id) VALUES(9201, 70021), (9201, 70022)`,
			},
			2, 2, 2, 2, 2,
		},
		{ // Tests the trigger of deletion for filters
			[]string{
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
				c.Assert(numberOfIspFilters, Equals, testCase.expectedNumberOfIspFilters)
				c.Assert(numberOfProvinceFilters, Equals, testCase.expectedNumberOfProvinceFilters)
				c.Assert(numberOfCityFilters, Equals, testCase.expectedNumberOfCityFilters)
				c.Assert(numberOfNameTagFilters, Equals, testCase.expectedNumberOfNameTagFilters)
				c.Assert(numberOfGroupTagFilters, Equals, testCase.expectedNumberOfGroupTagFilters)
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

// Tests the view used to load enabled targets by ping task
func (suite *TestDbNqmSuite) Test_vw_enabled_targets_by_ping_task(c *C) {
	testCases := []*struct {
		pingTaskId           int
		expectedIdsOfTargets []int32
	}{
		{47301, []int32{72001}},                                    // Matched by ISP
		{47302, []int32{72002}},                                    // Matched by province
		{47303, []int32{72003}},                                    // Matched by city
		{47304, []int32{72004}},                                    // Matched by name tag
		{47305, []int32{72005}},                                    // Matched by group tag
		{47311, []int32{72011}},                                    // Matched by all of the properties
		{47312, []int32{72001, 72002, 72003, 72004, 72005, 72011}}, // empty ping task(all of the targets)
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		targetIds := make([]int32, 0)
		DbFacade.SqlxDbCtrl.Select(
			&targetIds,
			`
			SELECT DISTINCT tg_id FROM vw_enabled_targets_by_ping_task
			WHERE tg_pt_id = ?
			ORDER BY tg_id ASC
			`,
			testCase.pingTaskId,
		)

		c.Assert(targetIds, DeepEquals, testCase.expectedIdsOfTargets, comment)
	}
}

func (s *TestDbNqmSuite) SetUpTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestDbNqmSuite.TestRefreshAgentInfo":
		executeInTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(10531, 'rhs1.99.hostname', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address, ag_status)
			VALUES (90871, 10531, 'agrh33-3@19.33', 'old-host-33-3', 0x1f340673, FALSE)
			`,
		)
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
			VALUES (47301, 11, true), -- Just ISP filter
				(47302, 12, true), -- Just province filter
				(47303, 13, true), -- Just city filter
				(47304, 14, true), -- Just filter of name tag
				(47305, 15, true), -- Just filter of group tag
				(47306, 16, false), -- disabled filter
				(47311, 17, true), -- Has all of the filters
				(47312, 18, true) -- Empty filter(includes all of the targets)
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
	case "TestDbNqmSuite.TestGetCacheLogOfPingList":
		executeInTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(45701, 'cpl-host-1', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(80981, 45701, 'cl-01', 'cpl1.ccc@19.10', 0x0A033D4B)
			`,
			`
			INSERT INTO nqm_cache_agent_ping_list_log(
				apll_ag_id, apll_number_of_targets, apll_time_access, apll_time_refresh
			)
			VALUES(80981, 20, NOW(), NOW())
			`,
		)
	case "TestDbNqmSuite.TestBuildCacheOfPingList":
		executeInTx(
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host, tg_isp_id, tg_probed_by_all, tg_available, tg_status)
			VALUES
				(80921, 'pl-tg-1', '10.81.7.1', 5, TRUE, TRUE, TRUE), -- Probed by all
				(80922, 'pl-tg-2', '10.81.7.2', 5, FALSE, TRUE, TRUE),
				(80923, 'pl-tg-3', '10.81.7.3', 5, FALSE, TRUE, TRUE),
				(80924, 'pl-tg-4', '10.81.7.4', 6, FALSE, TRUE, TRUE), -- Another ISP
				(80928, 'pl-tg-9-1', '10.81.7.5', 5, FALSE, FALSE, TRUE), -- Not-Available
				(80929, 'pl-tg-9-2', '10.81.7.6', 5, FALSE, TRUE, FALSE) -- Disabled
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_id, pt_period
			)
			VALUES
				(83051, 20), # All of the targets
				(83052, 30), # Has ISP filter
				(83053, 40) # Match none except probed by all
			`,
			`
			INSERT INTO nqm_pt_target_filter_isp(
				tfisp_pt_id, tfisp_isp_id
			)
			VALUES (83052, 5), (83053, 9)
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(88071, 'cpl-host-1', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(50761, 88071, 'cl-01', 'cpl1.ccc@19.10', 0x0A193D4B), -- All Targets
				(50762, 88071, 'cl-02', 'cpl2.ccc@19.11', 0x0A193D4C), -- ISP Filter(Refreshed)
				(50763, 88071, 'cl-03', 'cpl3.ccc@19.12', 0x0A193D4D) -- Only probed by all
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES(50761, 83051), (50762, 83052), (50763, 83053)
			`,
			`
			INSERT INTO nqm_cache_agent_ping_list_log(apll_ag_id, apll_number_of_targets, apll_time_access, apll_time_refresh)
			VALUES(50762, 4, '2012-05-05 10:10:23', '2012-05-05 10:10:23')
			`,
			`
			-- apl_time_access: 2014-05-06T13:24:00+00:00
			INSERT INTO nqm_cache_agent_ping_list(apl_apll_ag_id, apl_tg_id, apl_min_period, apl_time_access)
			VALUES(50762, 80921, 23, FROM_UNIXTIME(1399382640)),
				(50762, 80922, 23, FROM_UNIXTIME(1399382640)),
				(50762, 80923, 23, FROM_UNIXTIME(1399382640)),
				(50762, 80924, 23, FROM_UNIXTIME(1399382640))
			`,
		)
	case "TestDbNqmSuite.TestGetPingList":
		executeInTx(
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(30521, 'bmw-1'), (30522, 'bmw-2'), (30523, 'bmw-3'), (30524, 'bmw-4')
			`,
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3041, 'gd-01')
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(30291, 'tl-01', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(70071, 30291, 'tl-01', 'ccb1.ccc', 0x0A193D4B),
				(70072, 30291, 'tl-02', 'ccb2.ccc', 0x0A193D4C)
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
			)
			VALUES
				(99021, 'tgn-1', '1.2.3.4', 1, 3, 245, -1),
				(99022, 'tgn-2', '1.2.3.5', 2, 4, 1, 3041),
				(99023, 'tgn-3', '1.2.3.6', 1, 3, 246, -1),
				(99024, 'tgn-4', '20.91.83.101', 5, 3, 246, -1) -- As same as agent's IP
			`,
			`
			INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
			VALUES (99021, 30521), (99021, 30522),
				 (99022, 30522), (99022, 30523)
			`,
			`
			INSERT INTO nqm_cache_agent_ping_list_log(
				apll_ag_id, apll_number_of_targets, apll_time_access, apll_time_refresh
			)
			VALUES(70071, 3, '2012-03-04T06:12:33', '2012-03-04T06:12:33'),
				/**
				 * The period is not elapsed
				 *
				 * 2015-06-21T20:33:43Z
				 */
				(70072, 3, FROM_UNIXTIME(1434918823), FROM_UNIXTIME(1434918823))
			`,
			`
			INSERT INTO nqm_cache_agent_ping_list(
				apl_apll_ag_id, apl_tg_id, apl_min_period, apl_time_access
			) VALUES
				(70071, 99021, 5, FROM_UNIXTIME(0)),
				(70071, 99022, 5, FROM_UNIXTIME(0)),
				(70071, 99023, 5, FROM_UNIXTIME(0)),
				(70071, 99024, 5, FROM_UNIXTIME(0)),
				/**
				 *  Use 2015-06-21T20:30:00Z as access time
				 */
				(70072, 99021, 5, FROM_UNIXTIME(1434918600)),
				(70072, 99022, 5, FROM_UNIXTIME(1434918600)),
				(70072, 99023, 5, FROM_UNIXTIME(1434918600)),
				(70072, 99024, 5, FROM_UNIXTIME(1434918600))
				-- :~)
			`,
		)
	case "TestDbNqmSuite.TestUpdateAccessTime":
		executeInTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(98013, 'acc-host-1', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES (33404, 98013, 'acc-host-1@98.91', 'acc-host-1@98.91', 0x0A193D4B)
			`,
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host)
			VALUES
				(77061, 'tgn-1', '1.2.3.4'),
				(77062, 'tgn-2', '1.2.3.5'),
				(77063, 'tgn-3', '1.2.3.6')
			`,
			`
			/**
			 * 2015-06-21T20:33:43Z
			 */
			INSERT INTO nqm_cache_agent_ping_list_log(
				apll_ag_id, apll_number_of_targets, apll_time_access, apll_time_refresh
			)
			VALUES (33404, 3, FROM_UNIXTIME(1434918823), FROM_UNIXTIME(1434918823))
			`,
			`
			/**
			 * 2015-06-21T20:33:43Z
			 */
			INSERT INTO nqm_cache_agent_ping_list(
				apl_apll_ag_id, apl_tg_id, apl_min_period, apl_time_access
			) VALUES
				(33404, 77061, 60, FROM_UNIXTIME(1434918823)),
				(33404, 77062, 60, FROM_UNIXTIME(1434918823)),
				(33404, 77063, 60, FROM_UNIXTIME(1434918823))
			`,
		)
	}
}
func (s *TestDbNqmSuite) TearDownTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestDbNqmSuite.Test_vw_enabled_targets_by_ping_task":
		executeInTx(
			`DELETE FROM nqm_ping_task WHERE pt_id >= 47301 AND pt_id <= 47312`,
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
			"DELETE FROM nqm_agent WHERE ag_connection_id LIKE 'agrh33-%'",
			"DELETE FROM host WHERE hostname LIKE 'rhs1.%'",
		)
	case "TestDbNqmSuite.TestGetCacheLogOfPingList":
		executeInTx(
			"DELETE FROM nqm_cache_agent_ping_list_log WHERE apll_ag_id = 80981",
			"DELETE FROM nqm_agent WHERE ag_id = 80981",
			"DELETE FROM host WHERE id = 45701",
		)
	case "TestDbNqmSuite.TestBuildCacheOfPingList":
		executeInTx(
			"DELETE FROM nqm_cache_agent_ping_list_log WHERE apll_ag_id >= 50761 AND apll_ag_id <= 50763",
			"DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 50761 AND apt_ag_id <= 50763",
			"DELETE FROM nqm_agent WHERE ag_id >= 50761 AND ag_id <= 50763",
			"DELETE FROM host WHERE id = 88071",
			"DELETE FROM nqm_ping_task WHERE pt_id >= 83051 AND pt_id <= 83053",
			"DELETE FROM nqm_target WHERE tg_id >= 80921 AND tg_id <= 80929",
		)
	case "TestDbNqmSuite.TestGetPingList":
		executeInTx(
			"DELETE FROM nqm_cache_agent_ping_list_log WHERE apll_ag_id >= 70071 AND apll_ag_id <= 70072",
			"DELETE FROM nqm_agent WHERE ag_id >= 70071 AND ag_id <= 70072",
			"DELETE FROM host WHERE id = 30291",
			"DELETE FROM nqm_target WHERE tg_id >= 99021 ANd tg_id <= 99024",
			"DELETE FROM owl_name_tag WHERE nt_id = 3041",
			"DELETE FROM owl_group_tag WHERE gt_id >= 30521 AND gt_id <= 30524",
		)
	case "TestDbNqmSuite.TestUpdateAccessTime":
		executeInTx(
			`
			DELETE FROM nqm_cache_agent_ping_list_log
			WHERE apll_ag_id = 33404
			`,
			`
			DELETE FROM nqm_target
			WHERE tg_id >= 77061 AND tg_id <= 77063
			`,
			`DELETE FROM nqm_agent WHERE ag_id = 33404`,
			`DELETE FROM host WHERE id = 98013`,
		)
	}
}

func (s *TestDbNqmSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}
func (s *TestDbNqmSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
