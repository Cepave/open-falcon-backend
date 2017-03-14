package restful

import (
	"net/http"
	"strconv"

	testingOwlDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	rdb "github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	testingDb "github.com/Cepave/open-falcon-backend/modules/nqm-mng/testing"
	"github.com/dghubble/sling"
	"github.com/mikelue/cepave-owl/common/json"

	. "gopkg.in/check.v1"
)

type TestPingtaskItSuite struct{}

var _ = Suite(&TestPingtaskItSuite{})

func (s *TestPingtaskItSuite) SetUpSuite(c *C) {
	testingDb.InitRdb(c)
}
func (s *TestPingtaskItSuite) TearDownSuite(c *C) {
	testingDb.ReleaseRdb(c)
}

func (s *TestPingtaskItSuite) SetUpTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case
		"TestPingtaskItSuite.TestGetPingtaskById",
		"TestPingtaskItSuite.TestListPingtasks",
		"TestPingtaskItSuite.TestModifyPingtask":
		inTx(testingOwlDb.InsertPingtaskSQL)
	case
		"TestPingtaskItSuite.TestAddPingtaskToAgentForAgent",
		"TestPingtaskItSuite.TestRemovePingtaskToAgentForAgent",
		"TestPingtaskItSuite.TestAddPingtaskToAgentForPingtask",
		"TestPingtaskItSuite.TestRemovePingtaskToAgentForPingtask":
		inTx(testingOwlDb.InitNqmAgentAndPingtaskSQL...)
	}
}
func (s *TestPingtaskItSuite) TearDownTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case
		"TestPingtaskItSuite.TestGetPingtaskById",
		"TestPingtaskItSuite.TestListPingtasks",
		"TestPingtaskItSuite.TestModifyPingtask",
		"TestPingtaskItSuite.TestAddNewPingtask":
		inTx(testingOwlDb.DeletePingtaskSQL)
	case
		"TestPingtaskItSuite.TestAddPingtaskToAgentForAgent",
		"TestPingtaskItSuite.TestRemovePingtaskToAgentForAgent",
		"TestPingtaskItSuite.TestAddPingtaskToAgentForPingtask",
		"TestPingtaskItSuite.TestRemovePingtaskToAgentForPingtask":
		inTx(testingOwlDb.CleanNqmAgentAndPingtaskSQL...)
	}
}

func (suite *TestPingtaskItSuite) TestGetPingtaskById(c *C) {
	testCases := []*struct {
		inputID           int
		expectedStatus    int
		expectedErrorCode int
	}{
		{10119, http.StatusOK, -1},
		{10121, http.StatusNotFound, -1},
	}
	for i, v := range testCases {
		c.Logf("case[%d]:", i)
		client := sling.New().Get(httpClientConfig.String()).
			Path("/api/v1/nqm/pingtask/" + strconv.Itoa(v.inputID))

		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResult := slintChecker.GetJsonBody(v.expectedStatus)

		c.Logf("[Get A Pingtask By ID] JSON Result:\n%s", json.MarshalPrettyJSON(jsonResult))
		switch v.expectedStatus {
		case http.StatusOK:
			c.Assert(jsonResult.Get("id").MustInt(), Equals, v.inputID)
		case http.StatusNotFound:
			c.Assert(jsonResult.Get("error_code").MustInt(), Equals, v.expectedErrorCode)
		}
	}
}

func (suite *TestPingtaskItSuite) TestAddNewPingtask(c *C) {
	testCases := []*struct {
		expectedStatus int
	}{
		{http.StatusCreated},
		{http.StatusCreated},
	}
	for i, v := range testCases {
		c.Logf("case[%d]:", i)
		client := sling.New().Post(httpClientConfig.String()).
			Path("/api/v1/nqm/pingtask").BodyJSON(json.UnmarshalToJson([]byte(`
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
			`)))

		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResult := slintChecker.GetJsonBody(v.expectedStatus)
		c.Logf("[Add Pingtask] JSON Result:\n%s", json.MarshalPrettyJSON(jsonResult))
		c.Assert(jsonResult.Get("comment").MustString(), Equals, "This is for some purpose")
	}
}

func (suite *TestPingtaskItSuite) TestListPingtasks(c *C) {
	client := sling.New().Get(httpClientConfig.String()).
		Path("/api/v1/nqm/pingtasks")

	slintChecker := testingHttp.NewCheckSlint(c, client)

	slintChecker.AssertHasPaging()
	message := slintChecker.GetJsonBody(http.StatusOK)

	c.Logf("[List Pingtasks] JSON Result:\n%s", json.MarshalPrettyJSON(message))
	c.Assert(len(message.MustArray()), Equals, 2)
}

func (suite *TestPingtaskItSuite) TestModifyPingtask(c *C) {
	testCases := []*struct {
		inputID           int
		expectedStatus    int
		expectedErrorCode int
	}{
		{10120, http.StatusOK, 1},
		{10121, http.StatusInternalServerError, -1},
	}
	for i, v := range testCases {
		c.Logf("case[%d]:", i)
		client := sling.New().Put(httpClientConfig.String()).
			Path("/api/v1/nqm/pingtask/" + strconv.Itoa(v.inputID)).BodyJSON(json.UnmarshalToJson([]byte(`
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
			`)))

		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResult := slintChecker.GetJsonBody(v.expectedStatus)
		c.Logf("[Modify Pingtask] JSON Result:\n%s", json.MarshalPrettyJSON(jsonResult))
		switch v.expectedStatus {
		case http.StatusOK:
			c.Assert(len(jsonResult.Get("filter").Get("isps").MustArray()), Equals, 2)
			c.Assert(len(jsonResult.Get("filter").Get("provinces").MustArray()), Equals, 3)
			c.Assert(len(jsonResult.Get("filter").Get("cities").MustArray()), Equals, 1)
			c.Assert(jsonResult.Get("id").MustInt(), Equals, v.inputID)
		case http.StatusInternalServerError:
			c.Assert(jsonResult.Get("error_code").MustInt(), Equals, v.expectedErrorCode)
		}
	}
}

func (suite *TestPingtaskItSuite) TestAddPingtaskToAgentForAgent(c *C) {
	testCases := []*struct {
		inputAID       int
		inputPID       int
		expectedStatus int
	}{
		{24021, 10119, http.StatusCreated},
		{24022, 10119, http.StatusCreated},
		{24021, 10120, http.StatusCreated},
		// i > 2: cases for panic
		{24024, 10121, http.StatusInternalServerError},
		{24025, 10120, http.StatusInternalServerError},
		{24026, 10121, http.StatusInternalServerError},
	}
	for i, v := range testCases {
		c.Logf("case[%d]:", i)
		client := sling.New().Post(httpClientConfig.String()).
			Path("/api/v1/nqm/agent/" + strconv.Itoa((v.inputAID)) + "/pingtask?pingtask_id=" + strconv.Itoa(v.inputPID))
		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResult := slintChecker.GetJsonBody(v.expectedStatus)
		c.Logf("[Modify Pingtask] JSON Result:\n%s", json.MarshalPrettyJSON(jsonResult))
		c.Logf("%+v", jsonResult)
	}
}

func (suite *TestPingtaskItSuite) TestRemovePingtaskToAgentForAgent(c *C) {
	testCases := []*struct {
		inputAID       int
		inputPID       int
		expectedStatus int
	}{
		{24021, 10119, http.StatusOK},
		{24022, 10119, http.StatusOK},
		{24021, 10120, http.StatusOK},
		{24024, 10121, http.StatusOK},
		// i > 3: cases for panic
		{24025, 10120, http.StatusOK},
		{24026, 10121, http.StatusOK},
	}
	for _, v := range testCases {
		req, _ := sling.New().Post(httpClientConfig.String()).
			Path("/api/v1/nqm/agent/" + strconv.Itoa((v.inputAID)) + "/pingtask?pingtask_id=" + strconv.Itoa(v.inputPID)).
			Request()
		client := &http.Client{}
		client.Do(req)
	}
	for i, v := range testCases {
		c.Logf("case[%d]:", i)
		client := sling.New().Delete(httpClientConfig.String()).
			Path("/api/v1/nqm/agent/" + strconv.Itoa((v.inputAID)) + "/pingtask/=" + strconv.Itoa(v.inputPID))
		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResult := slintChecker.GetJsonBody(v.expectedStatus)
		c.Logf("[Modify Pingtask] JSON Result:\n%s", json.MarshalPrettyJSON(jsonResult))
		c.Logf("%+v", jsonResult)
	}
}

func (suite *TestPingtaskItSuite) TestAddPingtaskToAgentForPingtask(c *C) {
	testCases := []*struct {
		inputPID       int
		inputAID       int
		expectedStatus int
	}{
		{10119, 24021, http.StatusCreated},
		{10119, 24022, http.StatusCreated},
		{10120, 24023, http.StatusCreated},
		// i > 2: cases for panic
		{10121, 24024, http.StatusInternalServerError},
		{10120, 24025, http.StatusInternalServerError},
		{10121, 24026, http.StatusInternalServerError},
	}
	for i, v := range testCases {
		c.Logf("case[%d]:", i)
		client := sling.New().Post(httpClientConfig.String()).
			Path("/api/v1/nqm/pingtask/" + strconv.Itoa((v.inputPID)) + "/agent?agent_id=" + strconv.Itoa(v.inputAID))
		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResult := slintChecker.GetJsonBody(v.expectedStatus)
		c.Logf("[Modify Pingtask] JSON Result:\n%s", json.MarshalPrettyJSON(jsonResult))
		c.Logf("%+v", jsonResult)
	}
}

func (suite *TestPingtaskItSuite) TestRemovePingtaskToAgentForPingtask(c *C) {
	testCases := []*struct {
		inputPID       int
		inputAID       int
		expectedStatus int
	}{
		{10119, 24021, http.StatusOK},
		{10119, 24022, http.StatusOK},
		{10120, 24023, http.StatusOK},
		// i > 2: cases for panic
		{10121, 24024, http.StatusOK},
		{10120, 24025, http.StatusOK},
		{10121, 24026, http.StatusOK},
	}
	for _, v := range testCases {
		req, _ := sling.New().Post(httpClientConfig.String()).
			Path("/api/v1/nqm/pingtask/" + strconv.Itoa((v.inputPID)) + "/agent?agent_id=" + strconv.Itoa(v.inputAID)).
			Request()
		client := &http.Client{}
		client.Do(req)
	}
	for i, v := range testCases {
		c.Logf("case[%d]:", i)
		client := sling.New().Delete(httpClientConfig.String()).
			Path("/api/v1/nqm/pingtask/" + strconv.Itoa((v.inputPID)) + "/agent/=" + strconv.Itoa(v.inputAID))
		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResult := slintChecker.GetJsonBody(v.expectedStatus)
		c.Logf("[Modify Pingtask] JSON Result:\n%s", json.MarshalPrettyJSON(jsonResult))
		c.Logf("%+v", jsonResult)
	}
}
