package metric_parser

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	. "gopkg.in/check.v1"
)

type TestMetricGrammarSuite struct{}

var _ = Suite(&TestMetricGrammarSuite{})

// Tests the grammar
func (suite *TestMetricGrammarSuite) TestParse(c *C) {
	testCases := []*struct {
		code string
		expectedResult bool
	} {
		{ "$min >= 30", true },
		{ "($min >= 30)", true },
		{ "(($min >= 30))", true },
		{ "( $min < 30 )", false },
		{ "(( $min < 30 ))", false },
		{ "(( (( $min < 30 )) ))", false },
		{ "$min < 30 or $max == 80", true },
		{ "$min < 30 and $max == 80", false },
		{ "$max > $min and $max != 80 or $med >= 50", false },
		{ "$max > $min and ($max == 80 or $med >= 50)", true },
		{ "$max > $min or $max != 80 and $med >= 50", true },
		{ "$max > $min or ($max == 80 and $med >= 50)", true },
		{ "4 == 5 or 4 > 5 or 9 < 8", false },
		{ "((\t\t\t($max > 70)) and ( $max >   $min   and   $max > $avg )   and\t$avg >\t30 and $med > 30\t\t\t)", true },
	}

	sampleMetrics := &nqm.Metrics {
		Min: 30, Max: 80, Avg: 50,
		Med: 45, Mdev: 5.6,
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		/**
		 * Parses the code to filter
		 */
		c.Logf("Curret Code: %s", testCase.code)
		result, err := Parse("Metrics", ([]byte)(testCase.code))
		c.Assert(err, IsNil, comment)
		resultFilter := result.(nqm.MetricFilter)
		// :~)

		c.Assert(resultFilter.IsMatch(sampleMetrics), Equals, testCase.expectedResult, comment)
	}
}

// Tests the error for factor
func (suite *TestMetricGrammarSuite) TestError(c *C) {
	testCases := []*struct {
		code string
		matchError string
	} {
		{ "30 == $minor", ".*Unknown factor.*" },
		{ "30a99 == 30", ".*Unknown factor.*" },
		{ "30 == $max98", ".*Unknown factor.*" },
		{ "avg2 == 39", ".*Unknown factor.*" },
		{ "30 > ", ".*Need right.*" },
		{ "!= 33", ".*Unknown factor.*" },
		{ ">", ".*Unknown factor.*" },
		{ "($min > 20", ".*Need right.*" },
		{ "$min > 20)", ".*no match.*" },
		{ "$min > 20 or", ".*OR what.*" },
		{ "$min > 20 and", ".*AND what.*" },
		{ "$min ==!", ".*Unknown operator.*" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		_, err := doParse(testCase.code)

		combinedError := combineError(c, err)
		c.Assert(combinedError, ErrorMatches, testCase.matchError, comment)
	}
}

func combineError(c *C, err error) error {
	c.Assert(err, NotNil)

	list := err.(errList)
	errorMessage := ""
	for _, err := range list {
		c.Logf("|Got error for factor| --> %v", err)
		errorMessage += err.(*parserError).Error()
	}

	return fmt.Errorf("%s", errorMessage)
}
func doParse(code string) (interface{}, error) {
	return Parse("Metrics", ([]byte)(code))
}
