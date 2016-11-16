package owl

import (
	"reflect"

	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestIspSuite struct{}

var _ = Suite(&TestIspSuite{})

func (suite *TestIspSuite) TestGetISPByName(c *C) {
	testCases := []struct {
		input    string
		expected []*owlModel.Isp
	}{
		{"北", []*owlModel.Isp{{Id: 1, Name: "北京三信时代"}, {Id: 13, Name: "北京宽捷"}}},
		{"方", []*owlModel.Isp{{Id: 8, Name: "方正宽带"}}},
		{"幹", []*owlModel.Isp{}},
	}

	for _, v := range testCases {
		got := GetISPsByName(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			c.Error("Error:", got, "!=", v.expected)
		} else {
			c.Log(got, "==", v.expected)
		}
	}

	got := GetISPsByName("")
	if len(got) >= 32 {
		c.Log("Case for \"\": PASS")
	} else {
		c.Error("Case for \"\": Checking len(got) >= 32...FAIL")
	}
}

func (s *TestIspSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestIspSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
