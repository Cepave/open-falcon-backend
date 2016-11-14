package owl

import (
	"reflect"

	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestIspSuite struct{}

var _ = Suite(&TestIspSuite{})

func (suite *TestIspSuite) TestGetISPByName(c *C) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"北", []string{"北京三信时代", "北京宽捷"}},
		{"方", []string{"方正宽带"}},
		{"幹", []string{}},
	}

	for _, v := range testCases {
		got := GetISPsByName(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			c.Error("Error:", got, "!=", v.expected)
		} else {
			c.Log(got, "==", v.expected)
		}
	}
}

func (s *TestIspSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestIspSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
