package nqm_parser

import (
	. "gopkg.in/check.v1"
	"sort"
	"time"
)

type TestQueryParamsSuite struct{}

var _ = Suite(&TestQueryParamsSuite{})

// Tests the error for both "province" and "city" are set
type locationErrorTestCase struct {
	params     *QueryParams
	matchError string
	hasError   bool
}

func (suite *TestQueryParamsSuite) TestCheckParamsWithLocationError(c *C) {
	testCases := []*locationErrorTestCase{
		{
			&QueryParams{
				AgentFilter: NodeFilter{
					MatchProvinces: []string{"p1", "p2"},
					MatchCities:    []string{"c1", "c2"},
				},
			},
			".*Agent.*", true,
		},
		{
			&QueryParams{
				TargetFilter: NodeFilter{
					MatchProvinces: []string{"p1", "p2"},
					MatchCities:    []string{"c1", "c2"},
				},
			},
			".*Target.*", true,
		},
	}

	for _, testCase := range testCases {
		err := testCase.params.checkParams()

		if testCase.hasError {
			c.Logf("Error for check parameters: %v", err)
			c.Assert(err, ErrorMatches, testCase.matchError)
		} else {
			c.Assert(err, IsNil)
		}
	}
}

// Tests the elimination for duplicated values
func (suite *TestQueryParamsSuite) TestCheckParamsWithDuplicatedValue(c *C) {
	testedParams_1 := &QueryParams{
		AgentFilter: NodeFilter{
			MatchProvinces: []string{"p1", "p2", "p1", "p2"},
			MatchIsps:      []string{"i1", "i2", "i1", "i2"},
		},
		TargetFilter: NodeFilter{
			MatchProvinces: []string{"p3", "p4", "p3", "p4"},
			MatchIsps:      []string{"i3", "i4", "i3", "i4"},
		},
	}
	testedParams_2 := &QueryParams{
		AgentFilter: NodeFilter{
			MatchCities: []string{"c1", "c2", "c1", "c2"},
		},
		TargetFilter: NodeFilter{
			MatchCities: []string{"c3", "c4", "c3", "c4"},
		},
	}

	c.Assert(testedParams_1.checkParams(), IsNil)
	c.Assert(testedParams_2.checkParams(), IsNil)

	sort.Strings(testedParams_1.AgentFilter.MatchProvinces)
	sort.Strings(testedParams_2.AgentFilter.MatchCities)
	sort.Strings(testedParams_1.AgentFilter.MatchIsps)
	sort.Strings(testedParams_1.TargetFilter.MatchProvinces)
	sort.Strings(testedParams_2.TargetFilter.MatchCities)
	sort.Strings(testedParams_1.TargetFilter.MatchIsps)

	c.Assert(testedParams_1.AgentFilter.MatchProvinces, DeepEquals, []string{"p1", "p2"})
	c.Assert(testedParams_2.AgentFilter.MatchCities, DeepEquals, []string{"c1", "c2"})
	c.Assert(testedParams_1.AgentFilter.MatchIsps, DeepEquals, []string{"i1", "i2"})
	c.Assert(testedParams_1.TargetFilter.MatchProvinces, DeepEquals, []string{"p3", "p4"})
	c.Assert(testedParams_2.TargetFilter.MatchCities, DeepEquals, []string{"c3", "c4"})
	c.Assert(testedParams_1.TargetFilter.MatchIsps, DeepEquals, []string{"i3", "i4"})
}

type checkRationalOfParametersTestCase struct {
	startTime time.Time
	endTime   time.Time
	assertion func(error)
}

// Tests the rational meaning of query parameters
func (suite *TestQueryParamsSuite) TestCheckRationalOfParameters(c *C) {
	testCases := []checkRationalOfParametersTestCase{
		/**
		 * Nothing failed
		 */
		{
			startTime: time.Now(),
			endTime:   time.Now().Add(1 * time.Minute),
			assertion: func(err error) {
				c.Assert(err, IsNil)
			},
		},
		// :~)
		/**
		 * The time range is meaningless
		 */
		{
			startTime: time.Now(),
			endTime:   time.Now().Add(-1 * time.Minute),
			assertion: func(err error) {
				c.Assert(err, ErrorMatches, "Start time is not valid.*")
			},
		},
		// :~)
	}

	for _, testCase := range testCases {
		testedParams := &QueryParams{
			StartTime: testCase.startTime,
			EndTime:   testCase.endTime,
		}

		testCase.assertion(testedParams.CheckRationalOfParameters())
	}
}
