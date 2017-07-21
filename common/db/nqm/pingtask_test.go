package nqm

import (
	"fmt"

	nqmTestingDb "github.com/Cepave/open-falcon-backend/common/db/nqm/testing"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"

	. "gopkg.in/check.v1"
)

type TestPingtaskSuite struct{}

var _ = Suite(&TestPingtaskSuite{})

func (s *TestPingtaskSuite) SetUpTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestPingtaskSuite.TestAddAndGetPingtask":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(7801, 'pt-cp-1'),
				(7802, 'pt-cp-2')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(10201, 'pt-ii-1'),
				(10202, 'pt-ii-2')
			`,
		)
	case "TestPingtaskSuite.TestUpdateAndGetPingtask":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(6931, 'pt-cp-1'),
				(6932, 'pt-cp-2')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(28441, 'pt-ii-1'),
				(28442, 'pt-ii-2')
			`,
			`
			INSERT INTO nqm_ping_task(pt_id, pt_name, pt_period)
			VALUES(8702, 'pt-update-test-1', 40)
			`,
		)
	case
		"TestPingtaskSuite.TestGetPingtaskById",
		"TestPingtaskSuite.TestListPingtasks":
		inTx(nqmTestingDb.InsertPingtaskSQL)
	case
		"TestPingtaskSuite.TestAssignPingtaskToAgentForAgent",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForAgent",
		"TestPingtaskSuite.TestAssignPingtaskToAgentForPingtask",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForPingtask":
		inTx(nqmTestingDb.InitNqmAgentAndPingtaskSQL...)
	}
}

func (s *TestPingtaskSuite) TearDownTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestPingtaskSuite.TestAddAndGetPingtask":
		inTx(
			"DELETE FROM nqm_ping_task WHERE pt_name LIKE 'add-pt-%'",
			"DELETE FROM owl_name_tag WHERE nt_id >= 7801 AND nt_id <= 7802",
			"DELETE FROM owl_group_tag WHERE gt_id >= 10201 AND gt_id <= 10202",
		)
	case "TestPingtaskSuite.TestUpdateAndGetPingtask":
		inTx(
			"DELETE FROM nqm_ping_task WHERE pt_id = 8702",
			"DELETE FROM owl_name_tag WHERE nt_id >= 6931 AND nt_id <= 6932",
			"DELETE FROM owl_group_tag WHERE gt_id >= 28441 AND gt_id <= 28442",
		)
	case
		"TestPingtaskSuite.TestGetPingtaskById",
		"TestPingtaskSuite.TestListPingtasks":
		inTx(nqmTestingDb.DeletePingtaskSQL)
	case
		"TestPingtaskSuite.TestAssignPingtaskToAgentForAgent",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForAgent",
		"TestPingtaskSuite.TestAssignPingtaskToAgentForPingtask",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForPingtask":
		inTx(nqmTestingDb.CleanNqmAgentAndPingtaskSQL...)
	}
}

func (s *TestPingtaskSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
	owlDb.DbFacade = DbFacade
}
func (s *TestPingtaskSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
	owlDb.DbFacade = nil
}

func (suite *TestPingtaskSuite) TestAssignPingtaskToAgentForAgent(c *C) {
	testCases := []*struct {
		inputAID                      int32
		inputPID                      int32
		expectedNumOfEnabledPingtasks int32
		expectedAgent                 *nqmModel.Agent
		expectedErr                   error
	}{
		{24021, 10119, 1, GetAgentById(24021), nil},
		{24022, 10119, 1, GetAgentById(24022), nil},
		{24021, 10120, 2, GetAgentById(24021), nil},
		// i > 2: cases for panic
		{24024, 10121, 0, nil, nil},
		{24025, 10120, 1, nil, nil},
		{24026, 10121, -1, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i > 2 {
			c.Assert(func() (*nqmModel.Agent, error) { return AssignPingtaskToAgentForAgent(v.inputAID, v.inputPID) }, PanicMatches, `*.FOREIGN KEY.*`)
			continue
		}
		actualAgent, actualErr := AssignPingtaskToAgentForAgent(v.inputAID, v.inputPID)
		c.Assert(actualAgent, NotNil)
		c.Assert(actualAgent.NumOfEnabledPingtasks, Equals, v.expectedNumOfEnabledPingtasks)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestRemovePingtaskFromAgentForAgent(c *C) {
	testCases := []*struct {
		inputAID                      int32
		inputPID                      int32
		expectedNumOfEnabledPingtasks int32
		expectedAgent                 *nqmModel.Agent
		expectedErr                   error
	}{
		{24021, 10119, 1, GetAgentById(24021), nil},
		{24022, 10119, 0, GetAgentById(24022), nil},
		{24021, 10120, 0, GetAgentById(24021), nil},
		{24024, 10121, 0, GetAgentById(24024), nil},
		// i > 3: Not deleting
		{24025, 10120, -1, nil, nil},
		{24026, 10121, -1, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i > 3 {
			actualAgent, _ := RemovePingtaskFromAgentForAgent(v.inputAID, v.inputPID)
			c.Assert(actualAgent, IsNil)
			continue
		}
		if i == 3 {
			actualAgent, actualErr := RemovePingtaskFromAgentForAgent(v.inputAID, v.inputPID)
			c.Assert(actualAgent, NotNil)
			c.Assert(actualAgent.NumOfEnabledPingtasks, Equals, v.expectedNumOfEnabledPingtasks)
			c.Assert(actualErr, IsNil)
			continue
		}
		AssignPingtaskToAgentForAgent(v.inputAID, v.inputPID)
	}
	for i, v := range testCases {
		if i > 2 {
			break
		}
		actualAgent, actualErr := RemovePingtaskFromAgentForAgent(v.inputAID, v.inputPID)
		c.Assert(actualAgent, NotNil)
		c.Assert(actualAgent.NumOfEnabledPingtasks, Equals, v.expectedNumOfEnabledPingtasks)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestGetPingtaskById(c *C) {
	testCases := []*struct {
		input int32
	}{
		{10119}, // NotNil
		{10120}, // NotNil
		// i > 1: cases for panic
		{10121}, //IsNil
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		actual := GetPingtaskById(v.input)
		if i > 1 {
			c.Assert(actual, IsNil)
			continue
		}
		c.Assert(actual, NotNil)
	}
}

func (suite *TestPingtaskSuite) TestAssignPingtaskToAgentForPingtask(c *C) {
	testCases := []*struct {
		inputAID                   int32
		inputPID                   int32
		expectedNumOfEnabledAgents int32
		expectedPingtask           *nqmModel.PingtaskView
		expectedErr                error
	}{
		{24021, 10119, 1, GetPingtaskById(10119), nil},
		{24022, 10119, 2, GetPingtaskById(10119), nil},
		{24023, 10120, 1, GetPingtaskById(10120), nil},
		// i > 2: cases for panic
		{24024, 10121, -1, nil, nil},
		{24025, 10120, 1, GetPingtaskById(10120), nil},
		{24026, 10121, -1, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i > 2 {
			c.Assert(func() (*nqmModel.PingtaskView, error) {
				return AssignPingtaskToAgentForPingtask(v.inputAID, v.inputPID)
			}, PanicMatches, `*.FOREIGN KEY.*`)
			continue
		}
		actualPingtask, actualErr := AssignPingtaskToAgentForPingtask(v.inputAID, v.inputPID)
		c.Assert(actualPingtask, NotNil)
		c.Assert(actualPingtask.NumOfEnabledAgents, Equals, v.expectedNumOfEnabledAgents)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestRemovePingtaskFromAgentForPingtask(c *C) {
	testCases := []*struct {
		inputAID                   int32
		inputPID                   int32
		expectedNumOfEnabledAgents int32
		expectedPingtask           *nqmModel.PingtaskView
		expectedErr                error
	}{
		{24021, 10119, 1, GetPingtaskById(10119), nil},
		{24022, 10119, 0, GetPingtaskById(10119), nil},
		{24023, 10120, 0, GetPingtaskById(10120), nil},
		// i > 2: Not deleting
		{24024, 10121, -1, nil, nil},
		{24025, 10120, 1, nil, nil},
		{24026, 10121, -1, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i == 3 || i == 5 {
			actualPingtask, _ := RemovePingtaskFromAgentForPingtask(v.inputAID, v.inputPID)
			c.Assert(actualPingtask, IsNil)
			continue
		}
		if i == 4 {
			actualPingtask, actualErr := RemovePingtaskFromAgentForPingtask(v.inputAID, v.inputPID)
			c.Assert(actualPingtask, NotNil)
			c.Assert(actualPingtask.NumOfEnabledAgents, Equals, v.expectedNumOfEnabledAgents)
			c.Assert(actualErr, IsNil)
			continue
		}
		AssignPingtaskToAgentForPingtask(v.inputAID, v.inputPID)
	}
	for i, v := range testCases {
		if i > 2 {
			break
		}
		actualPingtask, actualErr := RemovePingtaskFromAgentForPingtask(v.inputAID, v.inputPID)
		c.Assert(actualPingtask, NotNil)
		c.Assert(actualPingtask.NumOfEnabledAgents, Equals, v.expectedNumOfEnabledAgents)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestListPingtasks(c *C) {
	orderBy := []*commonModel.OrderByEntity{
		{"id", commonModel.Descending},
		{"period", commonModel.Descending},
		{"name", commonModel.Ascending},
		{"enable", commonModel.Ascending},
		{"comment", commonModel.Ascending},
		{"num_of_enabled_agents", commonModel.Ascending},
	}

	testCases := []*struct {
		query                      *nqmModel.PingtaskQuery
		paging                     commonModel.Paging
		expectedCountOfCurrentPage int
		expectedCountOfAll         int32
	}{
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: orderBy},
			2, 2,
		},
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 1, Position: 1, OrderBy: orderBy},
			1, 2,
		},
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 10, Position: 10, OrderBy: orderBy},
			0, 2,
		},
		{
			&nqmModel.PingtaskQuery{
				Period: "3",
				Name:   "test2",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: orderBy},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Period: "40",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: orderBy},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Name: "test1",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: orderBy},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Enable: "true",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: orderBy},
			2, 2,
		},
	}

	for i, v := range testCases {
		fmt.Printf("%+v\n", v)
		actualResult, actualPaging := ListPingtasks(v.query, v.paging)
		c.Logf("case [%d]:", i)
		c.Logf("[List] Query condition: %v. Number of agents: %d", v.query, len(actualResult))

		for _, pingtask := range actualResult {
			c.Logf("[List] Pingtask: %v.", pingtask)
		}
		c.Assert(actualResult, HasLen, v.expectedCountOfCurrentPage)
		c.Assert(actualPaging.TotalCount, Equals, v.expectedCountOfAll)
	}

}

func (suite *TestPingtaskSuite) TestAddAndGetPingtask(c *C) {
	sPtr := func(v string) *string { return &v }
	var newPingTask = &nqmModel.PingtaskModify{
		Period:  30,
		Name:    sPtr("add-pt-廣東"),
		Enable:  true,
		Comment: sPtr("This is for some purpose"),
	}

	testCases := []*struct {
		filter *nqmModel.PingtaskModifyFilter
	}{
		{
			&nqmModel.PingtaskModifyFilter{
				IspIds:      []int16{},
				ProvinceIds: []int16{},
				CityIds:     []int16{},
				NameTagIds:  []int16{},
				GroupTagIds: []int32{},
			},
		},
		{
			&nqmModel.PingtaskModifyFilter{
				IspIds:      []int16{3, 4},
				ProvinceIds: []int16{11, 12, 13},
				CityIds:     []int16{51, 52},
				NameTagIds:  []int16{7801, 7802},
				GroupTagIds: []int32{10201, 10202},
			},
		},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		newPingTask.Filter = testCase.filter

		addedPingTask := AddAndGetPingtask(newPingTask)

		c.Assert(addedPingTask, NotNil, comment)
		c.Assert(addedPingTask.Period, Equals, newPingTask.Period, comment)
		c.Assert(addedPingTask.Name, DeepEquals, newPingTask.Name, comment)
		c.Assert(addedPingTask.Enable, Equals, newPingTask.Enable, comment)
		c.Assert(addedPingTask.Comment, DeepEquals, newPingTask.Comment, comment)

		/**
		 * Asserts filters
		 */
		filters := addedPingTask.Filter
		c.Assert(filters.IspFilters, HasLen, len(testCase.filter.IspIds), comment)
		c.Assert(filters.ProvinceFilters, HasLen, len(testCase.filter.ProvinceIds), comment)
		c.Assert(filters.CityFilters, HasLen, len(testCase.filter.CityIds), comment)
		c.Assert(filters.NameTagFilters, HasLen, len(testCase.filter.NameTagIds), comment)
		c.Assert(filters.GroupTagFilters, HasLen, len(testCase.filter.GroupTagIds), comment)
		// :~)
	}
}

func (suite *TestPingtaskSuite) TestUpdateAndGetPingtask(c *C) {
	sPtr := func(v string) *string { return &v }
	var modifiedPingTask = &nqmModel.PingtaskModify{
		Period:  78,
		Enable:  false,
		Name:    sPtr("up-name-88"),
		Comment: sPtr("up-comment-71"),
	}

	testCases := []*struct {
		filter *nqmModel.PingtaskModifyFilter
	}{
		{
			&nqmModel.PingtaskModifyFilter{
				IspIds:      []int16{3, 4},
				ProvinceIds: []int16{11, 12, 13},
				CityIds:     []int16{51, 52},
				NameTagIds:  []int16{6931, 6932},
				GroupTagIds: []int32{28441, 28442},
			},
		},
		{
			&nqmModel.PingtaskModifyFilter{
				IspIds:      []int16{},
				ProvinceIds: []int16{},
				CityIds:     []int16{},
				NameTagIds:  []int16{},
				GroupTagIds: []int32{},
			},
		},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		modifiedPingTask.Filter = testCase.filter

		updatedPingTask := UpdateAndGetPingtask(8702, modifiedPingTask)

		c.Assert(updatedPingTask.Period, Equals, modifiedPingTask.Period, comment)
		c.Assert(updatedPingTask.Name, DeepEquals, modifiedPingTask.Name, comment)
		c.Assert(updatedPingTask.Enable, Equals, modifiedPingTask.Enable, comment)
		c.Assert(updatedPingTask.Comment, DeepEquals, modifiedPingTask.Comment, comment)

		/**
		 * Asserts filters
		 */
		filters := updatedPingTask.Filter
		c.Assert(filters.IspFilters, HasLen, len(testCase.filter.IspIds), comment)
		c.Assert(filters.ProvinceFilters, HasLen, len(testCase.filter.ProvinceIds), comment)
		c.Assert(filters.CityFilters, HasLen, len(testCase.filter.CityIds), comment)
		c.Assert(filters.NameTagFilters, HasLen, len(testCase.filter.NameTagIds), comment)
		c.Assert(filters.GroupTagFilters, HasLen, len(testCase.filter.GroupTagIds), comment)
		// :~)
	}
}
