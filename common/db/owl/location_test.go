package owl

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestLocationSuite struct{}

var _ = Suite(&TestLocationSuite{})

// Tests the check for city over hierarchy of administrative region
func (suite *TestLocationSuite) TestCheckHierarchyForCity(c *C) {
	testCases := []*struct {
		provinceId int16
		cityId     int16
		hasError   bool
	}{
		{17, 27, false},
		{17, -1, false},
		{17, 20, true},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case : [%d]", i+1)
		err := CheckHierarchyForCity(testCase.provinceId, testCase.cityId)

		if testCase.hasError {
			c.Logf("Error: %v", err)
			c.Assert(err, NotNil, comment)
		} else {
			c.Assert(err, IsNil, comment)
		}
	}
}

// Tests the loading of province by id
func (suite *TestLocationSuite) TestGetProviceById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	}{
		{3, true},
		{-3, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := GetProvinceById(testCase.sampleId)

		if testCase.hasFound {
			c.Assert(testedResult, NotNil, comment)
		} else {
			c.Assert(testedResult, IsNil, comment)
		}
	}
}

func (suite *TestLocationSuite) TestGetProvincesByName(c *C) {
	testCases := []*struct {
		input    string
		expected []*owlModel.Province
	}{
		{"北", []*owlModel.Province{{Id: 4, Name: "北京"}}},
		{"河", []*owlModel.Province{{Id: 3, Name: "河北"}, {Id: 19, Name: "河南"}}},
		{"幹", []*owlModel.Province{}},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case : [%d]", i+1)
		c.Assert(GetProvincesByName(testCase.input), DeepEquals, testCase.expected, comment)
	}

	c.Assert(len(GetProvincesByName("")), ocheck.LargerThanOrEqualTo, 37, Commentf("Needs 37 provinces at least"))
}

func (suite *TestLocationSuite) TestGetCitiesByName(c *C) {
	testCases := []*struct {
		input    string
		expected []*city1view
	}{
		{"北",
			[]*city1view{
				{Id: 1, Name: "北京市", PostCode: "100000", Province: &owlModel.Province{Id: 4, Name: "北京"}},
				{Id: 116, Name: "北海市", PostCode: "536000", Province: &owlModel.Province{Id: 21, Name: "广西"}},
			},
		},
		{"海口市",
			[]*city1view{
				{Id: 71, Name: "海口市", PostCode: "570000", Province: &owlModel.Province{Id: 23, Name: "海南"}},
			},
		},
		{"幹", []*city1view{}},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case : [%d]", i+1)

		c.Assert(GetCitiesByName(testCase.input), DeepEquals, testCase.expected, comment)
	}

	c.Assert(len(GetCitiesByName("")), ocheck.LargerThanOrEqualTo, 295, Commentf("Needs 295 provinces at least"))
}

func (suite *TestLocationSuite) TestGetCitiesInProvinceByName(c *C) {
	testCases := []*struct {
		inputPvId int
		inputName string
		expected  []*owlModel.City2
	}{
		{24, "广", []*owlModel.City2{{Id: 74, Name: "广元市", PostCode: "628000"}, {Id: 81, Name: "广安市", PostCode: "638000"}}},
		{20, "茂名市", []*owlModel.City2{{Id: 20, Name: "茂名市", PostCode: "525000"}}},
		{20, "幹", []*owlModel.City2{}},
		{0, "", []*owlModel.City2{}},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case : [%d]", i+1)

		c.Assert(GetCitiesInProvinceByName(testCase.inputPvId, testCase.inputName), DeepEquals, testCase.expected, comment)
	}

	// Asserts the loading of cities of a province without any query condition
	c.Assert(len(GetCitiesInProvinceByName(2, "")), ocheck.LargerThanOrEqualTo, 11, Commentf("Needs 11 cities at least"))
}

// Tests the getting of City1 by id
func (suite *TestLocationSuite) TestGetCityById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	}{
		{3, true},
		{-3, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := GetCityById(testCase.sampleId)

		if testCase.hasFound {
			c.Assert(testedResult, NotNil, comment)
		} else {
			c.Assert(testedResult, IsNil, comment)
		}
	}
}

// Tests the getting of City2 by id
func (suite *TestLocationSuite) TestGetCity2ById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	}{
		{3, true},
		{-3, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := GetCity2ById(testCase.sampleId)

		if testCase.hasFound {
			c.Assert(testedResult, NotNil, comment)
		} else {
			c.Assert(testedResult, IsNil, comment)
		}
	}
}

// Tests getting of cities by name prefix
func (suite *TestLocationSuite) TestGetCit2sByName(c *C) {
	testCases := []*struct {
		sampleNamePrefix string
		leastNumber      int
	}{
		{"大", 3},
		{"NO!!", 0},
		{"", 295},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := GetCity2sByName(testCase.sampleNamePrefix)

		c.Assert(
			len(testedResult), ocheck.LargerThanOrEqualTo, testCase.leastNumber,
			comment,
		)
	}
}

func (s *TestLocationSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestLocationSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
