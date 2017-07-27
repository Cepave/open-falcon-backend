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

type TestGroupTagSuite struct{}

var _ = Suite(&TestGroupTagSuite{})

// Tests the getting of name tag by id
func (suite *TestGroupTagSuite) TestGetGroupTagById(c *C) {
	testCases := []*struct {
		sampleId int32
		hasFound bool
	}{
		{33061, true},
		{-10, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		c.Assert(
			GetGroupTagById(testCase.sampleId), ocheck.ViableValue, testCase.hasFound,
			comment,
		)
	}
}

// Tests the listing of group tags
func (suite *TestGroupTagSuite) TestListGroupTags(c *C) {
	testCases := []*struct {
		name               string
		pageSize           int32
		expectedIds        []int32
		expectedTotalCount int32
	}{
		{"", 5, []int32{23041, 23042, 23043, 23044}, 4},
		{"", 2, []int32{23041, 23042}, 4},
		{"ls-gt-gin", 5, []int32{23043, 23044}, 2},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		paging := &model.Paging{
			Size: testCase.pageSize,
		}
		testedResult := ListGroupTags(
			testCase.name, paging,
		)

		testedIds := utils.MakeAbstractArray(testedResult).
			MapTo(
				func(elem interface{}) interface{} {
					return elem.(*owlModel.GroupTag).Id
				},
				rt.TypeOfInt32,
			).GetArray()

		c.Assert(testedIds, DeepEquals, testCase.expectedIds, comment)
		c.Assert(paging.TotalCount, Equals, testCase.expectedTotalCount, comment)
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
	case "TestGroupTagSuite.TestListGroupTags":
		inTx(
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(23041, 'ls-gt-car-1'),
				(23042, 'ls-gt-car-2'),
				(23043, 'ls-gt-gin-3'),
				(23044, 'ls-gt-gin-4')
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
	case "TestGroupTagSuite.TestListGroupTags":
		inTx(
			`
			DELETE FROM owl_group_tag
			WHERE gt_id >= 23041 AND gt_id <= 23044
			`,
		)
	}
}

func (s *TestGroupTagSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}
func (s *TestGroupTagSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
