package db

import (
	"database/sql"
	"net"
	"sort"
	"time"

	commonModel "github.com/Cepave/common/model"
	"github.com/Cepave/hbs/model"
	hbstesting "github.com/Cepave/hbs/testing"
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
		{230001, []int{402001, 402002, 402003, 402004, 402005}}, // All of the targets
		{230002, []int{402001, 402002}},                         // Targets are matched by ISP
		{230003, []int{402001, 402003}},                         // Targets are matched by province
		{230004, []int{402001, 402004}},                         // Targets are matched by city
		{230005, []int{402001, 402005}},                         // Targets are matched by name tag
		{230006, []int{402001}},                                 // Nothing matched except probed by all
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
	checkTime string
	checker   Checker
}

func (suite *TestDbNqmSuite) TestGetAndRefreshNeedPingAgentForRpc(c *C) {
	testedCases := []getAndRefreshNeedPingAgentTestCase{
		{130001, "2115-08-08T00:00:00+00:00", IsNil},  // No ping task setting
		{130002, "2010-05-05T10:59:00+08:00", IsNil},  // The period is not elapsed yet
		{130003, "2013-10-01T00:00:00+08:00", NotNil}, // Never executed
		{130004, "2012-06-10T09:00:00+08:00", NotNil}, // The period is elapsed
		{130005, "2012-06-10T09:00:00+08:00", IsNil},  // Disabled agent
	}

	for _, v := range testedCases {
		testNeedPingAgent(c, v)
	}
}

func testNeedPingAgent(c *C, testCase getAndRefreshNeedPingAgentTestCase) {
	sampleCheckedTime, _ := time.Parse(time.RFC3339, testCase.checkTime)

	c.Logf("Agent Id: %d", testCase.agentId)

	testedAgent, err := GetAndRefreshNeedPingAgentForRpc(
		testCase.agentId, sampleCheckedTime,
	)

	/**
	 * Asserts the result data
	 */
	c.Assert(err, IsNil)
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
		SELECT UNIX_TIMESTAMP(pt_time_last_execute)
		FROM nqm_ping_task
		WHERE pt_ag_id = ?
		`,
		testCase.agentId,
	)

	c.Assert(unixTime, Equals, sampleCheckedTime.Unix())
	// :~)
}

/**
 * Tests the state of ping task
 */
func (suite *TestDbNqmSuite) TestGetPingTaskState(c *C) {
	testedCases := []struct {
		agentId        int
		expectedStatus int
	}{
		{2001, NO_PING_TASK},               // The agent has no ping task
		{2002, HAS_PING_TASK_ALL_MATCHING}, // The agent has ping task with all of the targets
		{2003, HAS_PING_TASK_WITH_FILTER},  // The agent has ping task with ISP filter
		{2004, HAS_PING_TASK_WITH_FILTER},  // The agent has ping task with province filter
		{2005, HAS_PING_TASK_WITH_FILTER},  // The agent has ping task with city filter
		{2006, HAS_PING_TASK_WITH_FILTER},  // The agent has ping task with name tag
	}

	for _, v := range testedCases {
		testedResult, err := getPingTaskState(v.agentId)

		c.Assert(err, IsNil)
		c.Assert(testedResult, Equals, v.expectedStatus)
	}
}

func (s *TestDbNqmSuite) SetUpTest(c *C) {
	if !hbstesting.HasDbEnvForMysqlOrSkip(c) {
		return
	}

	switch c.TestName() {
	case "TestDbNqmSuite.TestGetAndRefreshNeedPingAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			`SET time_zone = '+08:00'`,
			`
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id, ag_status)
			VALUES
				(130001, 'gc-1', 'tt1.org', 0x12345678, 3, DEFAULT),
				(130002, 'gc-2', 'tt2.org', 0x13345678, 3, DEFAULT),
				(130003, 'gc-3', 'tt3.org', 0x14345678, 3, DEFAULT),
				(130004, 'gc-4', 'tt4.org', 0x15345678, 3, DEFAULT),
				(130005, 'gc-5', 'tt5.org', 0x15345678, 3, b'00000000')
			`,
			`
			INSERT INTO nqm_ping_task(pt_ag_id, pt_period, pt_time_last_execute)
			VALUES
				(130002, 60, '2010-05-05 10:00:00'),
				(130003, 60, NULL),
				(130004, 60, '2012-06-10 08:00:00'),
				(130005, 60, NULL)
			`,
		)
	case "TestDbNqmSuite.TestGetPingTaskState":
		hbstesting.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(2001, 'pt-01', 'aaa1.ccc', 0x12345678),
				(2002, 'pt-02', 'aaa2.ccc', 0x13345678),
				(2003, 'pt-03', 'aaa3.ccc', 0x14345678),
				(2004, 'pt-04', 'aaa4.ccc', 0x14445678),
				(2005, 'pt-05', 'aaa5.ccc', 0x14745678),
				(2006, 'pt-06', 'aaa6.ccc', 0x14765678)
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_ag_id, pt_period
			)
			VALUES
				(2002, 20),
				(2003, 20),
				(2004, 20),
				(2005, 20),
				(2006, 20)
			`,
			`
			INSERT INTO nqm_pt_target_filter_isp(
				tfisp_pt_ag_id, tfisp_isp_id
			) VALUES (2003, 1)
			`,
			`
			INSERT INTO nqm_pt_target_filter_province(
				tfpv_pt_ag_id, tfpv_pv_id
			) VALUES (2004, 1)
			`,
			`
			INSERT INTO nqm_pt_target_filter_city(
				tfct_pt_ag_id, tfct_ct_id
			) VALUES (2005, 1)
			`,
			`
			INSERT INTO nqm_pt_target_filter_name_tag(
				tfnt_pt_ag_id, tfnt_name_tag
			) VALUES (2006, 'st-1')
			`,
		)
	case "TestDbNqmSuite.TestGetTargetsByAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES
				(230001, 'tl-01', 'ccb1.ccc', 0x12345678),
				(230002, 'tl-02', 'ccb2.ccc', 0x22345678),
				(230003, 'tl-03', 'ccb3.ccc', 0x32345678),
				(230004, 'tl-04', 'ccb4.ccc', 0x42345678),
				(230005, 'tl-05', 'ccb5.ccc', 0x52345678),
				(230006, 'tl-06', 'ccb6.ccc', 0x62345678)
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_probed_by_all, tg_name_tag,
				tg_status, tg_available
			)
			VALUES
				(402001, 'tgn-1', '1.2.3.4', -1, -1, -1, true, null, true, true),
				(402002, 'tgn-2', '1.2.3.5', 5, -1, -1, false, null, true, true),
				(402003, 'tgn-3', '1.2.3.6', -1, 5, -1, false, null, true, true),
				(402004, 'tgn-4', '1.2.3.7', -1, 20, 20, false, null, true, true),
				(402005, 'tgn-5', '1.2.3.8', -1, -1, -1, false, 'tag-1', true, true)
			`,
			`
			INSERT INTO nqm_ping_task(
				pt_ag_id, pt_period
			)
			VALUES
				(230001, 20),
				(230002, 20),
				(230003, 20),
				(230004, 20),
				(230005, 20),
				(230006, 20)
			`,
			`
			INSERT INTO nqm_pt_target_filter_isp(
				tfisp_pt_ag_id, tfisp_isp_id
			)
			VALUES (230002, 5)
			`,
			`
			INSERT INTO nqm_pt_target_filter_province(
				tfpv_pt_ag_id, tfpv_pv_id
			)
			VALUES (230003, 5)
			`,
			`
			INSERT INTO nqm_pt_target_filter_city(
				tfct_pt_ag_id, tfct_ct_id
			)
			VALUES (230004, 20)
			`,
			`
			INSERT INTO nqm_pt_target_filter_name_tag(
				tfnt_pt_ag_id, tfnt_name_tag
			)
			VALUES (230005, 'tag-1'),
				(230006, 'tag-nothing-matched')
			`,
		)
	}
}

func (s *TestDbNqmSuite) TearDownTest(c *C) {
	switch c.TestName() {
	case "TestDbNqmSuite.TestRefreshAgentInfo":
		hbstesting.ExecuteOrFail(
			"DELETE FROM nqm_agent WHERE ag_connection_id = 'refresh-1'",
		)
	case "TestDbNqmSuite.TestGetAndRefreshNeedPingAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			"SET FOREIGN_KEY_CHECKS=0",
			"DELETE FROM nqm_ping_task",
			"DELETE FROM nqm_agent",
			"SET FOREIGN_KEY_CHECKS=1",
		)
	case "TestDbNqmSuite.TestGetPingTaskState":
		hbstesting.ExecuteQueriesOrFailInTx(
			"SET FOREIGN_KEY_CHECKS=0",
			"DELETE FROM nqm_pt_target_filter_isp",
			"DELETE FROM nqm_pt_target_filter_province",
			"DELETE FROM nqm_pt_target_filter_city",
			"DELETE FROM nqm_pt_target_filter_name_tag",
			"DELETE FROM nqm_ping_task",
			"DELETE FROM nqm_agent",
			"SET FOREIGN_KEY_CHECKS=1",
		)
	case "TestDbNqmSuite.TestGetTargetsByAgentForRpc":
		hbstesting.ExecuteQueriesOrFailInTx(
			"SET FOREIGN_KEY_CHECKS=0",
			"DELETE FROM nqm_pt_target_filter_isp",
			"DELETE FROM nqm_pt_target_filter_province",
			"DELETE FROM nqm_pt_target_filter_city",
			"DELETE FROM nqm_pt_target_filter_name_tag",
			"DELETE FROM nqm_ping_task",
			"DELETE FROM nqm_target",
			"DELETE FROM nqm_agent",
			"SET FOREIGN_KEY_CHECKS=1",
		)
	}
}
