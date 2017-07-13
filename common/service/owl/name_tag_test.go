package owl

import (
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"time"
)

type TestNameTagSuite struct{}

var _ = Suite(&TestNameTagSuite{})

var testedNameTagService = NewNameTagService(
	cache.DataCacheConfig{
		MaxSize: 5, Duration: time.Minute * 5,
	},
)

// Tests the loading of name tag by id
func (suite *TestNameTagSuite) TestGetNameTagById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	}{
		{3021, true},
		{-10, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedNameTagService.GetNameTagById(testCase.sampleId)

		c.Assert(testedResult, ocheck.ViableValue, testCase.hasFound, comment)
		c.Assert(
			testedNameTagService.cache.Get(nameTagKeyById(testCase.sampleId)),
			ocheck.ViableValue,
			testCase.hasFound,
		)
	}
}

func (s *TestNameTagSuite) SetUpTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNameTagSuite.TestGetNameTagById":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3021, 'nt-cache-1')
			`,
		)
	}
}
func (s *TestNameTagSuite) TearDownTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNameTagSuite.TestGetNameTagById":
		inTx(
			`
			DELETE FROM owl_name_tag WHERE nt_id = 3021
			`,
		)
	}
}

func (s *TestNameTagSuite) SetUpSuite(c *C) {
	owlDb.DbFacade = dbTest.InitDbFacade(c)
}
func (s *TestNameTagSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, owlDb.DbFacade)
}
