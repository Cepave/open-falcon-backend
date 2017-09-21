package owl

import (
	"fmt"
	"net/http"
	"time"

	client "github.com/Cepave/open-falcon-backend/common/http/client"
	json "github.com/Cepave/open-falcon-backend/common/json"
	ogko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var gtResp = client.ToGentlemanResp

var _ = Describe("[POST] on /owl/query-object", itSkip.PrependBeforeEach(func() {
	AfterEach(func() {
		inTx(
			`
			DELETE FROM owl_query
			WHERE qr_named_id = 'query.object.it'
			`,
		)
	})

	It("Check generated uuid and content. And the access/creation time must be after a certain time of boundary", func() {
		now := time.Now().Unix()

		resp, err := httpClient.NewClient().
			Post().AddPath("/api/v1/owl/query-object").
			JSON(
				map[string]interface{}{
					"feature_name": "query.object.it",
					"content":      "O9Hlrlg4LWGjDfs9/VgjEbuCMPu99xJc3V3btld4TVw=",
					"md5_content":  "HheztFtq6WHjMgYKOjqB4g==",
				},
			).
			Send()

		Expect(err).To(Succeed())

		jsonBody := gtResp(resp).MustGetJson()

		GinkgoT().Logf("[/owl/query-object] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))

		Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
		Expect(jsonBody.Get("uuid").MustString()).To(MatchRegexp("^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"))
		Expect(jsonBody.Get("feature_name").MustString()).To(Equal("query.object.it"))
		Expect(jsonBody.Get("content").MustString()).To(Equal("O9Hlrlg4LWGjDfs9/VgjEbuCMPu99xJc3V3btld4TVw="))
		Expect(jsonBody.Get("md5_content").MustString()).To(Equal("HheztFtq6WHjMgYKOjqB4g=="))
		Expect(jsonBody.Get("access_time").MustInt64()).To(BeNumerically(">=", now))
		Expect(jsonBody.Get("creation_time").MustInt64()).To(BeNumerically(">=", now))
	})
}))

var _ = Describe("[GET] on /owl/query-object/:uuid", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(`
		INSERT INTO owl_query(
			qr_uuid, qr_named_id,
			qr_content, qr_md5_content,
			qr_time_creation, qr_time_access
		)
		VALUES(
			x'9ed424fb656ba68296132c2b5458f4ac', 'query.object.gg1',
			x'd841ccfe81ca02e1cceaa1e4cd8fc523', x'2070faa9bf90d6dfc91b41bedfe2500f',
			'2010-05-03', '2010-06-13'
		)`)
	})

	AfterEach(func() {
		inTx(
			`
			DELETE FROM owl_query
			WHERE qr_named_id = 'query.object.gg1'
			`,
		)
	})

	Context("Fetch existing query object", func() {
		It("Checks the named id", func() {
			resp, err := httpClient.NewClient().
				Get().AddPath("/api/v1/owl/query-object/9ed424fb-656b-a682-9613-2c2b5458f4ac").
				JSON(
					map[string]interface{}{
						"feature_name": "query.object.it",
						"content":      "O9Hlrlg4LWGjDfs9/VgjEbuCMPu99xJc3V3btld4TVw=",
						"md5_content":  "HheztFtq6WHjMgYKOjqB4g==",
					},
				).
				Send()

			Expect(err).To(Succeed())

			jsonBody := gtResp(resp).MustGetJson()
			GinkgoT().Logf("[/owl/query-object/:uuid] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))

			Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
			Expect(jsonBody.Get("feature_name").MustString()).To(Equal("query.object.gg1"))
		})
	})

	Context("Fetch non-existing query object", func() {
		It("Checks the 404 status", func() {
			resp, err := httpClient.NewClient().
				Get().AddPath("/api/v1/owl/query-object/70d420fb-6a6b-a682-9613-2c2b5458f4ac").
				JSON(
					map[string]interface{}{
						"feature_name": "query.object.it",
						"content":      "O9Hlrlg4LWGjDfs9/VgjEbuCMPu99xJc3V3btld4TVw=",
						"md5_content":  "HheztFtq6WHjMgYKOjqB4g==",
					},
				).
				Send()

			Expect(err).To(Succeed())

			jsonBody := gtResp(resp).MustGetJson()
			GinkgoT().Logf("[/owl/query-object/:uuid] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))

			Expect(resp).To(ogko.MatchHttpStatus(http.StatusNotFound))
			Expect(jsonBody.Get("http_status").MustInt()).To(Equal(404))
			Expect(jsonBody.Get("uuid").MustString()).To(Equal("70d420fb-6a6b-a682-9613-2c2b5458f4ac"))
		})
	})
}))

var _ = Describe("[POST] on /owl/query-object/vacuum", itSkip.PrependBeforeEach(func() {
	Context("Remove older query objects", func() {
		BeforeEach(func() {
			inTx(`
				INSERT INTO owl_query(
					qr_uuid, qr_named_id,
					qr_content, qr_md5_content,
					qr_time_access, qr_time_creation
				)
				VALUES
				(
					x'209f18f4f89b42568e1e5270987c057d', 'vacuum.f1',
					x'7011e902d4a848c184e242e8d71aa961', x'bdfa89f3df204071b2f48b359440e0fc',
					DATE_SUB(NOW(), INTERVAL 65 DAY), DATE_SUB(NOW(), INTERVAL 65 DAY)
				),
				(
					x'349f18f4f89b42568e1e5270987c057d', 'vacuum.f1',
					x'7011e902d4a848c184e242e8d71aa961', x'1afaa9f3df2d4071b2f48b359440e0fc',
					DATE_SUB(NOW(), INTERVAL 65 DAY), DATE_SUB(NOW(), INTERVAL 65 DAY)
				),
				(
					x'109f1cf4f89b42568e1e5270987c057d', 'vacuum.f2',
					x'7011e902d4a848c184e242e8d71aa961', x'12faa9f3df2d4071b9f48b359440e0fc',
					DATE_SUB(NOW(), INTERVAL 65 DAY), DATE_SUB(NOW(), INTERVAL 65 DAY)
				)
			`)
		})

		AfterEach(func() {
			inTx(`
				DELETE FROM owl_query
				WHERE qr_named_id IN ('vacuum.f1', 'vacuum.f2')
			`)
		})

		DescribeTable("Check number of affected rows",
			func(forDays int, expectedAffectedRows int) {
				resp, err := httpClient.NewClient().
					Post().AddPath("/api/v1/owl/query-object/vacuum").
					AddQuery("for_days", fmt.Sprintf("%d", forDays)).
					Send()

				Expect(err).To(Succeed())

				jsonBody := gtResp(resp).MustGetJson()
				GinkgoT().Logf("[/owl/query-object/vacuum] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))

				Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
				Expect(jsonBody.Get("affected_rows").MustInt()).To(BeEquivalentTo(expectedAffectedRows))
			},
			Entry("Vacuum is affecting data", 30, 3),
			Entry("Vacuum is not affecting data", 120, 0),
		)
	})
}))
