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
		expected []*city1view
	}{
		{"北", []*city1view{&city1view{Id: 1, Name: "北京市", PostCode: "100000", Province: &owlModel.Province{Id: 4, Name: "北京"}}, &city1view{Id: 116, Name: "北海市", PostCode: "536000", Province: &owlModel.Province{Id: 21, Name: "广西"}}}},
		{"海口市", []*city1view{&city1view{Id: 71, Name: "海口市", PostCode: "570000", Province: &owlModel.Province{Id: 23, Name: "海南"}}}},
		{"幹", []*city1view{}},
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

func (suite *TestLocationSuite) TestGetCitiesInProvinceByName(c *C) {
	testCases := []struct {
		inputPvId int
		inputName string
		expected  []*owlModel.City2
	}{
		{24, "广", []*owlModel.City2{{Id: 74, Name: "广元市", PostCode: "628000"}, {Id: 81, Name: "广安市", PostCode: "638000"}}},
		{20, "茂名市", []*owlModel.City2{{Id: 20, Name: "茂名市", PostCode: "525000"}}},
		{20, "幹", []*owlModel.City2{}},
		{0, "", []*owlModel.City2{}},
		{2, "", []*owlModel.City2{}},
	}

	for i, v := range testCases {
		if i > 3 {
			break
		}
		got := GetCitiesInProvinceByName(v.inputPvId, v.inputName)
		if !reflect.DeepEqual(got, v.expected) {
			c.Error("Error:", got, "!=", v.expected)
		} else {
			c.Log(got, "==", v.expected)
		}
	}

	got := GetCitiesInProvinceByName(testCases[4].inputPvId, testCases[4].inputName)
	if len(got) >= 11 {
		c.Log("Case ct_pv_id=2, ct_name=\"\": PASS")
	} else {
		c.Error("Case ct_pv_id=2, ct_name=\"\": Checking len(got) >= 11...FAIL")
	}
}

func (s *TestLocationSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestLocationSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
