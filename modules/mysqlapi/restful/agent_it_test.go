package restful

import (
	"fmt"
	"net/http"

	nqmSql "github.com/Cepave/open-falcon-backend/common/db/nqm/testing"
	json "github.com/Cepave/open-falcon-backend/common/json"
	ogko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func inTx(sql ...string) {
	dbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

func sPtr(v string) *string {
	return &v
}

var _ = Describe("Getting NQM agent by id", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
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
	})

	AfterEach(func() {
		inTx(
			`DELETE FROM nqm_agent WHERE ag_id = 36771`,
			`DELETE FROM host WHERE id = 25101`,
			`DELETE FROM owl_name_tag WHERE nt_id = 407`,
			`
			DELETE FROM owl_group_tag
			WHERE gt_id >= 50091 AND
				gt_id <= 50093
			`,
		)
	})

	It("Get data of a NQM agent", func() {
		resp := tHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().
				Get("api/v1/nqm/agent/36771"),
		)

		jsonBody := resp.GetBodyAsJson()
		GinkgoT().Logf("[Get A Agent] JSON Result: %s", json.MarshalPrettyJSON(jsonBody))

		Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
		Expect(jsonBody.Get("id").MustInt()).To(Equal(36771))
		Expect(jsonBody.Get("group_tags").MustArray()).To(HaveLen(3))
	})
}))

var _ = Describe("Adding new NQM agent", itSkip.PrependBeforeEach(func() {
	AfterEach(func() {
		inTx(
			"DELETE FROM nqm_agent WHERE ag_connection_id LIKE 'add-agent%'",
			"DELETE FROM host WHERE hostname = 'new-host-cccc'",
			"DELETE FROM owl_name_tag where nt_value LIKE 'add-agent-%'",
			"DELETE FROM owl_group_tag where gt_name LIKE 'pp-rest-tag-%'",
		)
	})

	reqBody := &struct {
		Name         string   `json:"name"`
		Comment      string   `json:"comment"`
		Hostname     string   `json:"hostname"`
		ConnectionId string   `json:"connection_id"`
		Status       bool     `json:status`
		IspId        int      `json:"isp_id"`
		ProvinceId   int      `json:"province_id"`
		CityId       int      `json:"city_id"`
		NameTag      *string  `json:"name_tag"`
		GroupTags    []string `json:"group_tags"`
	}{
		Name:       "ko-name-cc1",
		Comment:    "cc-name-cc1",
		Hostname:   "new-host-cccc",
		Status:     true,
		IspId:      8,
		ProvinceId: 9,
		CityId:     130,
		GroupTags:  []string{"pp-rest-tag-1", "pp-rest-tag-2"},
	}

	callApi := func(jsonBody interface{}) *tHttp.ResponseResult {
		return tHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().Post("api/v1/nqm/agent").
				BodyJSON(reqBody),
		)
	}

	DescribeTable("Success adding",
		func(connectionId string, nameTag *string) {
			reqBody.ConnectionId = "add-agent-" + connectionId
			reqBody.NameTag = nameTag

			respResult := callApi(reqBody)

			jsonBody := respResult.GetBodyAsJson()
			GinkgoT().Logf("[Add Agent] JSON Result: %s", json.MarshalPrettyJSON(jsonBody))

			Expect(jsonBody.Get("name").MustString()).To(Equal(reqBody.Name))
			Expect(jsonBody.Get("comment").MustString()).To(Equal(reqBody.Comment))
			Expect(jsonBody.Get("connection_id").MustString()).To(Equal(reqBody.ConnectionId))
			Expect(jsonBody.Get("ip_address").MustString()).To(Equal("0.0.0.0"))
			Expect(jsonBody.Get("hostname").MustString()).To(Equal(reqBody.Hostname))
			Expect(jsonBody.Get("status").MustBool()).To(Equal(reqBody.Status))
			Expect(jsonBody.Get("isp").Get("id").MustInt()).To(Equal(reqBody.IspId))
			Expect(jsonBody.Get("province").Get("id").MustInt()).To(Equal(reqBody.ProvinceId))
			Expect(jsonBody.Get("city").Get("id").MustInt()).To(Equal(reqBody.CityId))

			if nameTag != nil {
				Expect(jsonBody.Get("name_tag").Get("value").MustString()).To(Equal(*nameTag))
			} else {
				Expect(jsonBody.Get("name_tag").Get("id").MustInt()).To(Equal(-1))
			}

			Expect(jsonBody.Get("group_tags").MustArray()).To(HaveLen(2))
		},
		Entry("Viable name tag", "@192.33.9.1", sPtr("add-agent-nt-1")),
		Entry("Nil name tag", "@192.33.9.2", nil),
	)

	It("Adding error with conflict connection id", func() {
		reqBody.ConnectionId = "add-agent-@existing-gc01"
		reqBody.NameTag = nil

		By("Prepare existing data")
		prepareResult := callApi(reqBody)
		Expect(prepareResult).To(ogko.MatchHttpStatus(http.StatusOK))

		By("Asserts the conflict result")
		conflictResult := callApi(reqBody)
		Expect(conflictResult).To(ogko.MatchHttpStatus(http.StatusConflict))

		jsonBody := conflictResult.GetBodyAsJson()
		GinkgoT().Logf("Conflict body: %s", json.MarshalPrettyJSON(jsonBody))
		Expect(jsonBody.Get("error_code").MustInt()).To(Equal(1))
	})
}))

var _ = Describe("Listing agents", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
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
	})

	AfterEach(func() {
		inTx(
			"DELETE FROM nqm_agent WHERE ag_id >= 4321 AND ag_id <= 4323",
			"DELETE FROM host WHERE id = 22091",
		)
	})

	It("Listing without any conditions", func() {
		result := tHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().Get("api/v1/nqm/agents"),
		)
		Expect(result).To(ogko.MatchHttpStatus(http.StatusOK))

		jsonBody := result.GetBodyAsJson()
		GinkgoT().Logf("[List Agents] JSON Result: %s", json.MarshalPrettyJSON(jsonBody))

		Expect(jsonBody.MustArray()).To(HaveLen(3))
	})
}))

var _ = Describe("Listing targets of agent(ping list)", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(nqmSql.InitNqmCacheAgentPingList...)
	})

	AfterEach(func() {
		inTx(nqmSql.ClearNqmCacheAgentPingList...)
	})

	pInt64 := func(v int64) *int64 { return &v }
	fetchTargets := func(agentId int) *tHttp.ResponseResult {
		return tHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().Get(
				fmt.Sprintf("api/v1/nqm/agent/%d/targets", agentId),
			),
		)
	}

	DescribeTable("Normal NQM agent(existing)",
		func(agentId int, expectedRefreshTime *int64) {
			respResult := fetchTargets(agentId)

			Expect(respResult).To(ogko.MatchHttpStatus(http.StatusOK))

			jsonBody := respResult.GetBodyAsJson()
			GinkgoT().Logf("List of targets(JSON): %s", json.MarshalPrettyJSON(jsonBody))

			if expectedRefreshTime != nil {
				Expect(jsonBody.Get("cache_refresh_time").MustInt64()).To(Equal(*expectedRefreshTime))
			} else {
				Expect(jsonBody.Get("cache_refresh_time").Interface()).To(BeNil())
			}
		},
		Entry("Has cache", 24021, pInt64(1483200000)),
		Entry("Has cache", 24022, pInt64(1483200000)),
		Entry("Has no cache", 24023, nil),
	)

	It("Not existing NQM agent", func() {
		respResult := fetchTargets(99801)

		Expect(respResult).To(ogko.MatchHttpStatus(http.StatusNotFound))

		jsonBody := respResult.GetBodyAsJson()
		GinkgoT().Logf("Error content: %s", json.MarshalPrettyJSON(jsonBody))
		Expect(jsonBody.Get("error_code").MustInt()).To(Equal(-1))
	})
}))

var _ = Describe("Clearing cache of target(ping list) on a agent", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(nqmSql.InitNqmCacheAgentPingList...)
	})

	AfterEach(func() {
		inTx(nqmSql.ClearNqmCacheAgentPingList...)
	})

	callApi := func(agentId int) *tHttp.ResponseResult {
		return tHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().Post(
				fmt.Sprintf("api/v1/nqm/agent/%d/targets/clear", agentId),
			),
		)
	}

	DescribeTable("Normal agents",
		func(agentId int, expectedRowsAffected int) {
			callAndCheck := func(expectedAffectedRow int) {
				result := callApi(agentId)

				Expect(result).To(ogko.MatchHttpStatus(http.StatusOK))
				jsonBody := result.GetBodyAsJson()
				GinkgoT().Logf("Clear cache: %s", json.MarshalPrettyJSON(jsonBody))
				Expect(jsonBody.Get("rows_affected").MustInt()).To(Equal(expectedAffectedRow))
			}

			By("First time of clearing")
			callAndCheck(expectedRowsAffected)

			By("Second time of clearing")
			callAndCheck(0)
		},
		Entry("Cache has viable targets", 24021, 1),
		Entry("Cache has non-viable target", 24022, 1),
		Entry("No cache for the agent", 24023, 0),
	)

	It("Clear non-existing NQM agent", func() {
		result := callApi(90051)

		Expect(result).To(ogko.MatchHttpStatus(http.StatusNotFound))
		jsonBody := result.GetBodyAsJson()
		GinkgoT().Logf("Error result: %s", json.MarshalPrettyJSON(jsonBody))
		Expect(jsonBody.Get("error_code").MustInt()).To(Equal(-1))
	})
}))

var _ = Describe("Modifying NQM agent", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
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
	})

	AfterEach(func() {
		inTx(
			"DELETE FROM nqm_agent WHERE ag_id = 23041",
			"DELETE FROM host WHERE id = 4401",
			"DELETE FROM owl_name_tag WHERE nt_value LIKE 'rest-nt-%'",
			"DELETE FROM owl_group_tag WHERE gt_name LIKE 'rest-gt-%'",
		)
	})

	reqJson := &struct {
		Name       string   `json:"name"`
		Status     bool     `json:status`
		Comment    string   `json:"comment"`
		IspId      int      `json:"isp_id"`
		ProvinceId int      `json:"province_id"`
		CityId     int      `json:"city_id"`
		NameTag    *string  `json:"name_tag"`
		GroupTags  []string `json:"group_tags"`
	}{
		Name:       "Update-Agent-1",
		Status:     false,
		Comment:    "This is updated comment",
		IspId:      3,
		ProvinceId: 11,
		CityId:     230,
	}

	DescribeTable("",
		func(nameTag *string, groupTags []string) {
			reqJson.NameTag = nameTag
			reqJson.GroupTags = groupTags

			result := tHttp.NewResponseResultBySling(
				httpClientConfig.NewClient().Put(
					fmt.Sprintf("api/v1/nqm/agent/23041"),
				).BodyJSON(reqJson),
			)

			Expect(result).To(ogko.MatchHttpStatus(http.StatusOK))

			jsonBody := result.GetBodyAsJson()
			GinkgoT().Logf("Update agent: %s", json.MarshalPrettyJSON(jsonBody))

			Expect(jsonBody.Get("name").MustString()).To(Equal(reqJson.Name))
			Expect(jsonBody.Get("status").MustBool()).To(Equal(reqJson.Status))
			Expect(jsonBody.Get("isp").Get("id").MustInt()).To(Equal(reqJson.IspId))
			Expect(jsonBody.Get("province").Get("id").MustInt()).To(Equal(reqJson.ProvinceId))
			Expect(jsonBody.Get("city").Get("id").MustInt()).To(Equal(reqJson.CityId))

			if nameTag != nil {
				Expect(jsonBody.Get("name_tag").Get("value").MustString()).To(Equal(*nameTag))
			} else {
				Expect(jsonBody.Get("name_tag").Get("id").MustInt()).To(Equal(-1))
			}

			Expect(jsonBody.Get("group_tags").MustArray()).To(HaveLen(len(groupTags)))
		},
		Entry("Set name tag, group tags to viable value", sPtr("rest-nt-9"), []string{"rest-gt-91", "rest-gt-92", "rest-gt-93"}),
		Entry("Set name tag, group tags to non-viable value", nil, []string{}),
	)
}))
