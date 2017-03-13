package nqm

import (
	"github.com/Cepave/open-falcon-backend/common/conform"
	testV "github.com/Cepave/open-falcon-backend/common/testing/validator"
	. "gopkg.in/check.v1"
	"reflect"
)

type TestTargetSuite struct{}

var _ = Suite(&TestTargetSuite{})

// Tests validation of NQM target
func (suite *TestTargetSuite) TestConformOfTargetForAdding(c *C) {
	testedTarget := &TargetForAdding {
		Name: " name-1 ",
		Host: " host-1 ",
		Comment: " comment-1 ",
		NameTagValue: " name-tag-1 ",
		GroupTags: []string{ " gt-1 ", " gt-2 " },
	}

	conform.MustConform(testedTarget)

	c.Assert(testedTarget.Name, Equals, "name-1")
	c.Assert(testedTarget.Host, Equals, "host-1")
	c.Assert(testedTarget.Comment, Equals, "comment-1")
	c.Assert(testedTarget.NameTagValue, Equals, "name-tag-1")
	c.Assert(testedTarget.GroupTags, DeepEquals, []string{ "gt-1", "gt-2" })
}

// Tests the data validation of TargetForAdding
func (suite *TestTargetSuite) TestValidateOfTargetForAdding(c *C) {
	testCases := []*struct {
		fieldName string
		fieldValue interface{}
	} {
		{ "Name", "" },
		{ "Host", "" },
		{ "IspId", int16(0) },
		{ "ProvinceId", int16(0) },
		{ "CityId", int16(0) },
	}

	for _, testCase := range testCases {
		sampleTarget := &TargetForAdding{
			Name: "conn_id",
			Host: "hostname",
			IspId: -1,
			ProvinceId: -1,
			CityId: -1,
		}

		// Sets-up should-be-failed property
		reflect.ValueOf(sampleTarget).Elem().FieldByName(testCase.fieldName).
			Set(reflect.ValueOf(testCase.fieldValue))

		testV.AssertSingleErrorForField(
			c, Validator.Struct(sampleTarget),
			testCase.fieldName,
		)
	}
}
