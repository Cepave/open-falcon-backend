package nqm

import (
	qtest "github.com/Cepave/open-falcon-backend/modules/query/test"
	. "gopkg.in/check.v1"
)

type TestNqmRdbSuite struct{}

var _ = Suite(&TestNqmRdbSuite{})

func (s *TestNqmRdbSuite) SetUpSuite(c *C) {
	qtest.InitOrm()
	qtest.InitDb()
}

func (s *TestNqmRdbSuite) TearDownSuite(c *C) {
	qtest.ReleaseDb()
}

/**
 * Tests the listing for agents in city by id of province
 */
func (suite *TestNqmRdbSuite) TestListAgentsInCityByProvinceId(c *C) {
	testedResult := ListAgentsInCityByProvinceId(3)

	c.Assert(testedResult, HasLen, 2)

	for _, agentsInCity := range testedResult {
		c.Logf("Data for city[%v]: [%v]", agentsInCity.City, agentsInCity.Agents)

		switch agentsInCity.City.Id {
		case 2:
			c.Assert(agentsInCity.Agents, HasLen, 2)
			c.Assert(agentsInCity.Agents[0].Id, Equals, int32(14301))
			c.Assert(agentsInCity.Agents[1].Id, Equals, int32(14302))
		case 3:
			c.Assert(agentsInCity.Agents, HasLen, 2)
			c.Assert(agentsInCity.Agents[0].Id, Equals, int32(14303))
			c.Assert(agentsInCity.Agents[1].Id, Equals, int32(14304))
		default:
			c.Fatalf("Unexcepted city id: [%v]", agentsInCity.City.Id)
		}
	}
}

// Tests the getting of data for province by name
type getProvinceByNameTestCase struct {
	searchText   string
	expectedId   int32
	expectedName string
}

func (suite *TestNqmRdbSuite) TestGetProvinceByName(c *C) {
	testCases := []getProvinceByNameTestCase{
		{"天津", 12, "天津"},
		{"天", 12, "天津"},
		{"無此省份", UNKNOWN_ID_FOR_QUERY, "無此省份"},
		{"無此省份2", UNKNOWN_ID_FOR_QUERY, "無此省份2"},
		{"贵州", 22, "贵州"},
	}

	for _, testCase := range testCases {
		testedProvince := getProvinceByName(testCase.searchText)

		c.Logf("Test load province by name \"%v\". Got: %v", testCase.searchText, testedProvince)

		c.Assert(testedProvince, NotNil)
		c.Assert(testedProvince.Id, Equals, int16(testCase.expectedId))
		c.Assert(testedProvince.Name, Equals, testCase.expectedName)
	}
}

// Tests the getting of data for province by id
type getProvinceByIdTestCase struct {
	searchId     int16
	expectedId   int16
	expectedName string
}

func (suite *TestNqmRdbSuite) TestGetProvinceById(c *C) {
	testCases := []getProvinceByIdTestCase{
		{12, 12, "天津"},
		{23, 23, "海南"},
		{-1, -1, "<UNDEFINED>"},
		{-919, -919, UNKNOWN_NAME_FOR_QUERY},
	}

	for _, testCase := range testCases {
		testedProvince := getProvinceById(testCase.searchId)

		c.Logf("Test load province by id \"%d\". Got: %v", testCase.searchId, testedProvince)

		c.Assert(testedProvince, NotNil)
		c.Assert(testedProvince.Id, Equals, int16(testCase.expectedId))
		c.Assert(testedProvince.Id, Equals, testCase.expectedId)
	}
}

// Tests the getting of data for ISP by name
type getIspByNameTestCase struct {
	searchText   string
	expectedId   int32
	expectedName string
}

func (suite *TestNqmRdbSuite) TestGetIspByName(c *C) {
	testCases := []getIspByNameTestCase{
		{"方正", 8, "方正宽带"},
		{"方正宽带", 8, "方正宽带"},
		{"無此ISP", UNKNOWN_ID_FOR_QUERY, "無此ISP"},
		{"無此ISP2", UNKNOWN_ID_FOR_QUERY, "無此ISP2"},
		{"中信网络", 21, "中信网络"},
	}

	for _, testCase := range testCases {
		testedIsp := getIspByName(testCase.searchText)

		c.Logf("Test load ISP by name \"%v\". Got: %v", testCase.searchText, testedIsp)

		c.Assert(testedIsp, NotNil)
		c.Assert(testedIsp.Id, Equals, int16(testCase.expectedId))
		c.Assert(testedIsp.Name, Equals, testCase.expectedName)
	}
}

// Tests the getting of data for ISP by id
type getIspByIdTestCase struct {
	searchId     int16
	expectedId   int16
	expectedName string
}

func (suite *TestNqmRdbSuite) TestGetIspById(c *C) {
	testCases := []getIspByIdTestCase{
		{8, 8, "方正宽带"},
		{11, 11, "长城宽带"},
		{-1, -1, "<UNDEFINED>"},
		{-919, -919, UNKNOWN_NAME_FOR_QUERY},
	}

	for _, testCase := range testCases {
		testedIsp := getIspById(testCase.searchId)

		c.Logf("Test load ISP by id \"%d\". Got: %v", testCase.searchId, testedIsp)

		c.Assert(testedIsp, NotNil)
		c.Assert(testedIsp.Id, Equals, int16(testCase.expectedId))
		c.Assert(testedIsp.Id, Equals, testCase.expectedId)
	}
}

// Tests the getting of data for city by name
type getCityByNameTestCase struct {
	searchText       string
	expectedId       int32
	expectedName     string
	expectedPostCode string
}

func (suite *TestNqmRdbSuite) TestGetCityByName(c *C) {
	testCases := []getCityByNameTestCase{
		{"茂名", 20, "茂名市", "525000"},
		{"株洲市", 32, "株洲市", "412000"},
		{"無此city", UNKNOWN_ID_FOR_QUERY, "無此city", UNKNOWN_NAME_FOR_QUERY},
		{"無此city2", UNKNOWN_ID_FOR_QUERY, "無此city2", UNKNOWN_NAME_FOR_QUERY},
	}

	for _, testCase := range testCases {
		testedCity := getCityByName(testCase.searchText)

		c.Logf("Test load city by name \"%v\". Got: %v", testCase.searchText, testedCity)

		c.Assert(testedCity, NotNil)
		c.Assert(testedCity.Id, Equals, int16(testCase.expectedId))
		c.Assert(testedCity.Name, Equals, testCase.expectedName)
		c.Assert(testedCity.PostCode, Equals, testCase.expectedPostCode)
	}
}

// Tests the getting of data for city by id
type getCityByIdTestCase struct {
	searchId         int16
	expectedId       int16
	expectedName     string
	expectedPostCode string
}

func (suite *TestNqmRdbSuite) TestGetCityById(c *C) {
	testCases := []getCityByIdTestCase{
		{48, 48, "荆州市", "434100"},
		{33, 33, "娄底市", "417000"},
		{-1, -1, "<UNDEFINED>", "<UNDEFINED>"},
		{-919, -919, UNKNOWN_NAME_FOR_QUERY, UNKNOWN_NAME_FOR_QUERY},
	}

	for _, testCase := range testCases {
		testedCity := getCityById(testCase.searchId)

		c.Logf("Test load city by id \"%d\". Got: %v", testCase.searchId, testedCity)

		c.Assert(testedCity, NotNil)
		c.Assert(testedCity.Id, Equals, int16(testCase.expectedId))
		c.Assert(testedCity.Name, Equals, testCase.expectedName)
		c.Assert(testedCity.PostCode, Equals, testCase.expectedPostCode)
	}
}

// Tests the getting of data for target by id
type getTargetByIdTestCase struct {
	searchId     int32
	expectedId   int32
	expectedHost string
}

func (suite *TestNqmRdbSuite) TestGetTargetById(c *C) {
	testCases := []getTargetByIdTestCase{
		{19203, 19203, "100.20.50.3"},
		{19202, 19202, "100.20.50.2"},
		{28001, 28001, UNKNOWN_NAME_FOR_QUERY},
	}

	for _, testCase := range testCases {
		testedTarget := getTargetById(testCase.searchId)

		c.Logf("Test load target by id \"%v\". Got: %v", testCase.searchId, testedTarget)

		c.Assert(testedTarget, NotNil)
		c.Assert(testedTarget.Id, Equals, testCase.expectedId)
		c.Assert(testedTarget.Host, Equals, testCase.expectedHost)
	}
}

// Tests the getting of data for target by host
type getTargetByHostTestCase struct {
	searchText   string
	expectedId   int32
	expectedHost string
}

func (suite *TestNqmRdbSuite) TestGetTargetByHost(c *C) {
	testCases := []getTargetByHostTestCase{
		{"100.20.50.1", 19201, "100.20.50.1"},
		{"100.20.50.2", 19202, "100.20.50.2"},
		{"無此target", UNKNOWN_ID_FOR_QUERY, "無此target"},
	}

	for _, testCase := range testCases {
		testedTarget := getTargetByHost(testCase.searchText)

		c.Logf("Test load target by host \"%v\". Got: %v", testCase.searchText, testedTarget)

		c.Assert(testedTarget, NotNil)
		c.Assert(testedTarget.Id, Equals, testCase.expectedId)
		c.Assert(testedTarget.Host, Equals, testCase.expectedHost)
	}
}

func (s *TestNqmRdbSuite) SetUpTest(c *C) {
	if !qtest.HasDefaultOrmOnPortal(c) {
		return
	}

	switch c.TestName() {
	case "TestNqmRdbSuite.TestGetTargetByHost",
		"TestNqmRdbSuite.TestGetTargetById":
		qtest.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host) VALUES
			(19201, 'target-1', '100.20.50.1'),
			(19202, 'target-2', '100.20.50.2'),
			(19203, 'target-3', '100.20.50.3')
			`,
		)
	case "TestNqmRdbSuite.TestListAgentsInCityByProvinceId":
		qtest.ExecuteQueriesOrFailInTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(23401, 'agent-1', '', '')
			`,
			`
			INSERT INTO nqm_agent(
				ag_id, ag_hs_id, ag_name, ag_connection_id,
				ag_hostname, ag_ip_address,
				ag_isp_id, ag_pv_id, ag_ct_id, ag_status
			) VALUES
			(14301, 23401, 'agent-1', '92.20.50.1', 'ag-1-host', 0x5A143201, -1, 3, 2, true),
			(14302, 23401, 'agent-2', '92.20.50.2', 'ag-2-host', 0x5A143202, -1, 3, 2, true),
			(14303, 23401, 'agent-3', '92.20.50.3', 'ag-3-host', 0x5A143203, -1, 3, 3, true),
			(14304, 23401, 'agent-4', '92.20.50.4', 'ag-4-host', 0x5A143204, -1, 3, 3, true),
			(14305, 23401, 'agent-5', '92.20.50.5', 'ag-5-host', 0x5A143205, -1, 3, 3, false), # Agent disabled
			(14306, 23401, 'agent-6', '92.20.50.6', 'ag-6-host', 0x5A143206, -1, 3, 3, false), # Task disabled
			(14311, 23401, 'agent-11', '92.20.50.11', 'ag-11-host', 0x5A14320B, -1, 3, 3, true), # Has no task
			(14312, 23401, 'agent-12', '92.20.50.12', 'ag-12-host', 0x5A14320C, -1, 4, 3, true) # Not-matched province
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_period, pt_enable)
			VALUES
			(7701, 100, TRUE),
			(7702, 100, TRUE),
			(7703, 100, TRUE),
			(7704, 100, TRUE),
			(7705, 100, TRUE),
			(7706, 100, FALSE)
			`,
			`
			INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
			VALUES(14301, 7701),
				(14302, 7702),
				(14303, 7703),
				(14304, 7704),
				(14305, 7705),
				(14306, 7706)
			`,
		)
	}
}

func (s *TestNqmRdbSuite) TearDownTest(c *C) {
	if !qtest.HasTestDbConfig() {
		return
	}

	switch c.TestName() {
	case "TestNqmRdbSuite.TestGetTargetByHost",
		"TestNqmRdbSuite.TestGetTargetById":
		qtest.ExecuteQueriesOrFailInTx(
			`DELETE FROM nqm_target WHERE tg_id IN (19201, 19202, 19203)`,
		)
	case "TestNqmRdbSuite.TestListAgentsInCityByProvinceId":
		qtest.ExecuteQueriesOrFailInTx(
			`DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 14301 AND apt_ag_id <= 14320`,
			`DELETE FROM nqm_ping_task WHERE pt_id >= 7701 AND pt_id <= 7706`,
			`DELETE FROM nqm_agent WHERE ag_id >= 14301 AND ag_id <= 14320`,
			`DELETE FROM host WHERE id = 23401`,
		)
	}
}
