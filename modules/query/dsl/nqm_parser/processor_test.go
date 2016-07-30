package nqm_parser

import (
	. "gopkg.in/check.v1"
)

type TestParseProcessorSuite struct{}

var _ = Suite(&TestParseProcessorSuite{})

// Tests the parsing of Epoch time
func (suite *TestParseProcessorSuite) TestParseUnixTime(c *C) {
	testedResult, err := parseUnixTime(buildCurrent("1273053600"))

	c.Assert(err, IsNil)

	c.Logf("Parse UnixTime: %v", testedResult)

	c.Assert(testedResult.Unix(), Equals, int64(1273053600))
}

// Tests the parsing for node parameters
func (suite *TestParseProcessorSuite) TestNodeParams(c *C) {
	testCases := []*nodeParamsTestCase {
		&nodeParamsTestCase{ // Tests normal value of agent's property
			"agent.isp=i1,i2 agent.province=pv1,pv2 agent.city=ct1,ct2",
			func (testedQueryParam *QueryParams) {
				c.Assert(
					[]string{ "i1", "i2", "pv1", "pv2", "ct1", "ct2" },
					DeepEquals,
					[]string {
						testedQueryParam.AgentFilter.MatchIsps[0],
						testedQueryParam.AgentFilter.MatchIsps[1],
						testedQueryParam.AgentFilter.MatchProvinces[0],
						testedQueryParam.AgentFilter.MatchProvinces[1],
						testedQueryParam.AgentFilter.MatchCities[0],
						testedQueryParam.AgentFilter.MatchCities[1],
					},
				)
				c.Assert(testedQueryParam.IspRelation, Equals, UNKNOWN_RELATION)
				c.Assert(testedQueryParam.ProvinceRelation, Equals, UNKNOWN_RELATION)
				c.Assert(testedQueryParam.CityRelation, Equals, UNKNOWN_RELATION)
			},
		},
		&nodeParamsTestCase{ // Tests normal value of target's property
			"target.isp=i3,i4 target.province=pv3,pv4 target.city=ct3,ct4",
			func (testedQueryParam *QueryParams) {
				c.Assert(
					[]string{ "i3", "i4", "pv3", "pv4", "ct3", "ct4" },
					DeepEquals,
					[]string {
						testedQueryParam.TargetFilter.MatchIsps[0],
						testedQueryParam.TargetFilter.MatchIsps[1],
						testedQueryParam.TargetFilter.MatchProvinces[0],
						testedQueryParam.TargetFilter.MatchProvinces[1],
						testedQueryParam.TargetFilter.MatchCities[0],
						testedQueryParam.TargetFilter.MatchCities[1],
					},
				)
				c.Assert(testedQueryParam.IspRelation, Equals, UNKNOWN_RELATION)
				c.Assert(testedQueryParam.ProvinceRelation, Equals, UNKNOWN_RELATION)
				c.Assert(testedQueryParam.CityRelation, Equals, UNKNOWN_RELATION)
			},
		},
		&nodeParamsTestCase{ // Agent's auto-condition
			"agent.isp=%NOT_MATCH_ANOTHER% agent.province=%MATCH_ANOTHER% agent.city=%MATCH_ANOTHER%",
			func (testedQueryParam *QueryParams) {
				c.Assert(testedQueryParam.IspRelation, Equals, NOT_SAME_VALUE)
				c.Assert(testedQueryParam.ProvinceRelation, Equals, SAME_VALUE)
				c.Assert(testedQueryParam.CityRelation, Equals, SAME_VALUE)
			},
		},
		&nodeParamsTestCase{ // Agent's auto-condition
			"target.isp=%NOT_MATCH_ANOTHER% target.province=%MATCH_ANOTHER% target.city=%MATCH_ANOTHER%",
			func (testedQueryParam *QueryParams) {
				c.Assert(testedQueryParam.IspRelation, Equals, NOT_SAME_VALUE)
				c.Assert(testedQueryParam.ProvinceRelation, Equals, SAME_VALUE)
				c.Assert(testedQueryParam.CityRelation, Equals, SAME_VALUE)
			},
		},
		&nodeParamsTestCase{ // Duplicated condition
			"agent.isp=%NOT_MATCH_ANOTHER% target.isp=%MATCH_ANOTHER%",
			func (testedQueryParam *QueryParams) {
				c.Assert(testedQueryParam.IspRelation, Equals, SAME_VALUE)
			},
		},
	}

	for _, testCase := range testCases {
		c.Logf("Current DSL: %v", testCase.dsl)
		paramSetters, err := doParse(testCase.dsl)

		c.Assert(err, IsNil)

		var testedParams = NewQueryParams()
		testedParams.SetUpParams(paramSetters)

		testCase.assertionImpl(testedParams)
	}
}

// Tests the combination for string literals
type combineStringLiteralsTestCase struct {
	first string
	rest []interface{}
	expectedResult []string
}
func (suite *TestParseProcessorSuite) TestCombineStringLiterals(c *C) {
	testCases := []*combineStringLiteralsTestCase {
		&combineStringLiteralsTestCase{ "a1", []interface{}{ "b1", "b2" }, []string{ "a1", "b1", "b2" } },
		&combineStringLiteralsTestCase{ "a9", nil, []string{ "a9" } },
	}

	for _, testCase := range testCases {
		c.Assert(combineStringLiterals(testCase.first, testCase.rest), DeepEquals, testCase.expectedResult)
	}
}

// Tests the parsing for ISO8601 with various format of input
type parseIso8601Case struct {
	sampleValue string
	expectedYear int
	expectedHour int
}
func (suite *TestParseProcessorSuite) TestParseIso8601(c *C) {
	testCases := []parseIso8601Case {
		{ "2012-10-10T14:10+04:00", 2012, 14 },
		{ "2008-02-15T07:10", 2008, 7 },
		{ "2009-03-07T08", 2009, 8 },
		{ "2011-04-02", 2011, 0 },
	}

	for _, testCase := range testCases {
		testedResult, err := parseIso8601(buildCurrent(testCase.sampleValue))

		c.Assert(err, IsNil)

		c.Logf("Parse ISO8601: \"%v\" Result: \"%v\"", testCase.sampleValue, testedResult)
		c.Assert(testedResult.Year(), Equals, testCase.expectedYear)
		c.Assert(testedResult.Hour(), Equals, testCase.expectedHour)
	}
}

func buildCurrent(value string) *current {
	return &current {
		pos: position{},
		text: []byte(value),
	}
}
