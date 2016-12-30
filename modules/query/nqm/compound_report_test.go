package nqm

import (
	"fmt"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/satori/go.uuid"
	. "gopkg.in/check.v1"
)

type TestCompoundReportSuite struct{}
type TestCompoundReportSuiteOnDb struct{
	*dbTestSuite
}

var (
	_ = Suite(&TestCompoundReportSuite{})
	_ = Suite(&TestCompoundReportSuiteOnDb{ &dbTestSuite{} })
)

// Tests the convertion of query to detail information
func (suite *TestCompoundReportSuiteOnDb) TestToQueryDetail(c *C) {
	/**
	 * Sets-up sample query
	 */
	sampleQuery := model.NewCompoundQuery()

	sampleQuery.Filters.Metrics = "$max >= $min"
	sampleQuery.Grouping.Agent = []string{ model.AgentGroupingHostname, model.GroupingIsp }
	sampleQuery.Grouping.Target = []string{ model.TargetGroupingHost, model.GroupingProvince, model.GroupingIsp }
	sampleQuery.Output.Metrics = []string{ model.MetricAvg, model.MetricMax, model.MetricLoss }

	agentFilter := sampleQuery.Filters.Agent

	agentFilter.Name = []string { "Cool1", "Cool2" }
	agentFilter.Hostname = []string { "gc1.com", "gc2.com" }
	agentFilter.IpAddress = []string { "123.71.1", "123.71.2" }
	agentFilter.ConnectionId = []string { "c01", "c02" }
	agentFilter.IspIds = []int16 { 2, 3, 4, 1, 5, 6 }
	agentFilter.ProvinceIds = []int16 { 3, 4, 7, 8 }
	agentFilter.CityIds = []int16 { 13, 14 }
	agentFilter.NameTagIds = []int16 { 3375, 3376 }
	agentFilter.GroupTagIds = []int32 { 90801, 90802, 90803 }

	targetFilter := sampleQuery.Filters.Target

	targetFilter.Name = []string { "Zoo-1", "Zoo-2" }
	targetFilter.Host = []string { "kz1.com", "kz2.com" }
	targetFilter.IspIds = []int16 { 5, 7 }
	targetFilter.ProvinceIds = []int16 { 1, 3, 8 }
	targetFilter.CityIds = []int16 { 15, 16 }
	targetFilter.NameTagIds = []int16 { 3375, 3376 }
	targetFilter.GroupTagIds = []int32 { 90801, 90802, 90803 }
	// :~)

	testedDetail := ToQueryDetail(sampleQuery)

	c.Assert(string(testedDetail.Metrics), Equals, sampleQuery.Filters.Metrics)

	/**
	 * Asserts the query detail on agent conditions
	 */
	testedAgentDetail := testedDetail.Agent
	c.Assert(testedAgentDetail.Name, DeepEquals, agentFilter.Name)
	c.Assert(testedAgentDetail.Hostname, DeepEquals, agentFilter.Hostname)
	c.Assert(testedAgentDetail.IpAddress, DeepEquals, agentFilter.IpAddress)
	c.Assert(testedAgentDetail.ConnectionId, DeepEquals, agentFilter.ConnectionId)
	c.Assert(testedAgentDetail.Isps, HasLen, 6)
	c.Assert(testedAgentDetail.Provinces, HasLen, 4)
	c.Assert(testedAgentDetail.Cities, HasLen, 2)
	c.Assert(testedAgentDetail.NameTags, HasLen, 2)
	c.Assert(testedAgentDetail.GroupTags, HasLen, 3)
	// :~)

	/**
	 * Asserts the query detail on target conditions
	 */
	testedTargetDetail := testedDetail.Target
	c.Assert(testedTargetDetail.Name, DeepEquals, targetFilter.Name)
	c.Assert(testedTargetDetail.Host, DeepEquals, targetFilter.Host)
	c.Assert(testedTargetDetail.Isps, HasLen, 2)
	c.Assert(testedTargetDetail.Provinces, HasLen, 3)
	c.Assert(testedTargetDetail.Cities, HasLen, 2)
	c.Assert(testedTargetDetail.NameTags, HasLen, 2)
	c.Assert(testedTargetDetail.GroupTags, HasLen, 3)
	// :~)

	/**
	 * Asserts the output detail
	 */
	testedOutputDetail := testedDetail.Output
	c.Assert(testedOutputDetail.Agent, HasLen, 2)
	c.Assert(testedOutputDetail.Target, HasLen, 3)
	c.Assert(testedOutputDetail.Metrics, HasLen, 3)
	// :~)
}

// Tests the building of query object
func (suite *TestCompoundReportSuiteOnDb) TestBuildQuery(c *C) {
	sampleQuery := model.NewCompoundQuery()

	sampleJson := []byte(`
	{
		"filters": {
			"time": {
				"start_time": 1336608000,
				"end_time": 1336622400
			},
			"agent": {
				"name": [ "GD-1", "GD-2" ]
			},
			"target": {
				"host": [ "18.98.7.61", "google.com" ]
			}
		}
	}`)
	c.Assert(sampleQuery.UnmarshalJSON(sampleJson), IsNil)

	testedResult1 := BuildQuery(sampleQuery)

	c.Logf("[T-1] Query object: %s", testedResult1)
	c.Assert(testedResult1.NamedId, Equals, queryNamedId)

	/**
	 * Asserts sample query with same conditions
	 */
	testedResult2 := BuildQuery(sampleQuery)
	c.Logf("[T-2] Query object: %s", testedResult2)
	c.Assert(testedResult1, DeepEquals, testedResult2)
	// :~)
}

// Tests the loading of compound query by UUID
func (suite *TestCompoundReportSuiteOnDb) TestGetCompoundQueryByUuid(c *C) {
	sampleQuery := model.NewCompoundQuery()
	err := sampleQuery.UnmarshalJSON([]byte(`
	{
		"filters": {
			"time": {
				"start_time": 190807000,
				"end_time": 190827000
			}
		},
		"output": {
			"metrics": [ "max", "min", "avg", "loss" ]
		}
	}
	`))
	c.Assert(err, IsNil)

	/**
	 * Builds query object(persist the query) and load it with generated UUID
	 */
	queryObject := BuildQuery(sampleQuery)
	testedQuery := GetCompoundQueryByUuid(uuid.UUID(queryObject.Uuid))
	// :~)

	c.Assert(testedQuery.Filters.Time, DeepEquals, sampleQuery.Filters.Time)
	c.Assert(testedQuery.Output, DeepEquals, sampleQuery.Output)
}

// Tests the building of NQM dsl by compound query
func (suite *TestCompoundReportSuiteOnDb) TestBuildNqmDslByCompoundQuery(c *C) {
	type subCase [][2]interface{}

	allTestCases := []*struct {
		nodeProperty string
		queryProperty string
		testedProperty string

		sampleAndExpected subCase
	} {
		// For agents

		{
			"Agent", "Name", "IdsOfAgents",
			subCase {
				{ []string{}, []int32{}, },
				{ []string{ "GC-01" }, []int32{ 1041 }, },
				{ []string{ "No-1" }, []int32{ -2 }, },
			},
		},
		{
			"Agent", "ConnectionId", "IdsOfAgents",
			subCase {
				{ []string{}, []int32{}, },
				{ []string{ "eth3-gc-01" }, []int32{ 1041 }, },
				{ []string{ "No-1" }, []int32{ -2 }, },
			},
		},
		{
			"Agent", "Hostname", "IdsOfAgents",
			subCase {
				{ []string{}, []int32{}, },
				{ []string{ "KCB-01" }, []int32{ 1041 }, },
				{ []string{ "No-1" }, []int32{ -2 }, },
			},
		},
		{
			"Agent", "IpAddress", "IdsOfAgents",
			subCase {
				{ []string{}, []int32{}, },
				{ []string{ "10.91.8.33" }, []int32{ 1041 }, },
				{ []string{ "90.11.76.2" }, []int32{ -2 }, },
			},
		},
		{
			"Agent", "IspIds", "IdsOfAgentIsps",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 10, 20 }, []int16{ 10, 20 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Agent", "IspIds", "IspRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 10, 20 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Agent", "ProvinceIds", "IdsOfAgentProvinces",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 32, 31 }, []int16{ 32, 31 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Agent", "ProvinceIds", "ProvinceRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 7 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Agent", "CityIds", "IdsOfAgentCities",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 65, 72 }, []int16{ 65, 72 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Agent", "CityIds", "CityRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 80 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Agent", "NameTagIds", "IdsOfAgentNameTags",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 101, 114 }, []int16{ 101, 114 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Agent", "NameTagIds", "NameTagRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 29 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Agent", "GroupTagIds", "IdsOfAgentGroupTags",
			subCase {
				{ []int32{}, []int32{}, },
				{ []int32{ 1291, 1309 }, []int32{ 1291, 1309 }, },
				{ []int32{ model.RelationSame }, []int32{}, },
				{ []int32{ model.RelationNotSame }, []int32{}, },
			},
		},

		// For targets
		{
			"Target", "Name", "IdsOfTargets",
			subCase {
				{ []string{}, []int32{}, },
				{ []string{ "ZKP-01" }, []int32{ 2301 }, },
				{ []string{ "No-1" }, []int32{ -2 }, },
			},
		},
		{
			"Target", "Host", "IdsOfTargets",
			subCase {
				{ []string{}, []int32{}, },
				{ []string{ "ZKP-TTC-33" }, []int32{ 2301 }, },
				{ []string{ "No-1" }, []int32{ -2 }, },
			},
		},
		{
			"Target", "IspIds", "IdsOfTargetIsps",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 10, 20 }, []int16{ 10, 20 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Target", "IspIds", "IspRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 10, 20 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Target", "ProvinceIds", "IdsOfTargetProvinces",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 32, 31 }, []int16{ 32, 31 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Target", "ProvinceIds", "ProvinceRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 7 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Target", "CityIds", "IdsOfTargetCities",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 65, 72 }, []int16{ 65, 72 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Target", "CityIds", "CityRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 80 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Target", "NameTagIds", "IdsOfTargetNameTags",
			subCase {
				{ []int16{}, []int16{}, },
				{ []int16{ 101, 114 }, []int16{ 101, 114 }, },
				{ []int16{ model.RelationSame }, []int16{}, },
				{ []int16{ model.RelationNotSame }, []int16{}, },
			},
		},
		{
			"Target", "NameTagIds", "NameTagRelation",
			subCase {
				{ []int16{}, model.NoCondition, },
				{ []int16{ 29 }, model.NoCondition, },
				{ []int16{ model.RelationSame }, model.SameValue, },
				{ []int16{ model.RelationNotSame }, model.NotSameValue, },
			},
		},
		{
			"Target", "GroupTagIds", "IdsOfTargetGroupTags",
			subCase {
				{ []int32{}, []int32{}, },
				{ []int32{ 1291, 1309 }, []int32{ 1291, 1309 }, },
				{ []int32{ model.RelationSame }, []int32{}, },
				{ []int32{ model.RelationNotSame }, []int32{}, },
			},
		},
	}

	var newSampleQuery = func() *model.CompoundQuery {
		/**
		 * Prepares compound query of sample
		 */
		query := model.NewCompoundQuery()
		err := query.UnmarshalJSON(
			[]byte(`
			{
				"filters": {
					"time": {
						"start_time": 12088600,
						"end_time": 12088800
					}
				}
			}
			`),
		)
		c.Assert(err, IsNil)

		return query
	}

	for i, testCase := range allTestCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		sampleQuery := newSampleQuery()

		for _, subCase := range testCase.sampleAndExpected {
			c.Logf("\tSubCase sample: %v", subCase[0])

			oreflect.SetValueOfField(
				sampleQuery, subCase[0],
				"Filters", testCase.nodeProperty, testCase.queryProperty,
			)

			testedDsl := buildNqmDslByCompoundQuery(sampleQuery)

			/**
			 * Asserts time data
			 */
			c.Assert(testedDsl.StartTime, Equals, EpochTime(12088600), comment)
			c.Assert(testedDsl.EndTime, Equals, EpochTime(12088800), comment)
			// :~)

			/**
			 * Asserts the tested property
			 */
			c.Assert(
				oreflect.GetValueOfField(testedDsl, testCase.testedProperty),
				DeepEquals,
				subCase[1],
				comment,
			)
			// :~)
		}
	}
}

// Tests the building of group columns for DSL
func (suite *TestCompoundReportSuite) TestBuildGroupingColumnOfDsl(c *C) {
	testCases := []*struct {
		sampleGrouping *model.QueryGrouping
		expected []string
	} {
		{ // Agent + Target
			&model.QueryGrouping{
				Agent: []string { model.AgentGroupingName, model.AgentGroupingHostname, model.AgentGroupingIpAddress },
				Target: []string { model.TargetGroupingName, model.TargetGroupingHost },
			},
			[]string { "ag_id", "tg_id" },
		},
		{ // Other properties
			&model.QueryGrouping{
				Agent: []string { model.GroupingIsp, model.GroupingProvince, model.GroupingCity, model.GroupingNameTag, },
				Target: []string { model.GroupingIsp, model.GroupingProvince, model.GroupingCity, model.GroupingNameTag, },
			},
			[]string {
				"ag_isp_id", "ag_pv_id", "ag_ct_id", "ag_nt_id",
				"tg_isp_id", "tg_pv_id", "tg_ct_id", "tg_nt_id",
			},
		},
		{ // Agent + other property of target
			&model.QueryGrouping{
				Agent: []string { model.AgentGroupingName, model.GroupingIsp },
				Target: []string { model.GroupingIsp },
			},
			[]string { "ag_id", "tg_isp_id" },
		},
		{ // Other property of agent + target
			&model.QueryGrouping{
				Agent: []string { model.GroupingNameTag },
				Target: []string { model.TargetGroupingName, model.GroupingIsp },
			},
			[]string { "ag_nt_id", "tg_id" },
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedGrouping := buildGroupingColumnOfDsl(testCase.sampleGrouping)
		c.Assert(testedGrouping, DeepEquals, testCase.expected, comment)
	}
}

// Tests the setup of sorting properties
func (suite *TestCompoundReportSuite) TestSetupSorting(c *C) {
	testCases := []*struct {
		sampleEntities []*commonModel.OrderByEntity
		outputMetrics []string
		expectedResult []*commonModel.OrderByEntity
	} {
		{
			[]*commonModel.OrderByEntity{ { "agent_isp", commonModel.Descending } },
			[]string { "max", "min" },
			[]*commonModel.OrderByEntity{ { "agent_isp", commonModel.Descending }, { "loss", commonModel.Descending } },
		},
		{
			[]*commonModel.OrderByEntity{},
			[]string { "max", "min" },
			[]*commonModel.OrderByEntity{
				{ "max", commonModel.Descending },
				{ "min", commonModel.Descending },
				{ "loss", commonModel.Descending },
			},
		},
		{
			[]*commonModel.OrderByEntity{},
			[]string { "max", "avg", "loss" },
			[]*commonModel.OrderByEntity{
				{ "avg", commonModel.Descending },
				{ "loss", commonModel.Descending },
			},
		},
		{
			[]*commonModel.OrderByEntity{},
			[]string { "num_agent", "num_target" },
			[]*commonModel.OrderByEntity{
				{ "num_agent", commonModel.Descending },
				{ "num_target", commonModel.Descending },
				{ "loss", commonModel.Descending },
			},
		},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		sampleOutput := &model.QueryOutput{}
		sampleOutput.Metrics = testCase.outputMetrics

		samplePaging := commonModel.NewUndefinedPaging()
		samplePaging.OrderBy = testCase.sampleEntities

		setupSorting(samplePaging, sampleOutput)

		c.Assert(samplePaging.OrderBy, DeepEquals, testCase.expectedResult, comment)
	}
}

// Tests the less funcion on list of OrderByEntities
func (suite *TestCompoundReportSuite) TestLessByOrderByEntities(c *C) {
	testCases := []*struct {
		agentName []string
		targetIsp []string
		metricMax []int16
		expected bool
	} {
		{
			[]string{ "AG-1", "AG-2" },
			[]string{ "ISP-2", "ISP-1" },
			[]int16{ 10, 15 },
			true,
		},
		{
			[]string{ "AG-2", "AG-1" },
			[]string{ "ISP-2", "ISP-1" },
			[]int16{ 10, 15 },
			false,
		},
		{
			[]string{ "AG-2", "AG-2" },
			[]string{ "ISP-2", "ISP-1" },
			[]int16{ 10, 15 },
			true,
		},
		{
			[]string{ "AG-2", "AG-2" },
			[]string{ "ISP-2", "ISP-2" },
			[]int16{ 15, 10 },
			true,
		},
		{
			[]string{ "AG-2", "AG-2" },
			[]string{ "ISP-2", "ISP-2" },
			[]int16{ 10, 15 },
			false,
		},
	}

	sampleLessFunc := lessByOrderByEntities(
		[]*commonModel.OrderByEntity{
			{ "agent_name", utils.Ascending },
			{ "target_isp", utils.Descending },
			{ "max", utils.Descending },
		},
	)
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := sampleLessFunc.lessImpl(
			&model.DynamicRecord {
				Agent: &model.DynamicAgentProps {
					Name: testCase.agentName[0],
				},
				Target: &model.DynamicTargetProps {
					Isp: &owlModel.Isp {
						Name: testCase.targetIsp[0],
					},
				},
				Metrics: &model.DynamicMetrics {
					Metrics: &model.Metrics {
						Max: testCase.metricMax[0],
					},
				},
			},
			&model.DynamicRecord {
				Agent: &model.DynamicAgentProps {
					Name: testCase.agentName[1],
				},
				Target: &model.DynamicTargetProps {
					Isp: &owlModel.Isp {
						Name: testCase.targetIsp[1],
					},
				},
				Metrics: &model.DynamicMetrics {
					Metrics: &model.Metrics {
						Max: testCase.metricMax[1],
					},
				},
			},
		)

		c.Assert(testedResult, Equals, testCase.expected, comment)
	}
}

// Tests the filter of records
func (suite *TestCompoundReportSuite) TestFilterRecords(c *C) {
	testCases := []*struct {
		filter string
		expectedNumber int
	} {
		{ "", 3 },
		{ "$max >= 70", 3 },
		{ "$min > 25", 2 },
		{ "$max == 80 and $avg > 40", 1 },
		{ "$max == 70 or $avg < 25", 2 },
	}

	sampleRecords := []*model.DynamicRecord {
		{
			Metrics: &model.DynamicMetrics{
				Metrics: &model.Metrics {
					Max: 80, Min: 30, Avg: 45.59,
				},
			},
		},
		{
			Metrics: &model.DynamicMetrics{
				Metrics: &model.Metrics {
					Max: 70, Min: 30, Avg: 30.10,
				},
			},
		},
		{
			Metrics: &model.DynamicMetrics{
				Metrics: &model.Metrics {
					Max: 80, Min: 22, Avg: 22.33,
				},
			},
		},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := filterRecords(sampleRecords, testCase.filter)
		c.Assert(testedResult, HasLen, testCase.expectedNumber, comment)
	}
}

// Tests the retrieving of page(with sorted result)
func (suite *TestCompoundReportSuite) TestRetrievePage(c *C) {
	sampleRecords := []*model.DynamicRecord {
		{ Agent: &model.DynamicAgentProps{ Name: "AG-5" } },
		{ Agent: &model.DynamicAgentProps{ Name: "AG-4" } },
		{ Agent: &model.DynamicAgentProps{ Name: "AG-3" } },
		{ Agent: &model.DynamicAgentProps{ Name: "AG-2" } },
		{ Agent: &model.DynamicAgentProps{ Name: "AG-1" } },
	}

	samplePaging := &commonModel.Paging {
		Size: 2,
		Position: 2,
		OrderBy: []*commonModel.OrderByEntity {
			{ "agent_name", utils.DefaultDirection },
		},
	}

	testedResult := retrievePage(sampleRecords, samplePaging)
	c.Assert(testedResult[0].Agent.Name, Equals, "AG-3")
	c.Assert(testedResult[1].Agent.Name, Equals, "AG-4")
}

func (s *TestCompoundReportSuiteOnDb) SetUpTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestCompoundReportSuiteOnDb.TestBuildNqmDslByCompoundQuery":
		inTx(
			`
			INSERT INTO host(id, ip)
			VALUES(77621, '10.91.8.33')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES(1041, 77621, 'GC-01@10.91.8.33', 'eth3-gc-01@10.91.8.33', 'KCB-01.com.cn', x'0A5B0821')
			`,
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host)
			VALUES(2301, 'ZKP-01@abc.org', 'ZKP-TTC-33.easy.com.fr')
			`,
		)
	case "TestCompoundReportSuiteOnDb.TestToQueryDetail":
		inTx(
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(90801, 'gt-ab-1'), (90802, 'gt-ab-2'), (90803, 'gt-ab-3')
			`,
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3375, 'nt-ab-1'), (3376, 'nt-ab-2')
			`,
		)
	}
}
func (s *TestCompoundReportSuiteOnDb) TearDownTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestCompoundReportSuiteOnDb.TestBuildNqmDslByCompoundQuery":
		inTx(
			`DELETE FROM nqm_agent WHERE ag_id = 1041`,
			`DELETE FROM host WHERE id = 77621`,
			`DELETE FROM nqm_target WHERE tg_id = 2301`,
		)
	case "TestCompoundReportSuiteOnDb.TestToQueryDetail":
		inTx(
			`DELETE FROM owl_group_tag WHERE gt_id >= 90801 AND gt_id <= 90803`,
			`DELETE FROM owl_name_tag WHERE nt_id >= 3375 AND nt_id <= 3376`,
		)
	case "TestCompoundReportSuiteOnDb.TestGetCompoundQueryByUuid":
		inTx(
			fmt.Sprintf(`DELETE FROM owl_query WHERE qr_named_id = '%s'`, queryNamedId),
		)
	case "TestCompoundReportSuiteOnDb.TestBuildQuery":
		inTx(
			fmt.Sprintf(`DELETE FROM owl_query WHERE qr_named_id = '%s'`, queryNamedId),
		)
	}
}

func (s *TestCompoundReportSuiteOnDb) SetUpSuite(c *C) {
	s.dbTestSuite.SetUpSuite(c)
}
func (s *TestCompoundReportSuiteOnDb) TearDownSuite(c *C) {
	s.dbTestSuite.TearDownSuite(c)
}
