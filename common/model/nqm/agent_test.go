package nqm

import (
	"github.com/leebenson/conform"
	testV "github.com/Cepave/open-falcon-backend/common/testing/validator"
	. "gopkg.in/check.v1"
	"reflect"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

// Tests the checking of group tags for two agents
func (suite *TestAgentSuite) TestAreGroupTagsSameOfAgentForAdding(c *C) {
	testCases := []struct {
		oldGroupTags []string
		newGroupTags []string
		expectedResult bool
	} {
		{
			[]string{}, []string{}, true,
		},
		{
			[]string{ "GT-1", "GT-2" },
			[]string{ "GT-2", "GT-1" },
			true,
		},
		{
			[]string{ "GT-1", "GT-2" },
			[]string{ "GT-2", "GT-1", "GT-3" },
			false,
		},
		{
			[]string{},
			[]string{ "GT-2", "GT-1", "GT-3" },
			false,
		},
		{
			[]string{ "GT-1", "GT-2" },
			[]string{},
			false,
		},
	}

	for i, testCase := range testCases {
		leftAgent := &AgentForAdding {
			GroupTags: testCase.oldGroupTags,
		}
		rightAgent := &AgentForAdding {
			GroupTags: testCase.newGroupTags,
		}

		c.Assert(leftAgent.AreGroupTagsSame(rightAgent), Equals, testCase.expectedResult, Commentf("Case: %d", i + 1))
	}
}

// Tests the unique of group tags
func (suite *TestAgentSuite) TestUniqueGroupTagsOfAgentForAdding(c *C) {
	testCases := []struct {
		sampleGroupTags []string
		expectedGroupTags []string
	} {
		{ []string{}, []string{} },
		{ []string{ "T1", "T2" }, []string{ "T1", "T2" } },
		{ []string{ "T1", "T2", "T1", "T2" }, []string{ "T1", "T2" } },
		{ []string{ "T3", "T3", "T4", "T4" }, []string{ "T3", "T4" } },
	}

	for _, testCase := range testCases {
		testAgent := &AgentForAdding {
			GroupTags: testCase.sampleGroupTags,
		}

		testAgent.UniqueGroupTags()

		c.Assert(testAgent.GroupTags, DeepEquals, testCase.expectedGroupTags)
	}
}

// Tests validation of NQM agent
func (suite *TestAgentSuite) TestConformOfAgentForAdding(c *C) {
	testedAgent := &AgentForAdding {
		Name: " name-1 ",
		ConnectionId: " conn-id-1 ",
		Comment: " comment-1 ",
		Hostname: " hostname-1 ",
		NameTagValue: " name-tag-1 ",
		GroupTags: []string{ " gt-1 ", " gt-2 " },
	}

	conform.Strings(testedAgent)

	c.Assert(testedAgent.Name, Equals, "name-1")
	c.Assert(testedAgent.ConnectionId, Equals, "conn-id-1")
	c.Assert(testedAgent.Comment, Equals, "comment-1")
	c.Assert(testedAgent.Hostname, Equals, "hostname-1")
	c.Assert(testedAgent.NameTagValue, Equals, "name-tag-1")
	c.Assert(testedAgent.GroupTags, DeepEquals, []string{ "gt-1", "gt-2" })
}

// Tests the data validation of AgentForAdding
func (suite *TestAgentSuite) TestValidateOfAgentForAdding(c *C) {
	testCases := []struct {
		fieldName string
		fieldValue interface{}
	} {
		{ "ConnectionId", "" },
		{ "Hostname", "" },
		{ "IspId", int16(0) },
		{ "ProvinceId", int16(0) },
		{ "CityId", int16(0) },
	}


	for _, testCase := range testCases {
		sampleAgent := &AgentForAdding{
			ConnectionId: "conn_id",
			Hostname: "hostname",
			IspId: -1,
			ProvinceId: -1,
			CityId: -1,
		}

		// Sets-up should-be-failed property
		reflect.ValueOf(sampleAgent).Elem().FieldByName(testCase.fieldName).
			Set(reflect.ValueOf(testCase.fieldValue))

		testV.AssertSingleErrorForField(
			c, Validator.Struct(sampleAgent),
			testCase.fieldName,
		)
	}
}
