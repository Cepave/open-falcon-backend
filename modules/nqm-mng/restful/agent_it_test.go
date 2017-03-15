package restful

import (
	"net/http"

	json "github.com/Cepave/open-falcon-backend/common/json"
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	testingDb "github.com/Cepave/open-falcon-backend/modules/nqm-mng/testing"

	rdb "github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"

	"github.com/dghubble/sling"

	. "gopkg.in/check.v1"
)

type TestAgentItSuite struct{}

var _ = Suite(&TestAgentItSuite{})

// Tests the getting of agent by id
func (suite *TestAgentItSuite) TestGetAgentById(c *C) {
	client := sling.New().Get(httpClientConfig.String()).
		Path("/api/v1/nqm/agent/36771")

	slintChecker := testingHttp.NewCheckSlint(c, client)
	jsonResult := slintChecker.GetJsonBody(http.StatusOK)

	c.Logf("[Get A Agent] JSON Result: %s", json.MarshalPrettyJSON(jsonResult))
	c.Assert(jsonResult.Get("id").MustInt(), Equals, 36771)
}

// Tests the adding of new agent
func (suite *TestAgentItSuite) TestAddNewAgent(c *C) {
	jsonBody := &struct {
		Name         string   `json:"name"`
		Hostname     string   `json:"hostname"`
		ConnectionId string   `json:"connection_id"`
		Status       bool     `json:status`
		Comment      string   `json:"comment"`
		IspId        int      `json:"isp_id"`
		ProvinceId   int      `json:"province_id"`
		CityId       int      `json:"city_id"`
		NameTag      string   `json:"name_tag"`
		GroupTags    []string `json:"group_tags"`
	}{
		Name:         "new-agent-ccc",
		Hostname:     "new-host-cccc",
		Status:       true,
		ConnectionId: "new-agent@blue.12.91.33",
		Comment:      "This is new agent by blue 12.91 ***",
		IspId:        8,
		ProvinceId:   9,
		CityId:       130,
		NameTag:      "rest-nt-1",
		GroupTags:    []string{"pp-rest-tag-1", "pp-rest-tag-2"},
	}

	testCases := []*struct {
		expectedStatus    int
		expectedErrorCode int
	}{
		{http.StatusOK, -1},
		{http.StatusConflict, 1},
	}

	for _, testCase := range testCases {
		client := sling.New().Post(httpClientConfig.String()).
			Path("/api/v1/nqm/agent").
			BodyJSON(jsonBody)

		slintChecker := testingHttp.NewCheckSlint(c, client)

		jsonResp := slintChecker.GetJsonBody(testCase.expectedStatus)

		c.Logf("[Add Agent] JSON Result: %s", json.MarshalPrettyJSON(jsonResp))

		switch testCase.expectedStatus {
		case http.StatusConflict:
			c.Assert(jsonResp.Get("error_code").MustInt(), Equals, testCase.expectedErrorCode)
		}

		if testCase.expectedStatus != http.StatusOK {
			continue
		}

		c.Assert(jsonResp.Get("name").MustString(), Equals, jsonBody.Name)
		c.Assert(jsonResp.Get("connection_id").MustString(), Equals, jsonBody.ConnectionId)
		c.Assert(jsonResp.Get("ip_address").MustString(), Equals, "0.0.0.0")
		c.Assert(jsonResp.Get("hostname").MustString(), Equals, jsonBody.Hostname)
		c.Assert(jsonResp.Get("comment").MustString(), Equals, jsonBody.Comment)
		c.Assert(jsonResp.Get("status").MustBool(), Equals, jsonBody.Status)
		c.Assert(jsonResp.Get("isp").Get("id").MustInt(), Equals, jsonBody.IspId)
		c.Assert(jsonResp.Get("province").Get("id").MustInt(), Equals, jsonBody.ProvinceId)
		c.Assert(jsonResp.Get("city").Get("id").MustInt(), Equals, jsonBody.CityId)
		c.Assert(jsonResp.Get("name_tag").Get("value").MustString(), Equals, jsonBody.NameTag)
		c.Assert(jsonResp.Get("group_tags").MustArray(), HasLen, len(jsonBody.GroupTags))
	}
}

// Tests the listing of agents
func (suite *TestAgentItSuite) TestListAgents(c *C) {
	client := sling.New().Get(httpClientConfig.String()).
		Path("/api/v1/nqm/agents")

	slintChecker := testingHttp.NewCheckSlint(c, client)

	slintChecker.AssertHasPaging()
	message := slintChecker.GetJsonBody(http.StatusOK)

	c.Logf("[List Agents] JSON Result: %s", json.MarshalPrettyJSON(message))
	c.Assert(len(message.MustArray()), Equals, 3)
}

// Tests the modifying of agent
func (suite *TestAgentItSuite) TestModifyAgent(c *C) {
	jsonBody := &struct {
		Name       string   `json:"name"`
		Status     bool     `json:status`
		Comment    string   `json:"comment"`
		IspId      int      `json:"isp_id"`
		ProvinceId int      `json:"province_id"`
		CityId     int      `json:"city_id"`
		NameTag    string   `json:"name_tag"`
		GroupTags  []string `json:"group_tags"`
	}{
		Name:       "Update-Agent-1",
		Status:     false,
		Comment:    "This is updated comment",
		IspId:      3,
		ProvinceId: 11,
		CityId:     230,
		NameTag:    "rest-nt-9",
		GroupTags:  []string{"rest-gt-91", "rest-gt-92", "rest-gt-93"},
	}

	client := sling.New().Put(httpClientConfig.String()).
		Path("/api/v1/nqm/agent/23041").
		BodyJSON(jsonBody)

	slintChecker := testingHttp.NewCheckSlint(c, client)

	jsonResult := slintChecker.GetJsonBody(http.StatusOK)

	c.Logf("Update agent: %v", json.MarshalPrettyJSON(jsonResult))

	c.Assert(jsonResult.Get("name").MustString(), Equals, jsonBody.Name)
	c.Assert(jsonResult.Get("comment").MustString(), Equals, jsonBody.Comment)
	c.Assert(jsonResult.Get("status").MustBool(), Equals, jsonBody.Status)
	c.Assert(jsonResult.Get("isp").Get("id").MustInt(), Equals, jsonBody.IspId)
	c.Assert(jsonResult.Get("province").Get("id").MustInt(), Equals, jsonBody.ProvinceId)
	c.Assert(jsonResult.Get("city").Get("id").MustInt(), Equals, jsonBody.CityId)
	c.Assert(jsonResult.Get("name_tag").Get("value").MustString(), Equals, jsonBody.NameTag)
	c.Assert(jsonResult.Get("group_tags").MustArray(), HasLen, 3)
}

func (s *TestAgentItSuite) SetUpSuite(c *C) {
	testingDb.InitRdb(c)
}
func (s *TestAgentItSuite) TearDownSuite(c *C) {
	testingDb.ReleaseRdb(c)
}

func (s *TestAgentItSuite) SetUpTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentItSuite.TestGetAgentById":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(407, 'nt-rest-01')
			`,
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(25101, 'hs-rest-get-1', '', '')
			`,
			`
			-- IP: 87.90.6.55
			INSERT INTO nqm_agent(
				ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status,
				ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
			)
			VALUES(36771, 25101, 'ag-rest-name-1', 'ag-rest-1@87.90.6.55', 'ag-get-1.rest.com', x'575A0637', 1, 3, 3, 5, 407)
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(50091, 'BOCC-1'),
				(50092, 'BOCC-2'),
				(50093, 'BOCC-3')
			`,
			`
			INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
			VALUES(36771, 50091), (36771, 50092), (36771, 50093)
			`,
		)
	case "TestAgentItSuite.TestListAgents":
		inTx(
			`
			INSERT INTO host(id, hostname, agent_version, plugin_version)
			VALUES(22091, 'agent-it-01', '', '')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_isp_id)
			VALUES(4321, 22091, 'agent-it-01', 'agent-01@28.71.19.22', 'agent-01.fb.com', x'1C471316', 7),
				(4322, 22091, 'agent-it-02', 'agent-02@28.71.19.23', 'agent-02.fb.com', x'1C471317', 7),
				(4323, 22091, 'agent-it-03', 'agent-03@28.71.19.23', 'agent-03.fb.com', x'1C471318', 7)
			`,
		)
	case "TestAgentItSuite.TestModifyAgent":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(8461, 'rest-nt-1')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(20871, "rest-gt-1"), (20872, "rest-gt-2")
			`,
			`
			INSERT INTO host(id, hostname)
			VALUES(4401, '33.99.44.17')
			`,
			`
			INSERT INTO nqm_agent(ag_id, ag_hs_id, ag_connection_id, ag_hostname, ag_ip_address)
			VALUES(23041, 4401, 'update-agent@33.99.44.17', '33.99.44.17', x'21632C11')
			`,
			`
			INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
			VALUES(23041, 20871),(23041, 20872)
			`,
		)
	}
}
func (s *TestAgentItSuite) TearDownTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestAgentItSuite.TestGetAgentById":
		inTx(
			`
			DELETE FROM nqm_agent
			WHERE ag_id = 36771
			`,
			`
			DELETE FROM host
			WHERE id = 25101
			`,
			`
			DELETE FROM owl_name_tag
			WHERE nt_id = 407
			`,
			`
			DELETE FROM owl_group_tag
			WHERE gt_id >= 50091 AND
				gt_id <= 50093
			`,
		)
	case "TestAgentItSuite.TestListAgents":
		inTx(
			"DELETE FROM nqm_agent WHERE ag_id >= 4321 AND ag_id <= 4323",
			"DELETE FROM host WHERE id = 22091",
		)
	case "TestAgentItSuite.TestAddNewAgent":
		inTx(
			"DELETE FROM nqm_agent WHERE ag_connection_id = 'new-agent@blue.12.91.33'",
			"DELETE FROM host WHERE hostname = 'new-agent-blue-1'",
			"DELETE FROM owl_name_tag where nt_value = 'rest-nt-1'",
			"DELETE FROM owl_group_tag where gt_name LIKE 'pp-rest-tag-%'",
		)
	case "TestAgentItSuite.TestModifyAgent":
		inTx(
			"DELETE FROM nqm_agent WHERE ag_id = 23041",
			"DELETE FROM host WHERE id = 4401",
			"DELETE FROM owl_name_tag WHERE nt_value LIKE 'rest-nt-%'",
			"DELETE FROM owl_group_tag WHERE gt_name LIKE 'rest-gt-%'",
		)
	}
}
