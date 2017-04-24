package restful

import (
	"net/http"
	"strconv"

	json "github.com/Cepave/open-falcon-backend/common/json"
	"github.com/Cepave/open-falcon-backend/common/testing"
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	testingDb "github.com/Cepave/open-falcon-backend/modules/nqm-mng/testing"
	. "gopkg.in/check.v1"
)

type TestHeartbeatItSuite struct{}

var _ = Suite(&TestHeartbeatItSuite{})

func (s *TestHeartbeatItSuite) TestAgentHeartbeat(c *C) {
	testCases := []struct {
		hosts      []string
		timestamp  string
		updateOnly bool
		expect     int64
	}{
		{ // Add new 2 new hosts
			[]string{"001", "002"}, "2010-06-06T10:20:00+08:00", false, 2,
		},
		{ // Insert again
			[]string{"001", "002"}, "2010-06-06T10:20:00+08:00", false, 0,
		},
		{ // Update hosts and add a new one in updateOnly mode
			[]string{"002", "003"}, "2014-06-06T10:20:30+08:00", true, 1,
		},
		{ // Simulate 1 old and 1 new heatbeat in updateOnly mode
			[]string{"001", "002"}, "2011-06-06T10:20:30+08:00", true, 1,
		},
	}

	for _, testCase := range testCases {
		sampleHosts := make([]*model.AgentHeartbeat, len(testCase.hosts))
		sampleTime := testing.ParseTime(c, testCase.timestamp)
		for idx, hostName := range testCase.hosts {
			sampleNumber := strconv.Itoa(idx)
			sampleHosts[idx] = &model.AgentHeartbeat{
				Hostname:      "nqm-mng-it-tc1-" + hostName,
				UpdateTime:    sampleTime.Unix(),
				IP:            "127.0.0." + sampleNumber,
				AgentVersion:  "0.0." + sampleNumber,
				PluginVersion: "12345abcd" + sampleNumber,
			}
		}

		client := httpClientConfig.NewSlingByBase().Post("api/v1/agent/heartbeat?update_only=" + strconv.FormatBool(testCase.updateOnly)).
			BodyJSON(sampleHosts)
		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResp := slintChecker.GetJsonBody(http.StatusOK)

		c.Logf("[Agent heartbeat] JSON Result: %s", json.MarshalPrettyJSON(jsonResp))
		c.Assert(jsonResp.Get("rows_affected").MustInt64(), Equals, testCase.expect)
	}

}

func (s *TestHeartbeatItSuite) SetUpSuite(c *C) {
	testingDb.InitRdb(c)
}
func (s *TestHeartbeatItSuite) TearDownSuite(c *C) {
	testingDb.ReleaseRdb(c)
}

func (s *TestHeartbeatItSuite) TearDownTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestHeartbeatItSuite.TestAgentHeartbeat":
		inTx(
			`DELETE FROM host WHERE hostname LIKE 'nqm-mng-it-tc1-%'`,
		)
	}
}
