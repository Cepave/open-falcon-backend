package nqm

import (
	"reflect"

	"github.com/Cepave/open-falcon-backend/common/conform"

	otest "github.com/Cepave/open-falcon-backend/common/testing"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	testV "github.com/Cepave/open-falcon-backend/common/testing/validator"
	"github.com/Cepave/open-falcon-backend/common/utils"
	. "gopkg.in/check.v1"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

// Tests validation of NQM agent
func (suite *TestAgentSuite) TestConformOfAgentForAdding(c *C) {
	ps := func(v string) *string { return &v }

	testedAgent := &AgentForAdding{
		Name:         utils.PointerOfCloneString(" name-1 "),
		ConnectionId: " conn-id-1 ",
		Comment:      utils.PointerOfCloneString(" comment-1 "),
		Hostname:     " hostname-1 ",
		NameTagValue: ps(" name-tag-1 "),
		GroupTags:    []string{" gt-1 ", " gt-2 "},
	}

	conform.MustConform(testedAgent)

	c.Assert(testedAgent.Name, DeepEquals, utils.PointerOfCloneString("name-1"))
	c.Assert(testedAgent.ConnectionId, Equals, "conn-id-1")
	c.Assert(testedAgent.Comment, DeepEquals, utils.PointerOfCloneString("comment-1"))
	c.Assert(testedAgent.Hostname, Equals, "hostname-1")
	c.Assert(testedAgent.NameTagValue, DeepEquals, ps("name-tag-1"))
	c.Assert(testedAgent.GroupTags, DeepEquals, []string{"gt-1", "gt-2"})
}

// Tests the data validation of AgentForAdding
func (suite *TestAgentSuite) TestValidateOfAgentForAdding(c *C) {
	testCases := []*struct {
		fieldName  string
		fieldValue interface{}
	}{
		{"ConnectionId", ""},
		{"Hostname", ""},
		{"IspId", int16(0)},
		{"ProvinceId", int16(0)},
		{"CityId", int16(0)},
	}

	for _, testCase := range testCases {
		ocheck.LogTestCase(c, testCase)

		sampleAgent := &AgentForAdding{
			ConnectionId: "conn_id",
			Hostname:     "hostname",
			IspId:        -1,
			ProvinceId:   -1,
			CityId:       -1,
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

// Tests the getting of duration of time
func (suite *TestAgentSuite) TestGetDurationOfLastAccessOnPingListLog(c *C) {
	testCases := []*struct {
		checkedTime     string
		expectedMinutes int64
	}{
		{"2014-06-07T08:13:07+08:00", 13},
		{"2014-06-07T10:00:33+08:00", 120},
	}

	accessTime := otest.ParseTime(c, "2014-06-07T08:00:00+08:00")
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedLog := &PingListLog{AccessTime: accessTime}
		checkedTime := otest.ParseTime(c, testCase.checkedTime)

		c.Assert(testedLog.GetDurationOfLastAccess(checkedTime), Equals, testCase.expectedMinutes, comment)
	}
}
