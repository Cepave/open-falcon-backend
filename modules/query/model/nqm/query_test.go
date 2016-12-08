package nqm

import (
	sjson "github.com/bitly/go-simplejson"
	"github.com/Cepave/open-falcon-backend/common/utils"
	. "gopkg.in/check.v1"
	"reflect"
	"time"
)

type TestQuerySuite struct{}

var _ = Suite(&TestQuerySuite{})

// Tests the loading of filters.metrics
func (suite *TestQuerySuite) TestLoadMetricsOfFilters(c *C) {
	testCases := []struct {
		sampleJson string
		expectedResult string
	} {
		{
			`{ "filters": { "metrics": " $mAx > 20 aNd $min < 40 " } }`,
			"$max > 20 and $min < 40",
		},
		{
			`{ "filters": { "metrics": "" } }`, "",
		},
		{
			`{}`, "",
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedQuery := loadQueryObject(c, testCase.sampleJson, comment)
		c.Assert(testedQuery.Filters.Metrics, DeepEquals, testCase.expectedResult, comment)
	}
}

// Tests the loading of filters.time
func (suite *TestQuerySuite) TestLoadFiltersOfTime(c *C) {
	testCases := []struct {
		jsonSource string
		expectedResult *TimeFilter
	} {
		{
			`
			{
				"filters": {
					"time": {
						"start_time": 8977123,
						"end_time": 19082711
					}
				}
			}
			`,
			&TimeFilter {
				TimeRangeType: TIME_RANGE_ABSOLUTE,
				StartTime: time.Unix(8977123, 0),
				EndTime: time.Unix(19082711, 0),
			},
		},
		{
			`
			{
				"filters": {
					"time": {
						"to_now": {
							"unit": "m",
							"value": 3
						}
					}
				}
			}
			`,
			&TimeFilter {
				TimeRangeType: TIME_RANGE_RELATIVE,
				ToNow: &TimeWithUnit{
					Unit: TU_MONTH,
					Value: 3,
				},
			},
		},
		{
			` { "filters": {} } `, &TimeFilter{ TimeRangeType: 0 },
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := loadQueryObject(c, testCase.jsonSource, comment).Filters.Time
		expectedResult := testCase.expectedResult

		c.Assert(testedResult.TimeRangeType, Equals, expectedResult.TimeRangeType, comment)
		switch testedResult.TimeRangeType {
		case TIME_RANGE_ABSOLUTE:
			c.Assert(testedResult.StartTime, Equals, expectedResult.StartTime, comment)
			c.Assert(testedResult.EndTime, Equals, expectedResult.EndTime, comment)
		case TIME_RANGE_RELATIVE:
			testedTimeWithUnit := testedResult.ToNow
			c.Assert(testedTimeWithUnit.Unit, Equals, expectedResult.ToNow.Unit, comment)
			c.Assert(testedTimeWithUnit.Value, Equals, expectedResult.ToNow.Value, comment)
		}
	}
}

// Tests the loading of filters.agent
func (suite *TestQuerySuite) TestLoadFiltersOfAgent(c *C) {
	testCases := []struct {
		jsonSource string
		expectedFilter *AgentFilter
	} {
		{
			`
			{ "filters": {
				"agent": {
					"name": [ "G1", "C2", "K3", "G1"],
					"hostname": [ "hs-3", "hs-1", "hs-3" ],
					"ip_address": [ "10.20", "9.7", "10.20" ],
					"connection_id": [ "conn-id-3", "conn-id-1", "conn-id-3" ],
					"isp_ids": [ 20, 17, 20 ],
					"province_ids": [ 31, 22, 31 ],
					"city_ids": [ 32, 7, 32 ],
					"name_tag_ids": [ 77, 9, 77 ],
					"group_tag_ids": [ 16, 8, 16 ]
				}
			} }
			`,
			&AgentFilter {
				Name: []string { "C2", "G1", "K3" },
				Hostname: []string{ "hs-1", "hs-3" },
				IpAddress: []string{ "10.20", "9.7" },
				ConnectionId: []string{ "conn-id-1", "conn-id-3" },
				IspIds: []int16{ 17, 20 },
				ProvinceIds: []int16{ 22, 31 },
				CityIds: []int16{ 7, 32 },
				NameTagIds: []int16{ 9, 77 },
				GroupTagIds: []int32{ 8, 16 },
			},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := loadQueryObject(c, testCase.jsonSource, comment)
		c.Assert(testedResult.Filters.Agent, DeepEquals, testCase.expectedFilter, comment)
	}
}

// Tests the loading of filters.target
func (suite *TestQuerySuite) TestLoadFiltersOfTarget(c *C) {
	testCases := []struct {
		jsonSource string
		expectedFilter *TargetFilter
	} {
		{
			`
			{ "filters": {
				"target": {
					"name": [ "G1", "C2", "K3", "G1"],
					"host": [ "hs-3", "hs-1", "hs-3" ],
					"isp_ids": [ 20, 17, 20 ],
					"province_ids": [ 31, 22, 31 ],
					"city_ids": [ 32, 7, 32 ],
					"name_tag_ids": [ 77, 9, 77 ],
					"group_tag_ids": [ 16, 8, 16 ]
				}
			} }
			`,
			&TargetFilter {
				Name: []string { "C2", "G1", "K3" },
				Host: []string{ "hs-1", "hs-3" },
				IspIds: []int16{ 17, 20 },
				ProvinceIds: []int16{ 22, 31 },
				CityIds: []int16{ 7, 32 },
				NameTagIds: []int16{ 9, 77 },
				GroupTagIds: []int32{ 8, 16 },
			},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := loadQueryObject(c, testCase.jsonSource, comment)
		c.Assert(testedResult.Filters.Target, DeepEquals, testCase.expectedFilter, comment)
	}
}

// Tests the query of loading output
func (suite *TestQuerySuite) TestLoadOutput(c *C) {
	testCases := []*struct{
		sampleJson string
		expectedResult []string
	} {
		{
			`{ "output": { "metrics": [ "max", "min", "avg" ] } }`,
			[]string { "max", "min", "avg" },
		},
		{
			`{ "output": { "metrics": [] } }`,
			[]string {},
		},
		{ // No output property
			`{}`,
			[]string {},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedQuery := loadQueryObject(c, testCase.sampleJson, comment)

		c.Assert(testedQuery.Output.Metrics, DeepEquals, testCase.expectedResult, comment)
	}
}

// Tests the loading of grouping
func (suite *TestQuerySuite) TestLoadGrouping(c *C) {
	testCases := []struct {
		sampleJson string
		expectedAgentGrouping []string
		expectedTargetGrouping []string
	} {
		{
			`{ "grouping": { "agent": [ "isp", "province" ], "target": [ "name_tag"] } }`,
			[]string{ "isp", "province" },
			[]string{ "name_tag" },
		},
		{
			`{ "grouping": { "agent": [], "target": [] } }`,
			[]string{},
			[]string{},
		},
		{ // No output property
			`{}`,
			[]string{},
			[]string{},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedQuery := loadQueryObject(c, testCase.sampleJson, comment)

		c.Assert(testedQuery.Grouping.Agent, DeepEquals, testCase.expectedAgentGrouping, comment)
		c.Assert(testedQuery.Grouping.Target, DeepEquals, testCase.expectedTargetGrouping, comment)
	}
}

// Tests the purifying of json's array of numbers
func (suite *TestQuerySuite) TestPurifyNumberArrayOfJson(c *C) {
	testCases := []struct {
		jsonSource string
		targetType reflect.Type
		expectedResult interface{}
	} {
		{
			`[ 38, 29, 40, 38, 29 ]`,
			utils.TypeOfInt8,
			[]int8 { 29, 38, 40 },
		},
		{
			`[ 78781, 10981, 78781 ]`,
			utils.TypeOfUint64,
			[]uint64 { 10981, 78781 },
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := purifyNumberArrayOfJson(
			loadJson(c, testCase.jsonSource),
			testCase.targetType,
		)
		c.Assert(testedResult, DeepEquals, testCase.expectedResult, comment)
	}
}

// Tests the purifying for array of strings(by doamin)
func (suite *TestQuerySuite) TestPurifyStringArrayOfJsonForDomain(c *C) {
	testCases := []struct {
		jsonSource string
		sampleDomain map[string]bool
		expectedData []string
	} {
		{
			`[ "c9", " A1 ", "  ", "no-1", "C1", "c2", "no-3", "A1", " C1 " ]`,
			map[string]bool { "a1": true, "c1": true, "c9": true },
			[]string{ "c9", "a1", "c1" },
		},
		{ `[ "A1", "A2" ]`, map[string]bool {}, []string{}, },
		{ `[ "A1", "A2" ]`, nil, []string{}, },
		{ `[ "  ", "" ]`, map[string]bool { "k1": true }, []string{}, },
		{ "", map[string]bool { "g1": true }, nil, },
		{ "", nil, nil, },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := purifyStringArrayOfJsonForDomain(
			loadJson(c, testCase.jsonSource),
			testCase.sampleDomain,
		)

		c.Assert(testedResult, DeepEquals, testCase.expectedData, comment)
	}
}

// Tests the purifying for array of strings(for values)
func (suite *TestQuerySuite) TestPurifyStringArrayOfJsonForValues(c *C) {
	testCases := []struct {
		jsonSource string
		exepctedResult []string
	} {
		{
			`[ " A1 ", " b1 ", "A1", "B2", " a1 ", "B2", "C3" ]`,
			[]string{ "A1", "B2", "C3", "a1", "b1" },
		},
		{ `[ "  ", ""]`, []string{}, },
		{ "", nil, },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := purifyStringArrayOfJsonForValues(
			loadJson(c, testCase.jsonSource),
		)
		c.Assert(testedResult, DeepEquals, testCase.exepctedResult, comment)
	}
}

func loadJson(c *C, jsonSource string) *sjson.Json {
	if jsonSource == "" {
		return nil
	}

	sampleJson, err := sjson.NewJson(([]byte)(jsonSource))
	c.Assert(err, IsNil)

	return sampleJson
}

func loadQueryObject(c *C, json string, comment CommentInterface) *CompoundQuery {
	testedQuery := NewCompoundQuery()
	err := testedQuery.UnmarshalJSON(([]byte)(json))
	c.Assert(err, IsNil, comment)

	return testedQuery
}
