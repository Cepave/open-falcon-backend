package owl

import (
	"reflect"

	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
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
		expected []*owlModel.Province
	}{
		{"北", []*owlModel.Province{{Id: 4, Name: "北京"}}},
		{"河", []*owlModel.Province{{Id: 3, Name: "河北"}, {Id: 19, Name: "河南"}}},
		{"幹", []*owlModel.Province{}},
	}

	for _, v := range testCases {
		got := GetProvincesByName(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			c.Error("Error:", got, "!=", v.expected)
		} else {
			c.Log(got, "==", v.expected)
		}
	}

	got := GetProvincesByName("")
	if len(got) >= 37 {
		c.Log("Case for \"\": PASS")
	} else {
		c.Error("Case for \"\": Checking len(got) >= 37...FAIL")
	}

}

func (suite *TestLocationSuite) TestGetCitiesByName(c *C) {
	testCases := []struct {
		input    string
		expected []*owlModel.City
	}{
		{"北", []*owlModel.City{{Id: 1, ProvinceId: 4, Name: "北京市", PostCode: "100000"}, {Id: 116, ProvinceId: 21, Name: "北海市", PostCode: "536000"}}},
		{"海口市", []*owlModel.City{{Id: 71, ProvinceId: 23, Name: "海口市", PostCode: "570000"}}},
		{"幹", []*owlModel.City{}},
	}

	for _, v := range testCases {
		got := GetCitiesByName(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			c.Error("Error:", got, "!=", v.expected)
		} else {
			c.Log(got, "==", v.expected)
		}
	}

	got := GetCitiesByName("")
	if len(got) >= 295 {
		c.Log("Case for \"\": PASS")
	} else {
		c.Error("Case for \"\": Checking len(got) >= 295...FAIL")
	}
}

func (s *TestLocationSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestLocationSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
