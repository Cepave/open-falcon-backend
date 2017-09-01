package owl

import (
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/Cepave/open-falcon-backend/common/db"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"
	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"
	"github.com/Cepave/open-falcon-backend/common/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()

var _ = Describe("[RESTful client] Query Object", func() {
	mysqlApiConfig := gockConfig.NewRestfulClientConfig()

	testedSrv := NewQueryService(
		QueryServiceConfig{mysqlApiConfig},
	)

	AfterEach(func() {
		gockConfig.Off()
	})

	Context("Create or load query", func() {
		sampleReqBody := map[string]interface{}{
			"feature_name": "easy.cool.f1",
			"content":      "kvi4/PqaX8dVXMDwdkilAw==",
			"md5_content":  "kvi4/PqaX8dVXMDwdkilAw==",
		}
		sampleUuid := "1b9b9cf3-b6c6-46ca-8cf4-b5a7d115b932"
		sampleTime := 1426032000

		BeforeEach(func() {
			gockConfig.New().
				MatchType("json").
				JSON(sampleReqBody).
				Post("/api/v1/owl/query-object").
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"uuid":          sampleUuid,
					"creation_time": sampleTime,
					"access_time":   sampleTime,
				})
		})

		It("The UUID, creation time, and access time match expected one", func() {
			var content = &types.VarBytes{}
			content.MustFromBase64(sampleReqBody["content"].(string))

			var md5Content = &types.Bytes16{}
			md5Content.MustFromBase64(sampleReqBody["md5_content"].(string))

			testedQuery := &model.Query{
				NamedId:    "easy.cool.f1",
				Content:    []byte(*content),
				Md5Content: db.Bytes16(*md5Content),
			}
			testedSrv.CreateOrLoadQuery(testedQuery)

			uuidValue, _ := uuid.FromString(sampleUuid)
			expectedTime := time.Unix(int64(sampleTime), 0)
			Expect(testedQuery).To(
				PointTo(
					MatchFields(IgnoreExtras, Fields{
						"Uuid":         BeEquivalentTo(uuidValue),
						"CreationTime": BeTemporally("==", expectedTime),
						"AccessTime":   BeTemporally("==", expectedTime),
					}),
				),
			)
		})
	})

	Context("Get query by UUID", func() {
		Context("Get existing query object", func() {
			uuidString := "bd747923-0ae7-44c2-8e45-8f8e3cde7bce"
			expectedJson := map[string]interface{}{
				"uuid":          uuidString,
				"feature_name":  "get.uuid.f1",
				"content":       "JJ68Cik12gCBfkqtyKV5Yg==",
				"md5_content":   "cCy6iS7FzYN2otLGBHAmCA==",
				"creation_time": 38976511, // [UNIX TIME]
				"access_time":   72808242, // [UNIX TIME]
			}

			BeforeEach(func() {
				gockConfig.New().
					Get("/api/v1/owl/query-object/" + uuidString).
					Reply(http.StatusOK).
					JSON(expectedJson)
			})

			It("Query content matching expected one", func() {
				uuidValue, _ := uuid.FromString(uuidString)
				creationTime := time.Unix(int64(expectedJson["creation_time"].(int)), 0)
				accessTime := time.Unix(int64(expectedJson["access_time"].(int)), 0)

				testedQuery := testedSrv.LoadQueryByUuid(uuidValue)

				var content = &types.VarBytes{}
				content.MustFromBase64(expectedJson["content"].(string))

				var md5Content = &types.Bytes16{}
				md5Content.MustFromBase64(expectedJson["md5_content"].(string))

				Expect(testedQuery).To(
					PointTo(
						MatchFields(IgnoreExtras, Fields{
							"Uuid":         BeEquivalentTo(uuidValue),
							"NamedId":      Equal(expectedJson["feature_name"]),
							"Content":      BeEquivalentTo(*content),
							"Md5Content":   BeEquivalentTo(*md5Content),
							"CreationTime": BeTemporally("==", creationTime),
							"AccessTime":   BeTemporally("==", accessTime),
						}),
					),
				)
			})
		})

		Context("Get non-existing query object", func() {
			uuidString := "82fd6c62-c43f-40a1-b6da-0db7eccf0015"

			BeforeEach(func() {
				gockConfig.New().
					Get("/api/v1/owl/query-object/" + uuidString).
					Reply(http.StatusNotFound).
					JSON(map[string]interface{}{
						"uuid":        uuidString,
						"http_status": http.StatusNotFound,
						"error_code":  1,
						"uri":         "/owl/query-object/" + uuidString,
					})
			})

			It("Query object should be nil", func() {
				uuidValue, _ := uuid.FromString(uuidString)
				testedQuery := testedSrv.LoadQueryByUuid(uuidValue)

				Expect(testedQuery).To(BeNil())
			})
		})
	})
})
