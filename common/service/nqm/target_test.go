package nqm

import (
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"time"
)

type TestTargetSuite struct{}

var _ = Suite(&TestTargetSuite{})

var testedTargetService = NewTargetService(
	cache.DataCacheConfig{
		MaxSize:  10,
		Duration: time.Minute * 5,
	},
)

// Tests the getting of simple target by id
func (suite *TestTargetSuite) TestGetSimpleTarget1ById(c *C) {
	testCases := []*struct {
		sampleId int32
		hasFound bool
	}{
		{44023, true},
		{44024, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		/**
		 * Asserts the found data
		 */
		testedResult := testedTargetService.GetSimpleTarget1ById(testCase.sampleId)
		c.Assert(testedResult, ocheck.ViableValue, testCase.hasFound, comment)
		// :~)

		/**
		 * Asserts the cache
		 */
		testedCache := testedTargetService.cache.Get(getKeyByTargetId(testCase.sampleId))
		c.Assert(testedCache, ocheck.ViableValue, testCase.hasFound, comment)
		// :~)
	}

}

// Tests the loading of SimpleTarget1 by filter
func (suite *TestTargetSuite) TestGetSimpleTarget1sByFilter(c *C) {
	testCases := []*struct {
		sampleFilter   *nqmModel.TargetFilter
		expectedCache  []int32
		expectedNumber int
	}{
		{
			&nqmModel.TargetFilter{
				Host: []string{"no-such-1"},
			},
			[]int32{},
			0,
		},
		{
			&nqmModel.TargetFilter{
				Host: []string{"20.45.71"},
			},
			[]int32{32071, 32072},
			2,
		},
		{
			&nqmModel.TargetFilter{},
			[]int32{32071, 32072, 32073},
			3,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedTargetService.GetSimpleTarget1sByFilter(testCase.sampleFilter)
		c.Assert(testedResult, HasLen, testCase.expectedNumber, comment)

		for _, id := range testCase.expectedCache {
			comment = Commentf("%s. Needs cache by target id: [%d]", comment.CheckCommentString(), id)

			testedCache := testedTargetService.cache.Get(
				getKeyByTargetId(id),
			)

			c.Assert(testedCache, ocheck.ViableValue, true, comment)
		}

		testedTargetService.cache.Clear()
	}
}

func (s *TestTargetSuite) SetUpTest(c *C) {
	var inTx = nqmDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestTargetSuite.TestGetSimpleTarget1sByFilter":
		inTx(
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host)
			VALUES (32071, 'ftg-1-C01', '20.45.71.91'),
				(32072, 'ftg-1-C02', '20.45.71.92'),
				(32073, 'ftg-2-C01', '120.33.27.23')
			`,
		)
	case "TestTargetSuite.TestGetSimpleTarget1ById":
		inTx(
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host)
			VALUES (44023, 'gc-t-name-1', '11.48.76.51')
			`,
		)
	}
}
func (s *TestTargetSuite) TearDownTest(c *C) {
	var inTx = nqmDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestTargetSuite.TestGetSimpleTarget1sByFilter":
		inTx(
			`DELETE FROM nqm_target WHERE tg_id >= 32071 AND tg_id <= 32073`,
		)
	case "TestTargetSuite.TestGetSimpleTarget1ById":
		inTx(
			"DELETE FROM nqm_target WHERE tg_id = 44023",
		)
	}
}

func (s *TestTargetSuite) SetUpSuite(c *C) {
	nqmDb.DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestTargetSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, nqmDb.DbFacade)
}
