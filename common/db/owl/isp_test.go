package owl

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestIspSuite struct{}

var _ = Suite(&TestIspSuite{})

func (suite *TestIspSuite) TestGetISPByName(c *C) {
	testCases := []*struct {
		input    string
		expected []*owlModel.Isp
	}{
		{"北", []*owlModel.Isp{{Id: 1, Name: "北京三信时代", Acronym: "BJCIII"}, {Id: 13, Name: "北京宽捷", Acronym: "KJNET"}}},
		{"方", []*owlModel.Isp{{Id: 8, Name: "方正宽带", Acronym: "FBN"}}},
		{"幹", []*owlModel.Isp{}},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)
		c.Assert(GetISPsByName(testCase.input), DeepEquals, testCase.expected, comment)
	}

	/**
	 * Tests the query with empty string
	 */
	expectedAllIsps := GetISPsByName("")
	c.Assert(len(expectedAllIsps) >= 32, Equals, true, Commentf("Expected number of ISPs is not 32 at least"))
	// :~)
}

// Tests the getting of ISP by id
func (suite *TestIspSuite) TestGetIspById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasData  bool
	}{
		{1, true},
		{-10, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := GetIspById(testCase.sampleId)

		if testCase.hasData {
			c.Assert(testedResult, NotNil, comment)
		} else {
			c.Assert(testedResult, IsNil, comment)
		}
	}
}

func (s *TestIspSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestIspSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
