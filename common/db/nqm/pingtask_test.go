package nqm

import (
	"encoding/json"
	"fmt"

	nqmTestingDb "github.com/Cepave/open-falcon-backend/common/db/nqm/testing"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestPingtaskSuite struct{}

var _ = Suite(&TestPingtaskSuite{})

func (s *TestPingtaskSuite) SetUpTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case
		"TestPingtaskSuite.TestGetPingtaskById",
		"TestPingtaskSuite.TestListPingtasks",
		"TestPingtaskSuite.TestUpdateAndGetPingtask":
		inTx(nqmTestingDb.InsertPingtaskSQL)
	case
		"TestPingtaskSuite.TestAssignPingtaskToAgentForAgent",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForAgent",
		"TestPingtaskSuite.TestAssignPingtaskToAgentForPingtask",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForPingtask":
		inTx(nqmTestingDb.InitNqmAgentAndPingtaskSQL...)
	case
		"TestPingtaskSuite.TestAddAndGetPingtask":
	}
}

func (s *TestPingtaskSuite) TearDownTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case
		"TestPingtaskSuite.TestGetPingtaskById",
		"TestPingtaskSuite.TestListPingtasks",
		"TestPingtaskSuite.TestUpdateAndGetPingtask",
		"TestPingtaskSuite.TestAddAndGetPingtask":
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
		expectedNumOfEnabledAgents int8
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
		expectedNumOfEnabledAgents int8
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
	testCases := []*struct {
		query                      *nqmModel.PingtaskQuery
		paging                     commonModel.Paging
		expectedCountOfCurrentPage int
		expectedCountOfAll         int32
	}{
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			2, 2,
		},
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 1, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 2,
		},
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 10, Position: 10, OrderBy: []*commonModel.OrderByEntity{}},
			0, 2,
		},
		{
			&nqmModel.PingtaskQuery{
				Period: "3",
				Name:   "test2",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Period: "40",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Name: "test1",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Enable: "true",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
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
	var pm1 *nqmModel.PingtaskModify
	if err := json.Unmarshal([]byte(`
		{
		  "period" : 15,
		  "name" : "廣東",
		  "enable" : true,
		  "comment" : "This is for some purpose",
		  "filter" : {
		    "ids_of_isp" : [ 7, 8, 9 ],
		    "ids_of_province" : [ 2, 3 ],
		    "ids_of_city" : [ 11, 12, 13 ]
		  }
		}
	`), &pm1); err != nil {
		c.Error(err)
	}
	testCases := []*struct {
		inputPm *nqmModel.PingtaskModify
	}{
		{pm1},
	}
	for i, v := range testCases {
		actual := AddAndGetPingtask(v.inputPm)
		c.Logf("case [%d]: %+v\n", i+1, actual)
		c.Assert(actual, NotNil)
		c.Assert(actual.Period, Equals, int8(15))
		c.Assert(*actual.Name, Equals, "廣東")
		c.Assert(actual.Enable, Equals, true)
		c.Assert(*actual.Comment, Equals, "This is for some purpose")

		// Tricky, so I only test the lengths
		c.Assert(len(actual.Filter.IspFilters), Equals, 3)
		c.Assert(len(actual.Filter.ProvinceFilters), Equals, 2)
		c.Assert(len(actual.Filter.CityFilters), Equals, 3)
		c.Assert(len(actual.Filter.NameTagFilters), Equals, 0)
		c.Assert(len(actual.Filter.GroupTagFilters), Equals, 0)
	}
}

func (suite *TestPingtaskSuite) TestUpdateAndGetPingtask(c *C) {
	var pm1 *nqmModel.PingtaskModify
	if err := json.Unmarshal([]byte(`
		{
		  "period" : 15,
		  "name" : "廣東",
		  "enable" : true,
		  "comment" : "This is for some purpose",
		  "filter" : {
		    "ids_of_isp" : [ 7, 8, 9 ],
		    "ids_of_province" : [ 2, 3 ],
		    "ids_of_city" : [ 11, 12, 13 ]
		  }
		}
	`), &pm1); err != nil {
		c.Error(err)
	}
	var pm2 *nqmModel.PingtaskModify
	if err := json.Unmarshal([]byte(`
		{
		  "period" : 15,
		  "name" : "廣東",
		  "enable" : true,
		  "comment" : "This is for some purpose",
		  "filter" : {
		    "ids_of_isp" : [ 17, 18 ],
		    "ids_of_province" : [ 2, 3, 4 ],
		    "ids_of_city" : [ 3 ]
		  }
		}
	`), &pm2); err != nil {
		c.Error(err)
	}
	testCases := []*struct {
		inputPm             *nqmModel.PingtaskModify
		expectedIspLen      int
		expectedProvinceLen int
		expectedCityLen     int
		expectedNameTagLen  int
		expectedGroupTagLen int
	}{
		{pm1, 3, 2, 3, 0, 0},
		{pm2, 2, 3, 1, 0, 0},
	}
	for i, v := range testCases {
		actual := UpdateAndGetPingtask(10120, v.inputPm)
		c.Logf("case [%d]: %+v\n", i+1, actual)
		c.Assert(actual, NotNil)
		c.Assert(actual.Period, Equals, int8(15))
		c.Assert(*actual.Name, Equals, "廣東")
		c.Assert(actual.Enable, Equals, true)
		c.Assert(*actual.Comment, Equals, "This is for some purpose")

		// Tricky, so I only test the lengths
		c.Assert(len(actual.Filter.IspFilters), Equals, v.expectedIspLen)
		c.Assert(len(actual.Filter.ProvinceFilters), Equals, v.expectedProvinceLen)
		c.Assert(len(actual.Filter.CityFilters), Equals, v.expectedCityLen)
		c.Assert(len(actual.Filter.NameTagFilters), Equals, v.expectedNameTagLen)
		c.Assert(len(actual.Filter.GroupTagFilters), Equals, v.expectedGroupTagLen)
	}
}
