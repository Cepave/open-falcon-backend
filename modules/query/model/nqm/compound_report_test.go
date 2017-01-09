package nqm

import (
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	"github.com/Cepave/open-falcon-backend/common/utils"
	. "gopkg.in/check.v1"
)

type TestCompountReportSuite struct{}

var _ = Suite(&TestCompountReportSuite{})

// Tests the marshalling of JSON on metrics
func (suite *TestCompountReportSuite) TestMarshalJSONOnDynamicMetrics(c *C) {
	testCases := []*struct {
		columns []string
		expectedResult string
	} {
		{ // Everything
			[]string { MetricMax, MetricMin, MetricAvg, MetricMed, MetricMdev, MetricLoss, MetricCount, MetricPckSent, MetricPckReceived, MetricNumAgent, MetricNumTarget },
			`
			{
				"max": 78,
				"min": 21,
				"avg": 45.67,
				"med": 32,
				"mdev": 5.81,
				"loss": 0.04,
				"count": 100,
				"pck_sent": 2300,
				"pck_received": 2045,
				"num_agent": 10,
				"num_target": 15
			}
			`,
		},
		{ // Nothing
			[]string {},
			"{}",
		},
	}

	sampleMetrics := &DynamicMetrics {
		Metrics: &Metrics {
			Max: 78, Min: 21, Avg: 45.67, Med: 32, Mdev: 5.81, Loss: 0.04,
			Count: 100, NumberOfSentPackets: 2300, NumberOfReceivedPackets: 2045, NumberOfAgents: 10, NumberOfTargets: 15,
		},
	}
	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		sampleMetrics.Output = &testCase.columns

		c.Logf("Result JSON: %s", ojson.MarshalJSON(sampleMetrics))
		c.Assert(sampleMetrics, ocheck.JsonEquals, ojson.RawJsonForm(testCase.expectedResult), comment)
	}
}

// Tests the comparison of host of targets
func (suite *TestCompountReportSuite) TestCompareForHostOfTargets(c *C) {
	testCases := []*struct {
		leftHost string
		rightHost string
		expected int
	} {
		{ "22.20.30.40", "103.20.30.40", utils.SeqHigher },
		{ "10.20.30.40", "google.com", utils.SeqHigher },
		{ "10.20.30.40", "10.20.30.40", utils.SeqEqual },
		{ "wine.com.cn", "wine.com.cn", utils.SeqEqual },
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := CompareFunctions["target_host"](
			&DynamicRecord {
				Target: &DynamicTargetProps {
					Host: testCase.leftHost,
				},
			},
			&DynamicRecord {
				Target: &DynamicTargetProps {
					Host: testCase.rightHost,
				},
			},
			utils.Ascending,
		)
		c.Assert(testedResult, Equals, testCase.expected, comment)
	}
}
