package owl

import (
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestGroupTagSuite struct{}

var _ = Suite(&TestGroupTagSuite{})

// Tests the getting of name tag by id
func (suite *TestGroupTagSuite) TestGetGroupTagById(c *C) {
	testCases := []*struct {
		sampleId int32
		hasFound bool
	} {
		{ 33061, true },
		{ -10, false },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		c.Assert(
			GetGroupTagById(testCase.sampleId), ocheck.ViableValue, testCase.hasFound,
			comment,
		)
	}
}

func (s *TestGroupTagSuite) SetUpTest(c *C) {
	inTx := DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestGroupTagSuite.TestGetGroupTagById":
		inTx(
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(33061, 'gt-db-1')
			`,
		)
	}
}
func (s *TestGroupTagSuite) TearDownTest(c *C) {
	inTx := DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestGroupTagSuite.TestGetGroupTagById":
		inTx(
			`DELETE FROM owl_group_tag WHERE gt_id = 33061`,
		)
	}
}

func (s *TestGroupTagSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestGroupTagSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
