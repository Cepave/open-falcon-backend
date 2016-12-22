package basis

import (
	"database/sql"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

type testCaseOfUpdateAgent struct {
	ip string
	agentVersion string
	pluginVersion string
}

// Tests the refresh(insert or update) of agent information
func (suite *TestAgentSuite) TestUpdateAgent(c *C) {
	testCases := []testCaseOfUpdateAgent {
		{ "1.2.3.4", "1.0", "1.0" },
		{ "1.9.3.4", "1.1", "1.1" },
	}

	for _, testCase := range testCases {
		agentInfo := &commonModel.AgentUpdateInfo {
			0,
			&commonModel.AgentReportRequest {
				Hostname: "test-host-1",
				IP: testCase.ip,
				AgentVersion: testCase.agentVersion,
				PluginVersion: testCase.pluginVersion,
			},
		}

		c.Assert(UpdateAgent(agentInfo), IsNil)

		assertUpdateAgent(c, &testCase)
	}
}

func assertUpdateAgent(c *C, testCase *testCaseOfUpdateAgent) {
	DbFacade.SqlDbCtrl.QueryForRow(
		commonDb.RowCallbackFunc(func(row *sql.Row) {
			var ip, agentVersion, pluginVersion string

			err := row.Scan(&ip, &agentVersion, &pluginVersion)
			commonDb.DbPanic(err)

			c.Assert(ip, Equals, testCase.ip)
			c.Assert(agentVersion, Equals, testCase.agentVersion)
			c.Assert(pluginVersion, Equals, testCase.pluginVersion)
		}),
		`
		SELECT ip, agent_version, plugin_Version
		FROM host
		WHERE hostname = 'test-host-1'
		`,
	);
}

func (s *TestAgentSuite) SetUpSuite(c *C) {
	db.DbInit(dbTest.GetDbConfig(c))
}

func (s *TestAgentSuite) TearDownSuite(c *C) {
	db.Release()
}

func (s *TestAgentSuite) SetUpTest(c *C) {
	switch c.TestName() {
	case "TestAgentSuite.TestUpdateAgent":
		g.SetConfig(&g.GlobalConfig{
			Hosts: "",
		})
	}
}
func (s *TestAgentSuite) TearDownTest(c *C) {
	switch c.TestName() {
	case "TestAgentSuite.TestUpdateAgent":
		g.SetConfig(nil)
		DbFacade.SqlDbCtrl.Exec("DELETE FROM host WHERE hostname = 'test-host-1'")
	}
}
