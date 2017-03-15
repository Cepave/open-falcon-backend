package restful

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	commonNqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	. "gopkg.in/check.v1"
)

type TestTargetSuite struct{}

var _ = Suite(&TestTargetSuite{})

// Tests the building of query for targets
func (suite *TestTargetSuite) TestBuildQueryForList(c *C) {
	testCases := []*struct {
		params        string
		expectedQuery *commonNqmModel.TargetQuery
	}{
		{"", &commonNqmModel.TargetQuery{HasIspId: false, HasStatusCondition: false}}, // Nothing
		{ // With all of the supported parameters
			"name=tg-name-1&host=host-1&&isp_id=57&status=1",
			&commonNqmModel.TargetQuery{
				Name:               "tg-name-1",
				Host:               "host-1",
				HasIspId:           true,
				IspId:              57,
				HasStatusCondition: true,
				Status:             true,
			},
		},
	}

	sampleContext := &gin.Context{}

	for i, testCase := range testCases {
		sampleContext.Request, _ = http.NewRequest("GET", "/a?"+testCase.params, nil)

		c.Assert(buildQueryForListTargets(sampleContext), DeepEquals, testCase.expectedQuery, Commentf("Test Case: %d", i+1))
	}
}
