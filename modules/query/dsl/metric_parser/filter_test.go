package metric_parser

import (
	"github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	. "gopkg.in/check.v1"
)

type TestFilterSuite struct{}

var _ = Suite(&TestFilterSuite{})

// Tests the general filter
func (suite *TestFilterSuite) TestFilterImpl(c *C) {
	testCases := []*struct {
		testedFilter *filterImpl
		expectedResult bool
	} {
		/**
		 * Asserts every operator
		 */
		{ &filterImpl { MetricMax, ">", float64(88) }, true },
		{ &filterImpl { MetricMax, ">", float64(89) }, false },
		{ &filterImpl { MetricMax, "<", float64(90) }, true },
		{ &filterImpl { MetricMax, "<", float64(89) }, false },
		{ &filterImpl { MetricMax, ">=", float64(88) }, true },
		{ &filterImpl { MetricMax, ">=", float64(89) }, true },
		{ &filterImpl { MetricMax, ">=", float64(90) }, false },
		{ &filterImpl { MetricMax, "<=", float64(90) }, true },
		{ &filterImpl { MetricMax, "<=", float64(89) }, true },
		{ &filterImpl { MetricMax, "<=", float64(88) }, false },
		{ &filterImpl { MetricMax, "==", float64(89) }, true },
		{ &filterImpl { MetricMax, "==", float64(90) }, false },
		{ &filterImpl { MetricMax, "==", float64(88) }, false },
		{ &filterImpl { MetricMax, "!=", float64(88) }, true },
		{ &filterImpl { MetricMax, "!=", float64(90) }, true },
		{ &filterImpl { MetricMax, "!=", float64(89) }, false },
		// :~)
		/**
		 * Asserts every metrics
		 */
		{ &filterImpl { MetricMin, "==", float64(13) }, true },
		{ &filterImpl { MetricAvg, "==", float64(56.12) }, true },
		{ &filterImpl { MetricMed, "==", float64(45) }, true },
		{ &filterImpl { MetricMdev, "==", float64(6.2) }, true },
		{ &filterImpl { MetricLoss, "==", float64(0.031) }, true },
		{ &filterImpl { MetricCount, "==", float64(80) }, true },
		{ &filterImpl { MetricNumAgent, "==", float64(40) }, true },
		{ &filterImpl { MetricNumTarget, "==", float64(38) }, true },
		{ &filterImpl { MetricPckSent, "==", float64(3000) }, true },
		{ &filterImpl { MetricPckReceived, "==", float64(2870) }, true },
		// :~)
		/**
		 * Asserts for two metrics
		 */
		{ &filterImpl { MetricPckReceived, "<", MetricPckSent }, true },
		{ &filterImpl { MetricPckReceived, ">", MetricPckSent }, false },
		// :~)
	}

	sampleMetrics := &nqm.Metrics {
		Max: 89, Min: 13, Avg: 56.12,
		Med: 45, Mdev: 6.2, Loss: 0.031,
		Count: 80, NumberOfAgents: 40, NumberOfTargets: 38,
		NumberOfSentPackets: 3000, NumberOfReceivedPackets: 2870,
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d.", i + 1)

		testedResult := testCase.testedFilter.IsMatch(sampleMetrics)
		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the filter of boolean
func (suite *TestFilterSuite) TestBoolFilterImpl(c *C) {
	testCases := []*struct {
		boolOperator bool
		firstFactor bool
		restFactors []bool
		expectedResult bool
	} {
		{ // Or conditions
			true, true, []bool{}, true,
		},
		{ // Or conditions
			true, false, []bool{}, false,
		},
		{ // Or conditions
			true, true, []bool{ true }, true,
		},
		{ // Or conditions
			true, false, []bool{ true }, true,
		},
		{ // Or conditions
			true, false, []bool{ false, false }, false,
		},
		{ // And conditions
			false, true, []bool{}, true,
		},
		{ // And conditions
			false, false, []bool{}, false,
		},
		{ // And conditions
			false, true, []bool{ true }, true,
		},
		{ // And conditions
			false, false, []bool{ true, true }, false,
		},
		{ // And conditions
			false, false, []bool{ false }, false,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedFilter := newBoolFilterImpl(
			testCase.boolOperator,
			fakeFilter(testCase.firstFactor),
			buildFilters(testCase.restFactors...),
		)

		c.Assert(testedFilter.IsMatch(nil), Equals, testCase.expectedResult, comment)
	}
}

type fakeFilter bool
func (f fakeFilter) IsMatch(metrics *nqm.Metrics) bool {
	return bool(f)
}

func buildFilters(values ...bool) []interface{} {
	filters := make([]interface{}, 0)

	for _, v := range values {
		filters = append(filters, fakeFilter(v))
	}

	return filters
}
