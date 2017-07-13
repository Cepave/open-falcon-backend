package nqm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/satori/go.uuid"

	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	"github.com/Cepave/open-falcon-backend/common/utils"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"

	. "gopkg.in/check.v1"
)

type TestCompoundReportSuite struct{}
type TestCompoundReportSuiteOnDb struct {
	*dbTestSuite
}

var (
	_ = Suite(&TestCompoundReportSuite{})
	_ = Suite(&TestCompoundReportSuiteOnDb{&dbTestSuite{}})
)

// Tests the conversion of query to detail information
func (suite *TestCompoundReportSuiteOnDb) TestToQueryDetail(c *C) {
	/**
	 * Sets-up sample query
	 */
	sampleQuery := model.NewCompoundQuery()

	sampleQuery.Filters.Metrics = "$max >= $min"
	sampleQuery.Grouping.Agent = []string{model.AgentGroupingHostname, model.GroupingIsp}
	sampleQuery.Grouping.Target = []string{model.TargetGroupingHost, model.GroupingProvince, model.GroupingIsp}
	sampleQuery.Output.Metrics = []string{model.MetricAvg, model.MetricMax, model.MetricLoss}

	agentFilter := sampleQuery.Filters.Agent

	agentFilter.Name = []string{"Cool1", "Cool2"}
	agentFilter.Hostname = []string{"gc1.com", "gc2.com"}
	agentFilter.IpAddress = []string{"123.71.1", "123.71.2"}
	agentFilter.ConnectionId = []string{"c01", "c02"}
	agentFilter.IspIds = []int16{2, 3, 4, 1, 5, 6}
	agentFilter.ProvinceIds = []int16{3, 4, 7, 8}
	agentFilter.CityIds = []int16{13, 14}
	agentFilter.NameTagIds = []int16{3375, 3376}
	agentFilter.GroupTagIds = []int32{90801, 90802, 90803}

	targetFilter := sampleQuery.Filters.Target

	targetFilter.Name = []string{"Zoo-1", "Zoo-2"}
	targetFilter.Host = []string{"kz1.com", "kz2.com"}
	targetFilter.IspIds = []int16{5, 7}
	targetFilter.ProvinceIds = []int16{1, 3, 8}
	targetFilter.CityIds = []int16{15, 16}
	targetFilter.NameTagIds = []int16{3375, 3376}
	targetFilter.GroupTagIds = []int32{90801, 90802, 90803}
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

// Tests the conversion of query to deatil information on special conditions
//
// Some properties(isp_ids, province_ids, city_ids, name_tag_ids) supports special values:
//
// -11 - The property(ISP, location, etc.) should be same with another side
// -12 - The property(ISP, location, etc.) should not be same with another side
func (suite *TestCompoundReportSuiteOnDb) TestToQueryDetailOnSpecialValue(c *C) {
	/**
	 * Sets-up sample query
	 */
	sampleQuery := model.NewCompoundQuery()

	agentFilter := sampleQuery.Filters.Agent

	agentFilter.IspIds = []int16{-11, -12}
	agentFilter.ProvinceIds = []int16{-11, -12}
	agentFilter.CityIds = []int16{-11, -12}
	agentFilter.NameTagIds = []int16{-11, -12}

	targetFilter := sampleQuery.Filters.Target
	targetFilter.IspIds = []int16{-11, -12}
	targetFilter.ProvinceIds = []int16{-11, -12}
	targetFilter.CityIds = []int16{-11, -12}
	targetFilter.NameTagIds = []int16{-11, -12}
	// :~)

	testedDetail := ToQueryDetail(sampleQuery)

	c.Logf("%#v", testedDetail.Agent)
	c.Logf("%#v", testedDetail.Target)

	c.Assert(testedDetail.Agent.Isps[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Agent.Isps[1].Id, Equals, int16(-12))
	c.Assert(testedDetail.Agent.Provinces[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Agent.Provinces[1].Id, Equals, int16(-12))
	c.Assert(testedDetail.Agent.Cities[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Agent.Cities[1].Id, Equals, int16(-12))
	c.Assert(testedDetail.Agent.NameTags[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Agent.NameTags[1].Id, Equals, int16(-12))

	c.Assert(testedDetail.Target.Isps[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Target.Isps[1].Id, Equals, int16(-12))
	c.Assert(testedDetail.Target.Provinces[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Target.Provinces[1].Id, Equals, int16(-12))
	c.Assert(testedDetail.Target.Cities[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Target.Cities[1].Id, Equals, int16(-12))
	c.Assert(testedDetail.Target.NameTags[0].Id, Equals, int16(-11))
	c.Assert(testedDetail.Target.NameTags[1].Id, Equals, int16(-12))
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
	type assertFunc func(v interface{}, comment CommentInterface)

	testCases := []*struct {
		queryJson        string
		assertProperties map[string]interface{}
	}{
		{ // Fetch viable nodes
			`{
				"filters": {
					"agent": {
						"name": [ "GC-01" ],
						"connection_id": [ "eth3-gc-01" ],
						"hostname": [ "KCB-01" ],
						"ip_address": [ "10.91.8.33" ],
						"isp_ids": [ 10, 20 ],
						"province_ids": [ 32, 31 ],
						"city_ids": [ 65, 72 ],
						"name_tag_ids": [ 101, 114 ],
						"group_tag_ids": [ 1291, 1309 ]
					},
					"target": {
						"name": [ "ZKP-01" ],
						"host": [ "ZKP-TTC-33" ],
						"isp_ids": [ 10, 20 ],
						"province_ids": [ 32, 31 ],
						"city_ids": [ 65, 72 ],
						"name_tag_ids": [ 101, 114 ],
						"group_tag_ids": [ 1291, 1309 ]
					}
				}
			}`,
			map[string]interface{}{
				"IdsOfAgents":          []int32{1041},
				"IdsOfAgentIsps":       []int16{10, 20},
				"IdsOfAgentProvinces":  []int16{31, 32},
				"IdsOfAgentCities":     []int16{65, 72},
				"IdsOfAgentNameTags":   []int16{101, 114},
				"IdsOfAgentGroupTags":  []int32{1291, 1309},
				"IdsOfTargets":         []int32{2301},
				"IdsOfTargetIsps":      []int16{10, 20},
				"IdsOfTargetProvinces": []int16{31, 32},
				"IdsOfTargetCities":    []int16{65, 72},
				"IdsOfTargetNameTags":  []int16{101, 114},
				"IdsOfTargetGroupTags": []int32{1291, 1309},
			},
		},
		{ // Fetch non-viable nodes
			`{
				"filters": {
					"agent": {
						"name": [ "no-node" ],
						"connection_id": [ "no-node" ],
						"hostname": [ "no-node" ],
						"ip_address": [ "10.20.31.41" ]
					},
					"target": {
						"name": [ "no-node" ],
						"host": [ "no-node" ]
					}
				}
			}`,
			map[string]interface{}{
				"IdsOfAgents":          []int32{-2},
				"IdsOfAgentIsps":       []int16{},
				"IdsOfAgentProvinces":  []int16{},
				"IdsOfAgentCities":     []int16{},
				"IdsOfAgentNameTags":   []int16{},
				"IdsOfAgentGroupTags":  []int32{},
				"IdsOfTargets":         []int32{-2},
				"IdsOfTargetIsps":      []int16{},
				"IdsOfTargetProvinces": []int16{},
				"IdsOfTargetCities":    []int16{},
				"IdsOfTargetNameTags":  []int16{},
				"IdsOfTargetGroupTags": []int32{},
				"IspRelation":          model.NoCondition,
				"ProvinceRelation":     model.NoCondition,
				"CityRelation":         model.NoCondition,
				"NameTagRelation":      model.NoCondition,
			},
		},
		{ // Same relation(by agent)
			`{
				"filters": {
					"agent": {
						"isp_ids": [ -11 ],
						"province_ids": [ -11 ],
						"city_ids": [ -11 ],
						"name_tag_ids": [ -11 ],
						"group_tag_ids": [ -11 ]
					}
				}
			}`,
			map[string]interface{}{
				"IspRelation":      model.SameValue,
				"ProvinceRelation": model.SameValue,
				"CityRelation":     model.SameValue,
				"NameTagRelation":  model.SameValue,
			},
		},
		{ // Not same relation(by target)
			`{
				"filters": {
					"target": {
						"isp_ids": [ -12 ],
						"province_ids": [ -12 ],
						"city_ids": [ -12 ],
						"name_tag_ids": [ -12 ],
						"group_tag_ids": [ -12 ]
					}
				}
			}`,
			map[string]interface{}{
				"IspRelation":      model.NotSameValue,
				"ProvinceRelation": model.NotSameValue,
				"CityRelation":     model.NotSameValue,
				"NameTagRelation":  model.NotSameValue,
			},
		},
		{ // Fetch absolute time
			`{
				"filters": {
					"time": { "start_time": 70089020, "end_time": 70389020 }
				}
			}`,
			map[string]interface{}{
				"StartTime": toPointerOfEpochTime(70089020),
				"EndTime":   toPointerOfEpochTime(70389020),
			},
		},
		{ // Fetch relative time
			`
			{
				"filters": {
					"time": {
						"to_now": { "unit": "d", "value": 5 }
					}
				}
			}
			`,
			map[string]interface{}{
				"StartTime": assertFunc(func(v interface{}, comment CommentInterface) {
					timeValue := time.Unix(int64(*(v.(*EpochTime))), 0)
					c.Logf("Start Time(Relative): %s", timeValue.Format(time.RFC3339))

					c.Assert(timeValue, ocheck.TimeBefore, time.Now().AddDate(0, 0, -4), comment)
					c.Assert(timeValue, ocheck.TimeAfter, time.Now().AddDate(0, 0, -6), comment)
				}),
				"EndTime": assertFunc(func(v interface{}, comment CommentInterface) {
					timeValue := time.Unix(int64(*(v.(*EpochTime))), 0)
					c.Logf("End Time(Relative): %s", timeValue.Format(time.RFC3339))

					c.Assert(timeValue, ocheck.TimeBefore, time.Now(), comment)
					c.Assert(timeValue, ocheck.TimeAfter, time.Now().AddDate(0, 0, -2), comment)
				}),
			},
		},
		{ // Fetch relative time(multiple time ranges)
			`
			{
				"filters": {
					"time": {
						"to_now": { "unit": "d", "value": 3, "start_time_of_day": "03:45", "end_time_of_day": "04:55" }
					}
				}
			}
			`,
			map[string]interface{}{
				"TimeRanges": assertFunc(func(v interface{}, comment CommentInterface) {
					timeRanges := v.([]*TimeRangeOfDsl)
					c.Assert(timeRanges, HasLen, 3, comment)
				}),
			},
		},
	}

	var newSampleQuery = func(jsonQuery string) *model.CompoundQuery {
		query := model.NewCompoundQuery()
		c.Assert(query.UnmarshalJSON([]byte(jsonQuery)), IsNil)
		query.SetupDefault()
		return query
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedDsl := buildNqmDslByCompoundQuery(
			newSampleQuery(testCase.queryJson),
		)

		for name, value := range testCase.assertProperties {
			testedProperty := reflect.ValueOf(testedDsl).Elem().
				FieldByName(name)

			comment := Commentf("[Case: %s] Check Property: %s.", comment.CheckCommentString(), name)

			if assertImpl, ok := value.(assertFunc); ok {
				assertImpl(testedProperty.Interface(), comment)
				continue
			}

			c.Assert(testedProperty.Interface(), DeepEquals, value, comment)
		}
	}
}

// Tests the building of group columns for DSL
func (suite *TestCompoundReportSuite) TestBuildGroupingColumnOfDsl(c *C) {
	testCases := []*struct {
		sampleGrouping *model.QueryGrouping
		expected       []string
	}{
		{ // Agent + Target
			&model.QueryGrouping{
				Agent:  []string{model.AgentGroupingName, model.AgentGroupingHostname, model.AgentGroupingIpAddress},
				Target: []string{model.TargetGroupingName, model.TargetGroupingHost},
			},
			[]string{"ag_id", "tg_id"},
		},
		{ // Other properties
			&model.QueryGrouping{
				Agent:  []string{model.GroupingIsp, model.GroupingProvince, model.GroupingCity, model.GroupingNameTag},
				Target: []string{model.GroupingIsp, model.GroupingProvince, model.GroupingCity, model.GroupingNameTag},
			},
			[]string{
				"ag_isp_id", "ag_pv_id", "ag_ct_id", "ag_nt_id",
				"tg_isp_id", "tg_pv_id", "tg_ct_id", "tg_nt_id",
			},
		},
		{ // Agent + other property of target
			&model.QueryGrouping{
				Agent:  []string{model.AgentGroupingName, model.GroupingIsp},
				Target: []string{model.GroupingIsp},
			},
			[]string{"ag_id", "tg_isp_id"},
		},
		{ // Other property of agent + target
			&model.QueryGrouping{
				Agent:  []string{model.GroupingNameTag},
				Target: []string{model.TargetGroupingName, model.GroupingIsp},
			},
			[]string{"ag_nt_id", "tg_id"},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedGrouping := buildGroupingColumnOfDsl(testCase.sampleGrouping)
		c.Assert(testedGrouping, DeepEquals, testCase.expected, comment)
	}
}

// Tests the setup of sorting properties
func (suite *TestCompoundReportSuite) TestSetupSorting(c *C) {
	testCases := []*struct {
		sampleEntities []*commonModel.OrderByEntity
		outputMetrics  []string
		expectedResult []*commonModel.OrderByEntity
	}{
		{
			[]*commonModel.OrderByEntity{{"agent_isp", commonModel.Descending}},
			[]string{"max", "min"},
			[]*commonModel.OrderByEntity{{"agent_isp", commonModel.Descending}, {"loss", commonModel.Descending}},
		},
		{
			[]*commonModel.OrderByEntity{},
			[]string{"max", "min"},
			[]*commonModel.OrderByEntity{
				{"max", commonModel.Descending},
				{"min", commonModel.Descending},
				{"loss", commonModel.Descending},
			},
		},
		{
			[]*commonModel.OrderByEntity{},
			[]string{"max", "avg", "loss"},
			[]*commonModel.OrderByEntity{
				{"avg", commonModel.Descending},
				{"loss", commonModel.Descending},
			},
		},
		{
			[]*commonModel.OrderByEntity{},
			[]string{"num_agent", "num_target"},
			[]*commonModel.OrderByEntity{
				{"num_agent", commonModel.Descending},
				{"num_target", commonModel.Descending},
				{"loss", commonModel.Descending},
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
		expected  bool
	}{
		{
			[]string{"AG-1", "AG-2"},
			[]string{"ISP-2", "ISP-1"},
			[]int16{10, 15},
			true,
		},
		{
			[]string{"AG-2", "AG-1"},
			[]string{"ISP-2", "ISP-1"},
			[]int16{10, 15},
			false,
		},
		{
			[]string{"AG-2", "AG-2"},
			[]string{"ISP-2", "ISP-1"},
			[]int16{10, 15},
			true,
		},
		{
			[]string{"AG-2", "AG-2"},
			[]string{"ISP-2", "ISP-2"},
			[]int16{15, 10},
			true,
		},
		{
			[]string{"AG-2", "AG-2"},
			[]string{"ISP-2", "ISP-2"},
			[]int16{10, 15},
			false,
		},
	}

	sampleLessFunc := lessByOrderByEntities(
		[]*commonModel.OrderByEntity{
			{"agent_name", utils.Ascending},
			{"target_isp", utils.Descending},
			{"max", utils.Descending},
		},
	)
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := sampleLessFunc.lessImpl(
			&model.DynamicRecord{
				Agent: &model.DynamicAgentProps{
					Name: &testCase.agentName[0],
				},
				Target: &model.DynamicTargetProps{
					Isp: &owlModel.Isp{
						Name: testCase.targetIsp[0],
					},
				},
				Metrics: &model.DynamicMetrics{
					Metrics: &model.Metrics{
						Max: testCase.metricMax[0],
					},
				},
			},
			&model.DynamicRecord{
				Agent: &model.DynamicAgentProps{
					Name: &testCase.agentName[1],
				},
				Target: &model.DynamicTargetProps{
					Isp: &owlModel.Isp{
						Name: testCase.targetIsp[1],
					},
				},
				Metrics: &model.DynamicMetrics{
					Metrics: &model.Metrics{
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
		filter         string
		expectedNumber int
	}{
		{"", 3},
		{"$max >= 70", 3},
		{"$min > 25", 2},
		{"$max == 80 and $avg > 40", 1},
		{"$max == 70 or $avg < 25", 2},
	}

	sampleRecords := []*model.DynamicRecord{
		{
			Metrics: &model.DynamicMetrics{
				Metrics: &model.Metrics{
					Max: 80, Min: 30, Avg: 45.59,
				},
			},
		},
		{
			Metrics: &model.DynamicMetrics{
				Metrics: &model.Metrics{
					Max: 70, Min: 30, Avg: 30.10,
				},
			},
		},
		{
			Metrics: &model.DynamicMetrics{
				Metrics: &model.Metrics{
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
	sp := func(v string) *string { return &v }

	sampleRecords := []*model.DynamicRecord{
		{Agent: &model.DynamicAgentProps{Name: sp("AG-5")}},
		{Agent: &model.DynamicAgentProps{Name: sp("AG-4")}},
		{Agent: &model.DynamicAgentProps{Name: sp("AG-3")}},
		{Agent: &model.DynamicAgentProps{Name: sp("AG-2")}},
		{Agent: &model.DynamicAgentProps{Name: sp("AG-1")}},
	}

	samplePaging := &commonModel.Paging{
		Size:     2,
		Position: 2,
		OrderBy: []*commonModel.OrderByEntity{
			{"agent_name", utils.DefaultDirection},
		},
	}

	testedResult := retrievePage(sampleRecords, samplePaging)
	c.Assert(*testedResult[0].Agent.Name, Equals, "AG-3")
	c.Assert(*testedResult[1].Agent.Name, Equals, "AG-4")
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
