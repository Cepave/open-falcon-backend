package owl

import (
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"time"
)

type TestGroupTagSuite struct{}

var _ = Suite(&TestGroupTagSuite{})

var testedGroupTagService = NewGroupTagService(
	cache.DataCacheConfig{
		MaxSize: 5, Duration: time.Minute * 5,
	},
)

// Tests the loading of name tag by id
func (suite *TestGroupTagSuite) TestGetGroupTagById(c *C) {
	testCases := []*struct {
		sampleId int32
		hasFound bool
	}{
		{77031, true},
		{-10, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testedGroupTagService.GetGroupTagById(testCase.sampleId)

		c.Assert(testedResult, ocheck.ViableValue, testCase.hasFound, comment)
		c.Assert(
			testedGroupTagService.cache.Get(groupTagKeyById(testCase.sampleId)),
			ocheck.ViableValue,
			testCase.hasFound,
		)
	}
}

func (s *TestGroupTagSuite) SetUpTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestGroupTagSuite.TestGetGroupTagById":
		inTx(
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(77031, 'gt-cc-1')
			`,
		)
	}
}
func (s *TestGroupTagSuite) TearDownTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestGroupTagSuite.TestGetGroupTagById":
		inTx(
			`
			DELETE FROM owl_group_tag WHERE gt_id = 77031
			`,
		)
	}
}

func (s *TestGroupTagSuite) SetUpSuite(c *C) {
	owlDb.DbFacade = dbTest.InitDbFacade(c)
}
func (s *TestGroupTagSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, owlDb.DbFacade)
}
