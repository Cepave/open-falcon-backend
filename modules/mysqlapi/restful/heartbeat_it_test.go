package restful

import (
	"net/http"
	"strconv"
	"time"

	json "github.com/Cepave/open-falcon-backend/common/json"
	cModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/Cepave/open-falcon-backend/common/testing"
	ogko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/test"
	testingDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	ch "gopkg.in/check.v1"
)

type TestHeartbeatItSuite struct{}

var _ = ch.Suite(&TestHeartbeatItSuite{})

func (s *TestHeartbeatItSuite) TestFalconAgentHeartbeat(c *ch.C) {
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
		sampleHosts := make([]*cModel.FalconAgentHeartbeat, len(testCase.hosts))
		sampleTime := testing.ParseTime(c, testCase.timestamp)
		for idx, hostName := range testCase.hosts {
			sampleNumber := strconv.Itoa(idx)
			sampleHosts[idx] = &cModel.FalconAgentHeartbeat{
				Hostname:      "mysqlapi-it-tc1-" + hostName,
				UpdateTime:    sampleTime.Unix(),
				IP:            "127.0.0." + sampleNumber,
				AgentVersion:  "0.0." + sampleNumber,
				PluginVersion: "12345abcd" + sampleNumber,
			}
		}

		client := httpClientConfig.NewClient().Post("api/v1/agent/heartbeat?update_only=" + strconv.FormatBool(testCase.updateOnly)).
			BodyJSON(sampleHosts)
		slintChecker := testingHttp.NewCheckSlint(c, client)
		jsonResp := slintChecker.GetJsonBody(http.StatusOK)

		c.Logf("[Agent heartbeat] JSON Result: %s", json.MarshalPrettyJSON(jsonResp))
		c.Assert(jsonResp.Get("rows_affected").MustInt64(), ch.Equals, testCase.expect)
	}

}

func (s *TestHeartbeatItSuite) SetUpSuite(c *ch.C) {
	itSkipForGocheck(c)
	testingDb.InitRdb(c)
}
func (s *TestHeartbeatItSuite) TearDownSuite(c *ch.C) {
	testingDb.ReleaseRdb(c)
}

func (s *TestHeartbeatItSuite) TearDownTest(c *ch.C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestHeartbeatItSuite.TestFalconAgentHeartbeat":
		inTx(
			`DELETE FROM host WHERE hostname LIKE 'mysqlapi-it-tc1-%'`,
		)
	}
}

var _ = Describe("Test TestNqmAgentHeartbeat()", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			test.DeleteNqmAgentSQL,
			test.DeleteHostSQL,
			test.ResetAutoIncForNqmAgent,
			test.ResetAutoIncForHost,
			test.SetAutoIncForHost,
			test.SetAutoIncForNqmAgent,
			test.InsertHostSQL,
			test.InsertNqmAgentSQL,
		)
	})

	AfterEach(func() {
		inTx(test.ClearNqmAgent...)
	})

	DescribeTable("update an existent agent or instert a new agent", func(inputConnId string, inputHostname string, inputIPAddr string) {
		inputReq := &nqmModel.HeartbeatRequest{
			ConnectionId: inputConnId,
			Hostname:     inputHostname,
			IpAddress:    json.NewIP(inputIPAddr),
			Timestamp:    json.JsonTime(time.Now()),
		}
		resp := testingHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().
				Post("api/v1/heartbeat/nqm/agent").
				BodyJSON(inputReq),
		)
		jsonBody := resp.GetBodyAsJson()
		GinkgoT().Logf("[NQM Agent Heartbeat Response] JSON Result: %s", json.MarshalPrettyJSON(jsonBody))
		Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
		Expect(jsonBody.Get("connection_id").MustString()).To(Equal(inputReq.ConnectionId))
		Expect(jsonBody.Get("hostname").MustString()).To(Equal(inputReq.Hostname))
		Expect(jsonBody.Get("ip_address").MustString()).To(Equal(inputReq.IpAddress.String()))
		Expect(jsonBody.Get("last_heartbeat_time").Int64()).To(Equal(time.Time(inputReq.Timestamp).Unix()))
	},
		Entry("[update] existent agent", "ct-255-1@201.3.116.1", "ct-255-1", "201.3.116.1"),
		Entry("[update] existent agent with duplicated IP address", "ct-255-1@201.3.116.1", "new-ct-255-1", "201.3.116.1"),
		Entry("[update] existent agent with duplicated hostname", "ct-255-1@201.3.116.1", "ct-255-1", "201.3.116.11"),
		Entry("[update] existent agent with duplicated IP address and hostname", "ct-255-1@201.3.116.1", "new-ct-255-1", "201.3.116.11"),
		Entry("[insert] new agent", "new-ct-255-1@201.3.116.1", "new-ct-255-1", "201.3.116.11"),
		Entry("[insert] new agent with duplicated IP address", "new-ct-255-1@201.3.116.1", "new-ct-255-1", "201.3.116.1"),
		Entry("[insert] new agent with duplicated hostname", "new-ct-255-1@201.3.116.1", "ct-255-1", "201.3.116.11"),
		Entry("[insert] new agent with duplicated IP address and hostname", "new-ct-255-1@201.3.116.1", "ct-255-1", "201.3.116.1"),
	)
}))

var _ = Describe("Test TestNqmAgentHeartbeatTargetList()", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			test.DeleteNqmAgentSQL,
			test.DeleteHostSQL,
			test.ResetAutoIncForNqmAgent,
			test.ResetAutoIncForHost,
			test.SetAutoIncForHost,
			test.SetAutoIncForNqmAgent,
			test.InsertHostSQL,
			test.InsertNqmAgentSQL,
			test.InsertNqmtargetSQL,
			test.InsertNqmPingtaskSQL,
			test.InsertNqmAgentPingtaskSQL,
		)
	})

	AfterEach(func() {
		inTx(
			test.DeleteNqmAgentPingtaskSQL,
			test.DeletetNqmPingtaskSQL,
			test.DeleteNqmtargetSQL,
			test.DeleteNqmAgentSQL,
			test.DeleteHostSQL,
			test.ResetAutoIncForNqmAgent,
			test.ResetAutoIncForHost,
		)
	})

	It("gets an empyt target list of an existent agent[24021]", func() {
		resp := testingHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().
				Get("api/v1/heartbeat/nqm/agent/24021/targets"),
		)
		Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
	})

	It("gets the target list of 3 targets from an existent agent[24022]", func() {
		resp := testingHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().
				Get("api/v1/heartbeat/nqm/agent/24022/targets"),
		)
		Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))

		jsonBody := resp.GetBodyAsJson()
		GinkgoT().Logf("[List Agents] JSON Result: %s", json.MarshalPrettyJSON(jsonBody))
		Expect(jsonBody.MustArray()).To(HaveLen(3))
	})

	It("gets error code 500 for an nonexistent agent[24020]", func() {
		resp := testingHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().
				Get("api/v1/heartbeat/nqm/agent/24020/targets"),
		)
		Expect(resp).To(ogko.MatchHttpStatus(http.StatusInternalServerError))
	})
}))
