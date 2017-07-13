package owl

import (
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"time"
)

type TestLocationSuite struct{}

var _ = Suite(&TestLocationSuite{})

var testedProvinceService = NewProvinceService(
	cache.DataCacheConfig{
		MaxSize: 10, Duration: time.Minute * 5,
	},
)

// Tests the loading of province by id
func (suite *TestLocationSuite) TestGetProvinceById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	}{
		{1, true},
		{-10, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedProvinceService.GetProvinceById(testCase.sampleId)
		c.Assert(testedResult, ocheck.ViableValue, testCase.hasFound, comment)
		c.Assert( // Asserts the cache status
			testedProvinceService.cache.Get(provinceKeyById(testCase.sampleId)),
			ocheck.ViableValue,
			testCase.hasFound,
		)
	}
}

// Tests the loading of provinces by name
func (suite *TestLocationSuite) TestGetProvincesByName(c *C) {
	testCases := []*struct {
		sampleName    string
		expectedFound int
	}{
		{"广", 2},
		{"無此處", 0},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedProvinceService.GetProvincesByName(testCase.sampleName)
		c.Assert(testedResult, HasLen, testCase.expectedFound, comment)
		c.Assert(
			testedProvinceService.cache.Get(provinceKeyByName(testCase.sampleName)),
			ocheck.ViableValue,
			testCase.expectedFound > 0,
		)
	}
}

var testedCityService = NewCityService(
	cache.DataCacheConfig{
		MaxSize: 10, Duration: time.Minute * 5,
	},
)

// Tests the loading of city by id
func (suite *TestLocationSuite) TestGetCity2ById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	}{
		{1, true},
		{-10, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedCityService.GetCity2ById(testCase.sampleId)
		c.Assert(testedResult, ocheck.ViableValue, testCase.hasFound, comment)
		c.Assert( // Asserts the cache status
			testedCityService.cache.Get(cityKeyById(testCase.sampleId)),
			ocheck.ViableValue,
			testCase.hasFound,
		)
	}
}

// Tests the loading of citys by name
func (suite *TestLocationSuite) TestGetCity2sByName(c *C) {
	testCases := []*struct {
		sampleName    string
		expectedFound int
	}{
		{"长", 3},
		{"無此處", 0},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedCityService.GetCity2sByName(testCase.sampleName)
		c.Assert(testedResult, HasLen, testCase.expectedFound, comment)
		c.Assert(
			testedCityService.cache.Get(cityKeyByName(testCase.sampleName)),
			ocheck.ViableValue,
			testCase.expectedFound > 0,
		)
	}
}

func (s *TestLocationSuite) SetUpSuite(c *C) {
	owlDb.DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestLocationSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, owlDb.DbFacade)
}
