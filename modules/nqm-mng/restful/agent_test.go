package restful

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	commonNqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"

	. "gopkg.in/check.v1"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

// Tests building of query parameters
func (suite *TestAgentSuite) TestBuildQueryForListAgents(c *C) {
	testCases := []*struct {
		params        string
		expectedQuery *commonNqmModel.AgentQuery
	}{
		{"", &commonNqmModel.AgentQuery{HasIspId: false, HasStatusCondition: false}}, // Nothing
		{ // With all of the supported parameters
			"name=name-1&connection_id=gtk-01&hostname=host-1&ip_address=34.55&isp_id=406&status=1",
			&commonNqmModel.AgentQuery{
				Name:               "name-1",
				Hostname:           "host-1",
				ConnectionId:       "gtk-01",
				IpAddress:          "34.55",
				HasIspId:           true,
				IspId:              406,
				HasStatusCondition: true,
				Status:             true,
			},
		},
	}

	sampleContext := &gin.Context{}

	for i, testCase := range testCases {
		sampleContext.Request, _ = http.NewRequest("GET", "/a?"+testCase.params, nil)

		c.Assert(buildQueryForListAgents(sampleContext), DeepEquals, testCase.expectedQuery, Commentf("Test Case: %d", i+1))
	}
}
