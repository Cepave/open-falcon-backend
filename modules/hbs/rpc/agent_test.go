package rpc

import (
	"net/rpc"
	"github.com/Cepave/open-falcon-backend/common/model"
	testJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"

	"github.com/Cepave/open-falcon-backend/modules/hbs/db"

	. "gopkg.in/check.v1"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

// Tests the updating of information for basis agent
func (suite *TestAgentSuite) TestReportStatus(c *C) {
	var request = model.AgentReportRequest{
		Hostname: "test-g-01",
		IP: "123.45.61.81",
		AgentVersion: "4.5.31",
		PluginVersion: "1.2.12",
	}
	var response = model.SimpleRpcResponse{}

	testJsonRpc.OpenClient(c, func(client *rpc.Client) {
		err := client.Call("Agent.ReportStatus", request, &response)
		c.Assert(err, IsNil)

		var hostData = struct {
			Ip string `gorm:"column:ip"`
			AgentVersion string `gorm:"column:agent_version"`
			PluginVersion string `gorm:"column:plugin_version"`
		}{}

		/**
		 * Asserts the data in database
		 */
		DbFacade.GormDb.Raw(
			`
			SELECT * FROM host
			WHERE hostname = 'test-g-01'
			`,
		).Scan(&hostData)

		c.Assert(hostData.Ip, Equals, "123.45.61.81")
		c.Assert(hostData.AgentVersion, Equals, "4.5.31")
		c.Assert(hostData.PluginVersion, Equals, "1.2.12")
		// :~)
	})
}

func (s *TestAgentSuite) SetUpSuite(c *C) {
	db.DbInit(dbTest.GetDbConfig(c))
}
func (s *TestAgentSuite) TearDownSuite(c *C) {
	db.Release()
}
func (s *TestAgentSuite) SetUpTest(c *C) {
}
func (s *TestAgentSuite) TearDownTest(c *C) {
	var executeInTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentSuite.TestReportStatus":
		executeInTx(
			`
			DELETE FROM host
			WHERE hostname = 'test-g-01'
			`,
		)
	}
}
