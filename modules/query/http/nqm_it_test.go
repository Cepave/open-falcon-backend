package http

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/satori/go.uuid"

	oHttp "github.com/Cepave/open-falcon-backend/common/http/client"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	tg "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var itConfig = tHttp.NewHttpClientConfigByFlag()
var itClient = tHttp.GentlemanClientConf{itConfig}

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()
var fakeServerConfig = &tHttp.FakeServerConfig{"127.0.0.1", 6040}

var _ = Describe("[POST] on /nqm/icmp/compound-report", func() {
	skipItOnMySqlApi.BeforeEachSkip()

	var (
		mockMysqlApiServer *httptest.Server
		currentUuid        string
	)

	BeforeEach(func() {
		mockMysqlApiServer = gockConfig.HttpTest.NewServer(fakeServerConfig)
		mockMysqlApiServer.Start()
	})
	AfterEach(func() {
		mockMysqlApiServer.Close()
		mockMysqlApiServer = nil
		gockConfig.Off()
	})

	Context("Successful building of query", func() {
		BeforeEach(func() {
			var (
				currentContent    string
				currentMd5Content string
			)

			currentUuid = uuid.NewV4().String()
			gockConfig.New().Post("/owl/query-object").
				MatchType("json").
				Filter(
					func(req *http.Request) bool {
						json := ojson.UnmarshalToJson(req.Body)

						/**
						 * Loads fed binary data and output as same as source
						 */
						currentContent = json.Get("content").MustString()
						currentMd5Content = json.Get("md5_content").MustString()
						GinkgoT().Logf("Content: %s", currentContent)
						GinkgoT().Logf("Md5 Content: %s", currentMd5Content)
						// :~)

						return true
					},
				).
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"uuid":          currentUuid,
					"feature_name":  "nqm.compound.report",
					"content":       currentContent,
					"md5_content":   currentMd5Content,
					"creation_time": time.Now().Unix(),
					"access_time":   time.Now().Unix(),
				})
		})

		DescribeTable("The UUID should be same as generated",
			func(queryBody map[string]interface{}) {
				resp, err := itClient.NewClient().Path("/nqm/icmp/compound-report").
					Post().
					JSON(queryBody).
					Send()

				Expect(err).To(Succeed())
				defer resp.Close()

				Expect(resp).To(tg.MatchHttpStatus(http.StatusOK))

				jsonBody := oHttp.ToGentlemanResp(resp).MustGetJson()

				jsonText, _ := jsonBody.MarshalJSON()
				GinkgoT().Logf("JSON Result: %s", jsonText)

				Expect(jsonBody.Get("query_id").MustString()).To(Equal(currentUuid))
			},
			Entry("Get a new UUID by empty JSON body: \"{}\"", map[string]interface{}{}),
			Entry(`Get a new UUID by filter of agent: "connection_id": []string{ "conn-id-1", "conn-id-2" }`,
				map[string]interface{}{
					"filters": map[string]interface{}{
						"target": map[string]interface{}{
							"host": []string{"pc-01", "pc-02"},
						},
						//"agent": map[string]interface{} {
						//"connection_id": []string{ "conn-id-1", "conn-id-2" },
						//},
					},
				},
			),
		)
	})

	Context("Error of DSL(metric filter)", func() {
		It("Get status of \"Bad Request(400)\" and < \"error_code\": 1 >(JSON)", func() {
			resp, err := itClient.NewClient().Path("/nqm/icmp/compound-report").
				Post().
				JSON(map[string]interface{}{
					"filters": map[string]interface{}{
						"metrics": "$min >= ccl",
					},
				}).Send()

			Expect(err).To(Succeed())
			defer resp.Close()

			Expect(resp).To(tg.MatchHttpStatus(http.StatusBadRequest))

			jsonBody := oHttp.ToGentlemanResp(resp).MustGetJson()

			jsonText, _ := jsonBody.MarshalJSON()
			GinkgoT().Logf("JSON Result for error: %s", jsonText)

			Expect(jsonBody.Get("error_code").MustInt()).To(BeEquivalentTo(1))
		})
	})
})

var _ = Describe("[GET] on /nqm/icmp/compound-report/query/{uuid}", func() {
	skipItOnMySqlApi.BeforeEachSkip()

	var (
		mockMysqlApiServer *httptest.Server
		currentUuid        string
	)

	BeforeEach(func() {
		mockMysqlApiServer = gockConfig.HttpTest.NewServer(fakeServerConfig)
		mockMysqlApiServer.Start()
	})
	AfterEach(func() {
		mockMysqlApiServer.Close()
		mockMysqlApiServer = nil
		gockConfig.Off()
	})

	Context("Get viable query", func() {
		BeforeEach(func() {
			currentUuid = uuid.NewV4().String()
			gockConfig.New().Get("/api/v1/owl/query-object/" + currentUuid).
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"uuid":          currentUuid,
					"feature_name":  "nqm.compound.report",
					"content":       "tJDBbsQgDET/Zc5U6vaYX4lWCBHKWkoMApO2WvHvFSQtXanXPaGxB7+x73inVVzKmO4Q2lx7s5gk+lBc1lXB8fKgJWgOH81bmAQTblDYzVocpkutCsY7ltZn037NV4VbyDIURW2WJbmcD20Ds7NCgTUtpyVHTcvZjynsxNaNiiX5GqpN1mL8qPgUSvxbqgpiknf/5MI0I9qX1wtUf9/wHP7mJJHNmIB6Ooh9y3MebD6CPdxnwHFQ0bNhQFvan9U6JxSJpe/5S5yxmU8obMRQMLuHwhr6eBsKC661fgMAAP//AQAA//8=",
					"md5_content":   "iyunvQkZgMBp2hxsm4x5Kw==",
					"creation_time": time.Now().Unix(),
					"access_time":   time.Now().Unix(),
				})
		})

		It(`Filter of target matches host on ["pc01", "pc02"]`, func() {
			resp, err := itClient.NewClient().Path("/nqm/icmp/compound-report/query/" + currentUuid).
				Get().
				Send()

			Expect(err).To(Succeed())
			defer resp.Close()

			Expect(resp).To(tg.MatchHttpStatus(http.StatusOK))

			jsonBody := oHttp.ToGentlemanResp(resp).MustGetJson()

			jsonText, _ := jsonBody.MarshalJSON()
			GinkgoT().Logf("JSON Result: %s", jsonText)

			Expect(jsonBody.GetPath("target", "host").MustStringArray()).To(Equal([]string{"pc-01", "pc-02"}))
		})
	})

	Context("Get non-existing query", func() {
		BeforeEach(func() {
			currentUuid = uuid.NewV4().String()
			gockConfig.New().Get("/api/v1/owl/query-object/" + currentUuid).
				Reply(http.StatusNotFound).
				JSON(map[string]interface{}{
					"uuid":        currentUuid,
					"http_status": http.StatusNotFound,
					"error_code":  1,
					"uri":         "/owl/query-object/" + currentUuid,
				})
		})

		performTest := func() {
			resp, err := itClient.NewClient().Path("/nqm/icmp/compound-report/query/" + currentUuid).
				Get().
				Send()

			Expect(err).To(Succeed())
			defer resp.Close()

			Expect(resp).To(tg.MatchHttpStatus(http.StatusNotFound))

			jsonBody := oHttp.ToGentlemanResp(resp).MustGetJson()

			jsonText, _ := jsonBody.MarshalJSON()
			GinkgoT().Logf("JSON Result for error: %s", jsonText)

			Expect(jsonBody.GetPath("uri").MustString()).To(ContainSubstring(currentUuid))
			Expect(jsonBody.GetPath("error_code").MustInt()).To(BeEquivalentTo(1))
		}

		It(`Should got 404(Not Found)`, performTest)
		It(`Should got 404(Not Found) - again`, performTest)
	})
})
