package types

import (
	"fmt"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"

	. "github.com/Cepave/open-falcon-backend/common/reflect/types"
	. "gopkg.in/check.v1"
)

type TestBindingSuite struct{}

var _ = Suite(&TestBindingSuite{})

type gBox struct {
	name string
}
func (b *gBox) Bind(sourceObject interface{}) {
	switch typedV := sourceObject.(type) {
	case string:
		b.name = "string: " + typedV
	case int32:
		b.name = fmt.Sprintf("int: %d", typedV)
	default:
		panic("Nothing")
	}
}

// Tests the translation from binding to converter
func (suite *TestBindingSuite) TestBindingToConverter(c *C) {
	testCases := []*struct {
		sourceObj interface{}
		expectedName string
	} {
		{ "Easy", "string: Easy" },
		{ int32(91), "int: 91" },
	}

	converter, targetType := BindingToConverter(
		func() interface{} {
			return &gBox{}
		},
	)
	srv := NewDefaultConversionService()
	srv.AddConverter(TypeOfString, targetType, converter)
	srv.AddConverter(TypeOfInt32, targetType, converter)

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := srv.ConvertTo(testCase.sourceObj, targetType).(*gBox)

		c.Assert(testedValue.name, Equals, testCase.expectedName, comment)
	}
}
