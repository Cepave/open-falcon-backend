package http

import (
	"fmt"
	"net/http"
	"github.com/dghubble/sling"
	httpT "github.com/Cepave/open-falcon-backend/common/testing/http"
	t "github.com/Cepave/open-falcon-backend/modules/query/test"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	sjson "github.com/bitly/go-simplejson"
	. "gopkg.in/check.v1"
)

type TestNqmItSuite struct{}

var _ = Suite(&TestNqmItSuite{})

var httpClientConfig = httpT.NewHttpClientConfigByFlag()

// Tests the building of ICMP query
func (suite *TestNqmItSuite) TestBuildQueryOfIcmp(c *C) {
	testCases := []*struct {
		metricFilter string
		expectedStatus int
	} {
		{ "$min >= 55", http.StatusOK },
		{ "$min >= invalid", http.StatusBadRequest },
	}

	var loadFunc = func(metricFilter string, expectedStatus int) *sjson.Json {
		jsonBody := ojson.RawJsonForm(fmt.Sprintf(
			 // Uses default conditions
			`
			{
				"filters": {
					"time": {
						"start_time": 209807000,
						"end_time": 209847000
					},
					"metrics": "%s"
				}
			}
			`, metricFilter,
		))

		return httpT.NewCheckSlint(
			c,
			sling.New().
				Post(httpClientConfig.String()).
				Path("/nqm/icmp/compound-report").
				BodyJSON(jsonBody),
		).
			GetJsonBody(expectedStatus)
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedJson := loadFunc(testCase.metricFilter, testCase.expectedStatus)
		testedJsonString, _ := testedJson.MarshalJSON()

		c.Logf("JSON Response: %s", testedJsonString)

		switch testCase.expectedStatus {
		case http.StatusOK:
			c.Assert(testedJson.Get("query_id").MustString(), HasLen, 36, comment)
		case http.StatusBadRequest:
			c.Assert(testedJson.Get("error_code").MustInt(), Equals, 1, comment)
		default:
			c.Fail()
		}
	}
}
// Tests the getting of content for a query by UUID
func (suite *TestNqmItSuite) TestGetQueryContentOfIcmp(c *C) {
	testCases := []*struct {
		sampleUuid string
		expectedStatus int
	} {
		{ "1B57D5EA-A02D-4942-871E-C9AFC2916BB1", http.StatusOK },
		{ "046084CD-05F7-09C8-09E9-0209CF8D0543", http.StatusNotFound },
		{ "046084CD-05F7-09C8-09E9-0209CF8D0543", http.StatusNotFound }, // Repeats the non-existing query
		{ "zzoo--cc", http.StatusNotFound },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		/**
		 * Calls RESTful service and retrieve data
		 */
		slingObject := sling.New().Base(httpClientConfig.String()).
				Path("/nqm/icmp/compound-report/query/").
				Path(testCase.sampleUuid)

		clientChecker := httpT.NewCheckSlint(c, slingObject)

		jsonBody := clientChecker.GetJsonBody(testCase.expectedStatus)

		jsonString, _ := jsonBody.MarshalJSON()
		// :~)

		c.Logf("JSON Response: %s", jsonString)

		switch testCase.expectedStatus {
		case http.StatusOK:
			/**
			 * Asserts the normal(200) result
			 */
			c.Assert(
				jsonBody.GetPath("output", "metrics").MustArray(),
				HasLen,
				5,
				comment,
			)
			// :~)
		case http.StatusNotFound:
			/**
			 * Asserts the 404 result
			 */
			c.Assert(jsonBody.Get("error_code").MustInt(), Equals, 1, comment)
			c.Assert(jsonBody.Get("uri").MustString(), Matches, ".*" + testCase.sampleUuid + ".*", comment)
			// :~)
		default:
			c.Fail()
		}
	}
}

func (s *TestNqmItSuite) SetUpTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNqmItSuite.TestGetQueryContentOfIcmp":
		inTx(
			`
			INSERT INTO owl_query(qr_uuid, qr_named_id, qr_md5_content, qr_content, qr_time_access, qr_time_creation)
			VALUES(
				x'1B57D5EAA02D4942871EC9AFC2916BB1',
				'nqm.compound.report',
				x'CE402A1EB8126D0111CD4DF80ECC114C',
				x'B450CB6EC4200CFC179F73C8567D44FB2B518410A1D45262233069AB15FF5E48A2B2957ADD8B61C68F19FB06EFB8880D11AE37105C6D7DA3E820EA4097A7D797B76118FABE034BF31D3B5C9EFBCA0A2BE2CFDA970805AE001D6C7A49A5ACCFB903ED2C494D93AEADE3D4C107476908BDD2F31C6C8C07364C648D2093C2F92C89BEFCCFBC0FBC2119DB1883F2DD509DAC44BBC6B8C0C9DF53C556D9D1D97F7C3D4C70B512D0C47AA07C5620B96AE0BCD07838F97390260E872AECDEA089C2D476D97538894FFB62BF8A23ACFAABF4AC4825EACD95B8F03EDE702ADA53CE3F000000FFFF010000FFFF',
				NOW(), NOW()
			)
			`,
		)
	}
}
func (s *TestNqmItSuite) TearDownTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestNqmItSuite.TestGetQueryContentOfIcmp":
		inTx(
			"DELETE FROM owl_query WHERE qr_named_id = 'nqm.compound.report'",
		)
	case "TestNqmItSuite.TestBuildQueryOfIcmp":
		inTx(
			"DELETE FROM owl_query WHERE qr_named_id = 'nqm.compound.report'",
		)
	}
}

func (s *TestNqmItSuite) SetUpSuite(c *C) {
	t.InitDb(c)
}
func (s *TestNqmItSuite) TearDownSuite(c *C) {
	t.ReleaseDb(c)
}
