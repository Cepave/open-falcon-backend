package rdb

import (
	"strconv"
	"time"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	"github.com/Cepave/open-falcon-backend/common/testing"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb/test"
	ch "gopkg.in/check.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

type TestHeartbeatSuite struct{}

var _ = ch.Suite(&TestHeartbeatSuite{})

func (suite *TestHeartbeatSuite) TestFalconAgentHeartbeat(c *ch.C) {
	testCases := []struct {
		hosts      []string
		timestamp  string
		updateOnly bool
		expect     int64
	}{
		{ // Add new 3 new hosts
			[]string{"001", "002", "003"}, "2014-05-05T10:20:00+08:00", false, 3,
		},
		{ // Update existing hosts
			[]string{"001", "002", "003"}, "2014-05-05T11:20:30+08:00", false, 3,
		},
		{ // Simulate old heartbeat and a new one
			[]string{"001", "002", "003", "004"}, "2014-04-04T10:20:30+08:00", false, 1,
		},
		{ // Update hosts in updateOnly mode
			[]string{"001", "002", "003", "005"}, "2014-06-06T10:20:30+08:00", true, 3,
		},
		{ // Simulate old heatbeat in updateOnly mode
			[]string{"001", "002", "003", "004"}, "2014-03-03T10:20:30+08:00", true, 0,
		},
	}

	countStmt := DbFacade.SqlxDbCtrl.PreparexExt(`
		SELECT COUNT(*)
		FROM host
		WHERE update_at = FROM_UNIXTIME(?)
			AND hostname LIKE 'nqm-mng-tc1-%'
			AND ip = ?
			AND agent_version = ?
			AND plugin_version = ?
	`)

	for idx, testCase := range testCases {
		comment := ocheck.TestCaseComment(idx)
		ocheck.LogTestCase(c, testCase)

		sampleTime := testing.ParseTime(c, testCase.timestamp)
		sampleNumber := strconv.Itoa(idx + 1)
		sampleIP, sampleAgentVersion, samplePluginVersion :=
			"127.0.0."+sampleNumber, "0.0."+sampleNumber, "12345abcd"+sampleNumber

		sampleHosts := make([]*model.FalconAgentHeartbeat, len(testCase.hosts))
		for idx, hostName := range testCase.hosts {
			sampleHosts[idx] = &model.FalconAgentHeartbeat{
				Hostname:      "nqm-mng-tc1-" + hostName,
				UpdateTime:    sampleTime.Unix(),
				IP:            sampleIP,
				AgentVersion:  sampleAgentVersion,
				PluginVersion: samplePluginVersion,
			}
		}

		result := FalconAgentHeartbeat(sampleHosts, testCase.updateOnly)

		var dbResult int64
		countStmt.QueryRowxAndScan([]interface{}{sampleTime.Unix(), sampleIP, sampleAgentVersion, samplePluginVersion}, &dbResult)
		c.Assert(result.RowsAffected, ch.Equals, testCase.expect, comment)
		c.Assert(dbResult, ch.Equals, testCase.expect, comment)
	}
}

func (suite *TestHeartbeatSuite) TearDownTest(c *ch.C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestHeartbeatSuite.TestFalconAgentHeartbeat":
		inTx(
			`DELETE FROM host WHERE hostname LIKE 'nqm-mng-tc1-%'`,
		)
	}
}

func (suite *TestHeartbeatSuite) SetUpSuite(c *ch.C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (suite *TestHeartbeatSuite) TearDownSuite(c *ch.C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
	DbFacade = nil
}

func inTx(sql ...string) {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

var _ = Describe("Test UpdateNqmAgentHeartbeat()", ginkgoDb.NeedDb(func() {
	BeforeEach(func() {
		inTx(test.InitNqmAgent...)
	})

	AfterEach(func() {
		inTx(test.ClearNqmAgent...)
	})

	now := time.Now()
	yesterday := time.Now().Add(-24 * time.Hour)
	DescribeTable("for newly inserted agents", func(input time.Time) {
		reqs := []*model.NqmAgentHeartbeatRequest{
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-1@201.3.116.1",
				Hostname:     "ct-255-1",
				IpAddress:    "201.3.116.1",
				Timestamp:    ojson.JsonTime(input),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-2@201.3.116.2",
				Hostname:     "ct-255-2",
				IpAddress:    "201.3.116.2",
				Timestamp:    ojson.JsonTime(input),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-3@201.4.23.3",
				Hostname:     "ct-255-3",
				IpAddress:    "201.4.23.3",
				Timestamp:    ojson.JsonTime(input),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-63-1@201.77.23.3",
				Hostname:     "ct-63-1",
				IpAddress:    "201.77.23.3",
				Timestamp:    ojson.JsonTime(input),
			},
		}

		UpdateNqmAgentHeartbeat(reqs)
		for _, req := range reqs {
			agent := SelectNqmAgentByConnId(req.ConnectionId)
			Expect(agent.LastHeartBeat.Unix()).To(Equal(input.Unix()))
		}
	},
		Entry("case: Now", now),
		Entry("case yesterday", yesterday),
	)

	existentTime := time.Now().Add(-240 * time.Hour)
	earlierTime := time.Now().Add(-480 * time.Hour)
	DescribeTable("for existent agents", func(input time.Time, expected time.Time) {
		init := []*model.NqmAgentHeartbeatRequest{
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-1@201.3.116.1",
				Hostname:     "ct-255-1",
				IpAddress:    "201.3.116.1",
				Timestamp:    ojson.JsonTime(existentTime),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-2@201.3.116.2",
				Hostname:     "ct-255-2",
				IpAddress:    "201.3.116.2",
				Timestamp:    ojson.JsonTime(existentTime),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-3@201.4.23.3",
				Hostname:     "ct-255-3",
				IpAddress:    "201.4.23.3",
				Timestamp:    ojson.JsonTime(existentTime),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-63-1@201.77.23.3",
				Hostname:     "ct-63-1",
				IpAddress:    "201.77.23.3",
				Timestamp:    ojson.JsonTime(existentTime),
			},
		}
		UpdateNqmAgentHeartbeat(init)

		reqs := []*model.NqmAgentHeartbeatRequest{
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-1@201.3.116.1",
				Hostname:     "ct-255-1",
				IpAddress:    "201.3.116.1",
				Timestamp:    ojson.JsonTime(input),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-2@201.3.116.2",
				Hostname:     "ct-255-2",
				IpAddress:    "201.3.116.2",
				Timestamp:    ojson.JsonTime(input),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-255-3@201.4.23.3",
				Hostname:     "ct-255-3",
				IpAddress:    "201.4.23.3",
				Timestamp:    ojson.JsonTime(input),
			},
			&model.NqmAgentHeartbeatRequest{
				ConnectionId: "ct-63-1@201.77.23.3",
				Hostname:     "ct-63-1",
				IpAddress:    "201.77.23.3",
				Timestamp:    ojson.JsonTime(input),
			},
		}
		UpdateNqmAgentHeartbeat(reqs)

		for _, req := range reqs {
			agent := SelectNqmAgentByConnId(req.ConnectionId)
			Expect(agent.LastHeartBeat.Unix()).To(Equal(expected.Unix()))
		}
	},
		Entry("case: now", now, now),
		Entry("case: yesterday", yesterday, yesterday),
		Entry("case: earlier than existent value", earlierTime, existentTime),
	)
}))

var _ = Describe("Test SelectNqmAgentByConnId()", ginkgoDb.NeedDb(func() {
	BeforeEach(func() {
		inTx(test.InitNqmAgent...)
	})

	AfterEach(func() {
		inTx(test.ClearNqmAgent...)
	})

	DescribeTable("for existent agents", func(input string, expected types.GomegaMatcher) {
		r := SelectNqmAgentByConnId(input)
		Ω(r).Should(expected)
	},
		Entry("case:     existent", "ct-255-1@201.3.116.1", Not(BeNil())),
		Entry("case: not existent", "ct-255-1@201.3.116.", BeNil()),
		Entry("case:     existent", "ct-255-2@201.3.116.2", Not(BeNil())),
		Entry("case: not existent", "ct-255-2@201.3.116.", BeNil()),
		Entry("case:     existent", "ct-255-3@201.4.23.3", Not(BeNil())),
		Entry("case: not existent", "ct-255-3@201.4.23", BeNil()),
		Entry("case:     existent", "ct-63-1@201.77.23.3", Not(BeNil())),
		Entry("case: not existent", "ct-63-1@201.77.23", BeNil()),
	)
}))
