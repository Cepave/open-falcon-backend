package owl

import (
	"reflect"

	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestLocationSuite struct{}

var _ = Suite(&TestLocationSuite{})

// Tests the check for city over hierarchy of administrative region
func (suite *TestLocationSuite) TestCheckHierarchyForCity(c *C) {
	testCases := []struct {
		provinceId int16
		cityId     int16
		hasError   bool
	}{
		{17, 27, false},
		{17, -1, false},
		{17, 20, true},
	}

	for i, testCase := range testCases {
		err := CheckHierarchyForCity(testCase.provinceId, testCase.cityId)

		if testCase.hasError {
			c.Logf("Error: %v", err)
			c.Assert(err, NotNil, Commentf("Test Case: [%d]", i+1))
		} else {
			c.Assert(err, IsNil, Commentf("Test Case: [%d]", i+1))
		}
	}
}

func (suite *TestLocationSuite) TestGetProvincesByName(c *C) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"北", []string{"北京"}},
		{"河", []string{"河北", "河南"}},
		{"幹", []string{}},
	}

	for _, v := range testCases {
		got := GetProvincesByName(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			c.Error("Error:", got, "!=", v.expected)
		} else {
			c.Log(got, "==", v.expected)
		}
	}
}

func (suite *TestLocationSuite) TestGetCitiesByName(c *C) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"北", []string{"北京市", "北海市"}},
		{"海口市", []string{"海口市"}},
		{"幹", []string{}},
	}

	for _, v := range testCases {
		got := GetCitiesByName(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			c.Error("Error:", got, "!=", v.expected)
		} else {
			c.Log(got, "==", v.expected)
		}
	}
}

func (s *TestLocationSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestLocationSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
