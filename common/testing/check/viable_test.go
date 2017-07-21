package check

import (
	. "gopkg.in/check.v1"
)

type TestViableSuite struct{}

var _ = Suite(&TestViableSuite{})

type car struct{}

// Tests the Viable check for various type of objects
func (suite *TestViableSuite) TestViableValue(c *C) {
	testCases := []*struct {
		sampleValue interface{}
		expected    bool
	}{
		{&car{}, true},
		{(*car)(nil), false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult, _ := ViableValue.Check(
			[]interface{}{testCase.sampleValue, true},
			[]string{},
		)

		c.Assert(testedResult, Equals, testCase.expected, comment)
	}
}
