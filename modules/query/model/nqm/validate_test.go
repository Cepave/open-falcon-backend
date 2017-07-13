package nqm

import (
	"gopkg.in/go-playground/validator.v9"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestValidateSuite struct{}

var _ = Suite(&TestValidateSuite{})

// Tests the validation of time unit object
func (suite *TestValidateSuite) TestValidateTimeWithUnit(c *C) {
	sPtr := func(v string) *string { return &v }

	testCases := []*struct {
		sampleTimeUnit *TimeWithUnit
		hasError       bool
	}{
		{
			&TimeWithUnit{},
			false,
		},
		{
			&TimeWithUnit{Unit: TimeUnitDay, StartTimeOfDay: sPtr("02:00"), EndTimeOfDay: sPtr("04:00")},
			false,
		},
		{
			&TimeWithUnit{StartTimeOfDay: sPtr("02:00")},
			true,
		},
		{
			&TimeWithUnit{EndTimeOfDay: sPtr("02:00")},
			true,
		},
		{
			&TimeWithUnit{StartTimeOfDay: sPtr("0g:00"), EndTimeOfDay: sPtr("02:00")},
			true,
		},
		{
			&TimeWithUnit{Unit: TimeUnitHour, StartTimeOfDay: sPtr("02:00"), EndTimeOfDay: sPtr("04:00")},
			true,
		},
	}

	validator := validator.New()
	validator.RegisterStructValidation(
		ValidateTimeWithUnit, TimeWithUnit{},
	)

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		result := validator.Struct(testCase.sampleTimeUnit)

		if testCase.hasError {
			c.Logf("Validation Error: %v", result)
			c.Assert(result, NotNil, comment)
		} else {
			c.Assert(result, IsNil, comment)
		}
	}
}
