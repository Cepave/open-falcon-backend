package http

import (
	"net/http"
	"net/http/httptest"
	"time"

	"encoding/json"
	"github.com/bitly/go-simplejson"

	dsl "github.com/Cepave/open-falcon-backend/modules/query/dsl/nqm_parser"

	. "gopkg.in/check.v1"
)

type TestNqmSuite struct{}

var _ = Suite(&TestNqmSuite{})

type processDslTestCase struct {
	sampleDsl string
	assertion func(*dsl.QueryParams, error)
}

// Tests:
// 1. default value of time
// 2. default value of start time(end time provided)
func (suite *TestNqmSuite) TestProcessDsl(c *C) {
	testCases := []processDslTestCase{
		/**
		 * Normal situation
		 */
		{
			sampleDsl: "starttime=2010-05-01 endtime=2010-05-02",
			assertion: func(testedDsl *dsl.QueryParams, testedError error) {
				c.Assert(testedError, IsNil)
				c.Assert(testedDsl.StartTime.Year(), Equals, 2010)
				c.Assert(testedDsl.StartTime.Day(), Equals, 1)
				c.Assert(testedDsl.EndTime.Year(), Equals, 2010)
				c.Assert(testedDsl.EndTime.Day(), Equals, 2)
			},
		},
		// :~)
		/**
		 * DSL is empty
		 */
		{
			sampleDsl: "     ",
			assertion: func(testedDsl *dsl.QueryParams, testedError error) {
				now := time.Now()
				beforeNow := now.Add(before7Days)
				expectedEndTime := now.Add(24 * time.Hour)

				c.Assert(testedError, IsNil)
				c.Assert(testedDsl.StartTime.Year(), Equals, beforeNow.Year())
				c.Assert(testedDsl.StartTime.Day(), Equals, beforeNow.Day())
				c.Assert(testedDsl.EndTime.Year(), Equals, expectedEndTime.Year())
				c.Assert(testedDsl.EndTime.Day(), Equals, expectedEndTime.Day())
			},
		},
		// :~)
		/**
		 * DSL only has start time
		 */
		{
			sampleDsl: "starttime=2013-11-12",
			assertion: func(testedDsl *dsl.QueryParams, testedError error) {
				c.Assert(testedError, IsNil)
				c.Assert(testedDsl.StartTime.Year(), Equals, 2013)
				c.Assert(testedDsl.StartTime.Day(), Equals, 12)
				c.Assert(testedDsl.EndTime.Year(), Equals, 2013)
				c.Assert(testedDsl.EndTime.Day(), Equals, 12+defaultDaysForTimeRange)
			},
		},
		// :~)
		/**
		 * DSL only has end time
		 */
		{
			sampleDsl: "endtime=2012-07-15",
			assertion: func(testedDsl *dsl.QueryParams, testedError error) {
				c.Assert(testedError, IsNil)
				c.Assert(testedDsl.StartTime.Year(), Equals, 2012)
				c.Assert(testedDsl.StartTime.Day(), Equals, 15-defaultDaysForTimeRange)
				c.Assert(testedDsl.EndTime.Year(), Equals, 2012)
				c.Assert(testedDsl.EndTime.Day(), Equals, 15)
			},
		},
		// :~)
		/**
		 * DSL only has same value of start/end time
		 */
		{
			sampleDsl: "starttime=2012-03-03 endtime=2012-03-03",
			assertion: func(testedDsl *dsl.QueryParams, testedError error) {
				c.Assert(testedError, IsNil)
				c.Assert(testedDsl.StartTime.Year(), Equals, 2012)
				c.Assert(testedDsl.StartTime.Day(), Equals, 3)
				c.Assert(testedDsl.EndTime.Year(), Equals, 2012)
				c.Assert(testedDsl.EndTime.Day(), Equals, 4)
			},
		},
		// :~)
	}

	for _, testCase := range testCases {
		testCase.assertion(processDsl(testCase.sampleDsl))
	}
}

// Tests the error message rendered as JSON
func (suite *TestNqmSuite) TestErrorMessage(c *C) {
	engine := getGinRouter()

	/**
	 * Sets-up HTTP request and response
	 */
	sampleRequest, err := http.NewRequest(http.MethodGet, "/nqm/icmp/list/by-provinces?dsl=v1%3D10", nil)
	c.Assert(err, IsNil)
	respRecorder := httptest.NewRecorder()
	// :~)

	engine.ServeHTTP(respRecorder, sampleRequest)
	c.Logf("Response: %v", respRecorder)

	/**
	 * Asserts the status code of HTTP
	 */
	c.Assert(respRecorder.Code, Equals, 400)
	// :~)

	testedJsonBody := jsonDslError{}
	json.Unmarshal(respRecorder.Body.Bytes(), &testedJsonBody)

	/**
	 * Asserts the JSON body for error message
	 */
	c.Assert(testedJsonBody.Code, Equals, 1)
	c.Assert(testedJsonBody.Message, Matches, ".+Unknown parameter.+")
	// :~)
}

type sampleJsonData struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

// Tests the JSON for result with DSL
func (suite *TestNqmSuite) TestJsonForResultWithDsl(c *C) {
	sampleStartTime, sampleEndTime := time.Now(), time.Now().Add(3*time.Hour)

	testedResultWithDsl := &resultWithDsl{
		queryParams: &dsl.QueryParams{
			StartTime: sampleStartTime,
			EndTime:   sampleEndTime,
		},
		resultData: []sampleJsonData{
			{20, "Bob"},
			{30, "Joe"},
		},
	}

	jsonResult, err := testedResultWithDsl.MarshalJSON()

	c.Assert(err, IsNil)
	c.Logf("JSON: %v", string(jsonResult))

	jsonObject, _ := simplejson.NewJson(jsonResult)

	c.Assert(jsonObject.GetPath("dsl", "start_time").MustInt64(), Equals, sampleStartTime.Unix())
	c.Assert(jsonObject.GetPath("dsl", "end_time").MustInt64(), Equals, sampleEndTime.Unix())
	c.Assert(jsonObject.Get("result").MustArray(), HasLen, 2)
}
