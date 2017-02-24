package owl

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	rt "github.com/Cepave/open-falcon-backend/common/reflect/types"

	"github.com/Cepave/open-falcon-backend/common/utils"

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

// Tests the listing of name tags(with paging)
func (suite *TestNameTagSuite) TestListNameTag(c *C) {
	testCases := []*struct {
		value string
		pageSize int32
		expectedIds []int16
		expectedTotalCount int32
	} {
		{ "", 5, []int16{ 3703, 3704, 3701, 3702 }, 4 },
		{ "", 2, []int16{ 3703, 3704 }, 4 },
		{ "pg-tg-bird", 5, []int16{ 3703, 3704 }, 2 },
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		paging := &model.Paging {
			Size: testCase.pageSize,
		}
		testedResult := ListNameTags(
			testCase.value, paging,
		)

		testedIds := utils.MakeAbstractArray(testedResult).
			MapTo(
				func(elem interface{}) interface{} {
					return elem.(*owlModel.NameTag).Id
				},
				rt.TypeOfInt16,
			).GetArray()

		c.Assert(testedIds, DeepEquals, testCase.expectedIds, comment)
		c.Assert(paging.TotalCount, Equals, testCase.expectedTotalCount, comment)
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
	case "TestNameTagSuite.TestListNameTag":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3701, 'pg-tg-car-1'),
				(3702, 'pg-tg-car-2'),
				(3703, 'pg-tg-bird-3'),
				(3704, 'pg-tg-bird-4')
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
	case "TestNameTagSuite.TestListNameTag":
		inTx(
			`DELETE FROM owl_name_tag WHERE nt_id >= 3701 AND nt_id <= 3704`,
		)
	}
}

func (s *TestNameTagSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestNameTagSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
