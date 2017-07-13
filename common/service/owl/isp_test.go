package owl

import (
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"time"
)

type TestIspSuite struct{}

var _ = Suite(&TestIspSuite{})

var testedIspService = NewIspService(
	cache.DataCacheConfig{
		MaxSize: 10, Duration: time.Minute * 5,
	},
)

// Tests the loading of isp by id
func (suite *TestIspSuite) TestGetIspById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	}{
		{1, true},
		{-10, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedIspService.GetIspById(testCase.sampleId)
		c.Assert(testedResult, ocheck.ViableValue, testCase.hasFound, comment)
		c.Assert( // Asserts the cache status
			testedIspService.cache.Get(ispKeyById(testCase.sampleId)),
			ocheck.ViableValue,
			testCase.hasFound,
		)
	}
}

// Tests the loading of isps by name
func (suite *TestIspSuite) TestGetIspsByName(c *C) {
	testCases := []*struct {
		sampleName    string
		expectedFound int
	}{
		{"北", 2},
		{"無此 ISP", 0},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedIspService.GetIspsByName(testCase.sampleName)
		c.Assert(testedResult, HasLen, testCase.expectedFound, comment)
		c.Assert(
			testedIspService.cache.Get(ispKeyByName(testCase.sampleName)),
			ocheck.ViableValue,
			testCase.expectedFound > 0,
		)
	}
}

func (s *TestIspSuite) SetUpSuite(c *C) {
	owlDb.DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestIspSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, owlDb.DbFacade)
}
