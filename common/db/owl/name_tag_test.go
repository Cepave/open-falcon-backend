package owl

import (
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestNameTagSuite struct{}

var _ = Suite(&TestNameTagSuite{})

// Tests the getting of name tag by id
func (suite *TestNameTagSuite) TestGetNameTagById(c *C) {
	testCases := []*struct {
		sampleId int16
		hasFound bool
	} {
		{ 2901, true },
		{ 2902, false },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		c.Assert(
			GetNameTagById(testCase.sampleId), ocheck.ViableValue, testCase.hasFound,
			comment,
		)
	}
}

func (s *TestNameTagSuite) SetUpTest(c *C) {
	inTx := DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNameTagSuite.TestGetNameTagById":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(2901, 'here-we-1')
			`,
		)
	}
}
func (s *TestNameTagSuite) TearDownTest(c *C) {
	inTx := DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNameTagSuite.TestGetNameTagById":
		inTx(
			`
			DELETE FROM owl_name_tag WHERE nt_id = 2901
			`,
		)
	}
}

func (s *TestNameTagSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestNameTagSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
