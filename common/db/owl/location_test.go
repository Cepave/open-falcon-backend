package owl

import (
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"

	. "gopkg.in/check.v1"
)

type TestLocationSuite struct{}

var _ = Suite(&TestLocationSuite{})

// Tests the check for city over hierarchy of administrative region
func (suite *TestLocationSuite) TestCheckHierarchyForCity(c *C) {
	testCases := []struct {
		provinceId int16
		cityId int16
		hasError bool
	} {
		{ 17, 27, false },
		{ -1, 33, false },
		{ 17, 20, true },
	}

	for _, testCase := range testCases {
		err := CheckHierarchyForCity(testCase.provinceId, testCase.cityId)

		if testCase.hasError {
			c.Logf("Error: %v", err)
			c.Assert(err, NotNil)
		} else {
			c.Assert(err, IsNil)
		}
	}
}

func (s *TestLocationSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}
func (s *TestLocationSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
