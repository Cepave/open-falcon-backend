package nqm

import (
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	testHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	sjson "github.com/bitly/go-simplejson"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	. "gopkg.in/check.v1"
)

type TestNqmLogSuite struct{}

var _ = Suite(&TestNqmLogSuite{})

// Tests the query(calling of JSONRPC) by DSL for ICMP log
func (suite *TestNqmLogSuite) TestGetStatisticsOfIcmpByDsl(c *C) {
	testCases := []*struct {
		sampleIdOfAgentProvince int16
		expectedNumberOfResult  int
	} {
		{15, 1},
		{20, 2},
		{23, 0}, // No data
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		icmpParams := &NqmDsl{
			GroupingColumns:      []string{"ag_pv_id"},
			StartTime:            1328407200,
			EndTime:              1328493600,
			IdsOfAgentProvinces:  []int16{testCase.sampleIdOfAgentProvince},
			IdsOfAgentIsps:       []int16{16},
			IdsOfTargetProvinces: []int16{31},
			IdsOfTargetIsps:      []int16{87},
			ProvinceRelation:     -1,
		}

		testedResult, err := getStatisticsOfIcmpByDsl(icmpParams)

		/**
		 * Asserts the content of data
		 */
		c.Assert(err, IsNil, comment)
		c.Assert(len(testedResult), Equals, testCase.expectedNumberOfResult, comment)
		if testCase.expectedNumberOfResult > 0 {
			c.Assert(testedResult[0].grouping, DeepEquals, []int32{20, 40}, comment)
			c.Assert(testedResult[0].metrics.Avg, Equals, float64(45.78), comment)
			c.Assert(testedResult[0].metrics.Max, Equals, int16(88), comment)
			c.Assert(testedResult[0].metrics.NumberOfAgents, Equals, int32(50), comment)
			c.Assert(testedResult[0].metrics.NumberOfTargets, Equals, int32(37), comment)
		}
		// :~)
	}
}

func (s *TestNqmLogSuite) SetUpSuite(c *C) {
	if !testHttp.HasWebConfigOrSkip(c) {
		return
	}

	testHttp.StartGinWebServer(
		c,
		func(engine *gin.Engine) {
			engine.POST("/nqm/icmp/query/by-dsl", mockIcmpService)
		},
	)

	g.SetConfig(
		&g.GlobalConfig{
			NqmLog: &g.NqmLogConfig{
				ServiceUrl: testHttp.GetWebUrl(),
			},
		},
	)

	initIcmp()
}

func mockIcmpService(c *gin.Context) {
	jsonRequestBody := sjson.New()
	if err := c.BindJSON(jsonRequestBody); err != nil {
		panic(err)
	}

	jsonIcmpStatistics := sjson.New()

	jsonIcmpStatistics.Set("grouping", []int32{20, 40})
	jsonIcmpStatistics.Set("max", 88)
	jsonIcmpStatistics.Set("min", 33)
	jsonIcmpStatistics.Set("avg", 45.78)
	jsonIcmpStatistics.Set("mdev", 35.22)
	jsonIcmpStatistics.Set("med", 55)
	jsonIcmpStatistics.Set("count", 100)
	jsonIcmpStatistics.Set("loss", 0.1)
	jsonIcmpStatistics.Set("number_of_sent_packets", 9000)
	jsonIcmpStatistics.Set("number_of_received_packets", 8933)
	jsonIcmpStatistics.Set("number_of_agents", 50)
	jsonIcmpStatistics.Set("number_of_targets", 37)

	resultList := make([]*sjson.Json, 0)

	sampleCondition, err := jsonRequestBody.Get("ids_of_agent_provinces").
		MustArray()[0].(json.Number).Int64()

	if err != nil {
		panic(err)
	}

	switch sampleCondition {
	case 15:
		resultList = append(resultList, jsonIcmpStatistics)
	case 20:
		resultList = append(resultList, jsonIcmpStatistics)
		resultList = append(resultList, jsonIcmpStatistics)
	}

	c.JSON(http.StatusOK, resultList)
}
