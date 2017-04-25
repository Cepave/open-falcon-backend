package rdb

import (
	"strconv"

	"github.com/Cepave/open-falcon-backend/common/testing"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	. "gopkg.in/check.v1"
)

type TestUpdateOrInsertSuite struct{}

var _ = Suite(&TestUpdateOrInsertSuite{})

func (suite *TestUpdateOrInsertSuite) TestAgentHeartbeat(c *C) {
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

		sampleHosts := make([]*model.AgentHeartbeat, len(testCase.hosts))
		for idx, hostName := range testCase.hosts {
			sampleHosts[idx] = &model.AgentHeartbeat{
				Hostname:      "nqm-mng-tc1-" + hostName,
				UpdateTime:    sampleTime.Unix(),
				IP:            sampleIP,
				AgentVersion:  sampleAgentVersion,
				PluginVersion: samplePluginVersion,
			}
		}

		result := AgentHeartbeat(sampleHosts, testCase.updateOnly)

		var dbResult int64
		countStmt.QueryRowxAndScan([]interface{}{sampleTime.Unix(), sampleIP, sampleAgentVersion, samplePluginVersion}, &dbResult)
		c.Assert(result.RowsAffected, Equals, testCase.expect, comment)
		c.Assert(dbResult, Equals, testCase.expect, comment)
	}
}

func (suite *TestUpdateOrInsertSuite) TearDownTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestUpdateOrInsertSuite.TestAgentHeartbeat":
		inTx(
			`DELETE FROM host WHERE hostname LIKE 'nqm-mng-tc1-%'`,
		)
	}
}

func (suite *TestUpdateOrInsertSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (suite *TestUpdateOrInsertSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
	DbFacade = nil
}
