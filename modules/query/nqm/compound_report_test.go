package nqm

import (
	"reflect"
	"time"

	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	"github.com/Cepave/open-falcon-backend/common/types"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/satori/go.uuid"

	db "github.com/Cepave/open-falcon-backend/modules/query/database"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	ch "gopkg.in/check.v1"
)

type mockQueryService struct {
	sampleQuery *owlModel.Query
}

func (s *mockQueryService) LoadQueryByUuid(targetUuid uuid.UUID) *owlModel.Query {
	return s.sampleQuery
}
func (s *mockQueryService) CreateOrLoadQuery(query *owlModel.Query) {
	s.sampleQuery = query
}

var _ = Describe("Query object", func() {
	mockQueryService := &mockQueryService{}

	BeforeEach(func() {
		db.QueryObjectService = mockQueryService
	})
	AfterEach(func() {
		db.QueryObjectService = nil
	})

	Context("Build a new one", func() {
		sampleNqmQuery := buildSampleNqmQuery(
			"ab01", "ab02", "cd99",
		)

		It("The named id, content, and md5 digest should be set to exepcted one", func() {
			testedQuery := BuildQuery(sampleNqmQuery)
			expectedMd5Content := new(types.Bytes16)
			expectedMd5Content.FromVarBytes(sampleNqmQuery.GetDigestValue())

			Expect(testedQuery).To(PointTo(
				MatchFields(IgnoreExtras, Fields{
					"NamedId":    Equal(queryNamedId),
					"Content":    BeEquivalentTo(sampleNqmQuery.GetCompressedQuery()),
					"Md5Content": BeEquivalentTo(*expectedMd5Content),
				}),
			))
		})
	})

	Context("Load by UUID", func() {
		Context("Load viable query", func() {
			sampleNqmQuery := buildSampleNqmQuery(
				"pc01-com", "pc02-com", "srv-01-com",
			)

			BeforeEach(func() {
				mockQueryService.sampleQuery = &owlModel.Query{
					Content: sampleNqmQuery.GetCompressedQuery(),
				}
			})

			It("Hostname should be matched", func() {
				testedQuery := GetCompoundQueryByUuid(uuid.NewV4())
				Expect(testedQuery.Filters.Agent.Hostname).To(
					Equal(sampleNqmQuery.Filters.Agent.Hostname),
				)
			})
		})

		Context("Load nil value", func() {
			BeforeEach(func() {
				mockQueryService.sampleQuery = nil
			})

			It("Object should be nil", func() {
				testedQuery := GetCompoundQueryByUuid(uuid.NewV4())
				Expect(testedQuery).To(BeNil())
			})
		})
	})
})

func buildSampleNqmQuery(hostname ...string) *model.CompoundQuery {
	sampleNqmQuery := model.NewCompoundQuery()
	sampleNqmQuery.Filters.Agent.Hostname = hostname
	sampleNqmQuery.SetupDefault()
	return sampleNqmQuery
}

type TestCompoundReportSuite struct{}
type TestCompoundReportSuiteOnDb struct {
	*dbTestSuite
}

var (
	_ = ch.Suite(&TestCompoundReportSuite{})
	_ = ch.Suite(&TestCompoundReportSuiteOnDb{&dbTestSuite{}})
)

// Tests the conversion of query to detail information
func (suite *TestCompoundReportSuiteOnDb) TestToQueryDetail(c *ch.C) {
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

	c.Assert(string(testedDetail.Metrics), ch.Equals, sampleQuery.Filters.Metrics)

	/**
	 * Asserts the query detail on agent conditions
	 */
	testedAgentDetail := testedDetail.Agent
	c.Assert(testedAgentDetail.Name, ch.DeepEquals, agentFilter.Name)
	c.Assert(testedAgentDetail.Hostname, ch.DeepEquals, agentFilter.Hostname)
	c.Assert(testedAgentDetail.IpAddress, ch.DeepEquals, agentFilter.IpAddress)
	c.Assert(testedAgentDetail.ConnectionId, ch.DeepEquals, agentFilter.ConnectionId)
	c.Assert(testedAgentDetail.Isps, ch.HasLen, 6)
	c.Assert(testedAgentDetail.Provinces, ch.HasLen, 4)
	c.Assert(testedAgentDetail.Cities, ch.HasLen, 2)
	c.Assert(testedAgentDetail.NameTags, ch.HasLen, 2)
	c.Assert(testedAgentDetail.GroupTags, ch.HasLen, 3)
	// :~)

	/**
	 * Asserts the query detail on target conditions
	 */
	testedTargetDetail := testedDetail.Target
	c.Assert(testedTargetDetail.Name, ch.DeepEquals, targetFilter.Name)
	c.Assert(testedTargetDetail.Host, ch.DeepEquals, targetFilter.Host)
	c.Assert(testedTargetDetail.Isps, ch.HasLen, 2)
	c.Assert(testedTargetDetail.Provinces, ch.HasLen, 3)
	c.Assert(testedTargetDetail.Cities, ch.HasLen, 2)
	c.Assert(testedTargetDetail.NameTags, ch.HasLen, 2)
	c.Assert(testedTargetDetail.GroupTags, ch.HasLen, 3)
	// :~)

	/**
	 * Asserts the output detail
	 */
	testedOutputDetail := testedDetail.Output
	c.Assert(testedOutputDetail.Agent, ch.HasLen, 2)
	c.Assert(testedOutputDetail.Target, ch.HasLen, 3)
	c.Assert(testedOutputDetail.Metrics, ch.HasLen, 3)
	// :~)
}

// Tests the conversion of query to deatil information on special conditions
//
// Some properties(isp_ids, province_ids, city_ids, name_tag_ids) supports special values:
//
// -11 - The property(ISP, location, etc.) should be same with another side
// -12 - The property(ISP, location, etc.) should not be same with another side
func (suite *TestCompoundReportSuiteOnDb) TestToQueryDetailOnSpecialValue(c *ch.C) {
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

	c.Assert(testedDetail.Agent.Isps[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Agent.Isps[1].Id, ch.Equals, int16(-12))
	c.Assert(testedDetail.Agent.Provinces[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Agent.Provinces[1].Id, ch.Equals, int16(-12))
	c.Assert(testedDetail.Agent.Cities[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Agent.Cities[1].Id, ch.Equals, int16(-12))
	c.Assert(testedDetail.Agent.NameTags[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Agent.NameTags[1].Id, ch.Equals, int16(-12))

	c.Assert(testedDetail.Target.Isps[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Target.Isps[1].Id, ch.Equals, int16(-12))
	c.Assert(testedDetail.Target.Provinces[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Target.Provinces[1].Id, ch.Equals, int16(-12))
	c.Assert(testedDetail.Target.Cities[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Target.Cities[1].Id, ch.Equals, int16(-12))
	c.Assert(testedDetail.Target.NameTags[0].Id, ch.Equals, int16(-11))
	c.Assert(testedDetail.Target.NameTags[1].Id, ch.Equals, int16(-12))
}

// Tests the building of NQM dsl by compound query
func (suite *TestCompoundReportSuiteOnDb) TestBuildNqmDslByCompoundQuery(c *ch.C) {
	type assertFunc func(v interface{}, comment ch.CommentInterface)

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
				"StartTime": assertFunc(func(v interface{}, comment ch.CommentInterface) {
					timeValue := time.Unix(int64(*(v.(*EpochTime))), 0)
					c.Logf("Start Time(Relative): %s", timeValue.Format(time.RFC3339))

					c.Assert(timeValue, ocheck.TimeBefore, time.Now().AddDate(0, 0, -4), comment)
					c.Assert(timeValue, ocheck.TimeAfter, time.Now().AddDate(0, 0, -6), comment)
				}),
				"EndTime": assertFunc(func(v interface{}, comment ch.CommentInterface) {
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
				"TimeRanges": assertFunc(func(v interface{}, comment ch.CommentInterface) {
					timeRanges := v.([]*TimeRangeOfDsl)
					c.Assert(timeRanges, ch.HasLen, 3, comment)
				}),
			},
		},
	}

	var newSampleQuery = func(jsonQuery string) *model.CompoundQuery {
		query := model.NewCompoundQuery()
		c.Assert(query.UnmarshalJSON([]byte(jsonQuery)), ch.IsNil)
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

			comment := ch.Commentf("[Case: %s] Check Property: %s.", comment.CheckCommentString(), name)

			if assertImpl, ok := value.(assertFunc); ok {
				assertImpl(testedProperty.Interface(), comment)
				continue
			}

			c.Assert(testedProperty.Interface(), ch.DeepEquals, value, comment)
		}
	}
}

// Tests the building of group columns for DSL
func (suite *TestCompoundReportSuite) TestBuildGroupingColumnOfDsl(c *ch.C) {
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
		comment := ch.Commentf("Test Case: %d", i+1)

		testedGrouping := buildGroupingColumnOfDsl(testCase.sampleGrouping)
		c.Assert(testedGrouping, ch.DeepEquals, testCase.expected, comment)
	}
}

// Tests the setup of sorting properties
func (suite *TestCompoundReportSuite) TestSetupSorting(c *ch.C) {
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

		c.Assert(samplePaging.OrderBy, ch.DeepEquals, testCase.expectedResult, comment)
	}
}

// Tests the less funcion on list of OrderByEntities
func (suite *TestCompoundReportSuite) TestLessByOrderByEntities(c *ch.C) {
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

		c.Assert(testedResult, ch.Equals, testCase.expected, comment)
	}
}

// Tests the filter of records
func (suite *TestCompoundReportSuite) TestFilterRecords(c *ch.C) {
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
		c.Assert(testedResult, ch.HasLen, testCase.expectedNumber, comment)
	}
}

// Tests the retrieving of page(with sorted result)
func (suite *TestCompoundReportSuite) TestRetrievePage(c *ch.C) {
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
	c.Assert(*testedResult[0].Agent.Name, ch.Equals, "AG-3")
	c.Assert(*testedResult[1].Agent.Name, ch.Equals, "AG-4")
}

func (s *TestCompoundReportSuiteOnDb) SetUpTest(c *ch.C) {
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
func (s *TestCompoundReportSuiteOnDb) TearDownTest(c *ch.C) {
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
	}
}

func (s *TestCompoundReportSuiteOnDb) SetUpSuite(c *ch.C) {
	s.dbTestSuite.SetUpSuite(c)
}
func (s *TestCompoundReportSuiteOnDb) TearDownSuite(c *ch.C) {
	s.dbTestSuite.TearDownSuite(c)
}
