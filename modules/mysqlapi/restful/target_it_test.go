package restful

import (
	"net/http"

	json "github.com/Cepave/open-falcon-backend/common/json"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	testingDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/testing"

	rdb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"

	. "gopkg.in/check.v1"
)

type TestTargetItSuite struct{}

var _ = Suite(&TestTargetItSuite{})

// Tests the getting of agent by id
func (suite *TestTargetItSuite) TestGetTargetById(c *C) {
	client := httpClientConfig.NewSlingByBase().Get(
		"api/v1/nqm/target/40021",
	)

	slintChecker := testingHttp.NewCheckSlint(c, client)
	jsonResult := slintChecker.GetJsonBody(http.StatusOK)

	c.Logf("[Get A Target] JSON Result: %s", json.MarshalPrettyJSON(jsonResult))
	c.Assert(jsonResult.Get("id").MustInt(), Equals, 40021)
}

// Tests the adding of new target
func (suite *TestTargetItSuite) TestAddNewTarget(c *C) {
	jsonBody := &struct {
		Name        string   `json:"name"`
		Host        string   `json:"host"`
		Status      bool     `json:"status"`
		ProbedByAll bool     `json:"probed_by_all"`
		Comment     string   `json:"comment"`
		IspId       int      `json:"isp_id"`
		ProvinceId  int      `json:"province_id"`
		CityId      int      `json:"city_id"`
		NameTag     *string  `json:"name_tag"`
		GroupTags   []string `json:"group_tags"`
	}{
		Status:      true,
		ProbedByAll: true,
		Comment:     "This is new target by red 33.72 ***",
		IspId:       2,
		ProvinceId:  27,
		CityId:      206,
	}

	sPtr := func(v string) *string { return &v }
	testCases := []*struct {
		host              string
		nameTag           *string
		groupTags         []string
		expectedStatus    int
		expectedErrorCode int
	}{
		{"54", nil, []string{}, http.StatusOK, -1},
		{"55", sPtr("tg-nt-1"), []string{"tg-rest-tag-1", "tg-rest-tag-2"}, http.StatusOK, -1},
		{"54", nil, []string{}, http.StatusConflict, 1},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		jsonBody.Name = "new-tg-" + testCase.host
		jsonBody.Host = "87.6.44." + testCase.host
		jsonBody.NameTag = testCase.nameTag
		jsonBody.GroupTags = testCase.groupTags

		client := httpClientConfig.NewSlingByBase().
			Post("api/v1/nqm/target").
			BodyJSON(jsonBody)

		slintChecker := testingHttp.NewCheckSlint(c, client)

		jsonResp := slintChecker.GetJsonBody(testCase.expectedStatus)

		c.Logf("[Add Target] JSON Result: %s", json.MarshalPrettyJSON(jsonResp))

		switch testCase.expectedStatus {
		case http.StatusConflict:
			c.Assert(jsonResp.Get("error_code").MustInt(), Equals, testCase.expectedErrorCode, comment)
		}

		if testCase.expectedStatus != http.StatusOK {
			continue
		}

		c.Assert(jsonResp.Get("name").MustString(), Equals, jsonBody.Name, comment)
		c.Assert(jsonResp.Get("host").MustString(), Equals, jsonBody.Host, comment)
		c.Assert(jsonResp.Get("comment").MustString(), Equals, jsonBody.Comment, comment)
		c.Assert(jsonResp.Get("status").MustBool(), Equals, jsonBody.Status, comment)
		c.Assert(jsonResp.Get("probed_by_all").MustBool(), Equals, jsonBody.ProbedByAll, comment)
		c.Assert(jsonResp.Get("isp").Get("id").MustInt(), Equals, jsonBody.IspId, comment)
		c.Assert(jsonResp.Get("province").Get("id").MustInt(), Equals, jsonBody.ProvinceId, comment)
		c.Assert(jsonResp.Get("city").Get("id").MustInt(), Equals, jsonBody.CityId, comment)

		if testCase.nameTag != nil {
			c.Assert(jsonResp.Get("name_tag").Get("value").MustString(), Equals, *jsonBody.NameTag, comment)
		} else {
			c.Assert(jsonResp.Get("name_tag").Get("id").MustInt(), Equals, -1, comment)
		}

		c.Assert(jsonResp.Get("group_tags").MustArray(), HasLen, len(jsonBody.GroupTags), comment)
	}
}

// Tests the modifying of target
func (suite *TestTargetItSuite) TestModifyTarget(c *C) {
	jsonBody := &struct {
		Name        string   `json:"name"`
		Status      bool     `json:"status"`
		ProbedByAll bool     `json:"probed_by_all"`
		Comment     string   `json:"comment"`
		IspId       int      `json:"isp_id"`
		ProvinceId  int      `json:"province_id"`
		CityId      int      `json:"city_id"`
		NameTag     *string  `json:"name_tag"`
		GroupTags   []string `json:"group_tags"`
	}{
		Name:        "Updated-Target-1",
		Status:      false,
		ProbedByAll: false,
		Comment:     "[3981] This is updated target",
		IspId:       9,
		ProvinceId:  19,
		CityId:      164,
	}

	sPtr := func(v string) *string { return &v }
	testCases := []*struct {
		nameTag   *string
		groupTags []string
	}{
		{sPtr("tg-nt-3"), []string{"blue-utg-3", "blue-utg-4", "blue-utg-5"}},
		{nil, []string{}},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		jsonBody.NameTag = testCase.nameTag
		jsonBody.GroupTags = testCase.groupTags

		client := httpClientConfig.NewSlingByBase().
			Put("api/v1/nqm/target/39347").
			BodyJSON(jsonBody)

		slintChecker := testingHttp.NewCheckSlint(c, client)

		jsonResult := slintChecker.GetJsonBody(http.StatusOK)

		c.Logf("Update target: %v", json.MarshalPrettyJSON(jsonResult))

		c.Assert(jsonResult.Get("name").MustString(), Equals, jsonBody.Name, comment)
		c.Assert(jsonResult.Get("host").MustString(), Equals, "ab01.targets.com.cn", comment)
		c.Assert(jsonResult.Get("comment").MustString(), Equals, jsonBody.Comment, comment)
		c.Assert(jsonResult.Get("status").MustBool(), Equals, jsonBody.Status, comment)
		c.Assert(jsonResult.Get("isp").Get("id").MustInt(), Equals, jsonBody.IspId, comment)
		c.Assert(jsonResult.Get("province").Get("id").MustInt(), Equals, jsonBody.ProvinceId, comment)
		c.Assert(jsonResult.Get("city").Get("id").MustInt(), Equals, jsonBody.CityId, comment)

		if jsonBody.NameTag != nil {
			c.Assert(jsonResult.Get("name_tag").Get("value").MustString(), Equals, *jsonBody.NameTag, comment)
		} else {
			c.Assert(jsonResult.Get("name_tag").Get("id").MustInt(), Equals, -1, comment)
		}

		c.Assert(jsonResult.Get("group_tags").MustArray(), HasLen, len(jsonBody.GroupTags), comment)
	}
}

// Tests the listing of targets
func (suite *TestTargetItSuite) TestListTargets(c *C) {
	client := httpClientConfig.NewSlingByBase().
		Get("api/v1/nqm/targets")

	slintChecker := testingHttp.NewCheckSlint(c, client)

	slintChecker.AssertHasPaging()
	message := slintChecker.GetJsonBody(http.StatusOK)

	c.Logf("[List Targets] JSON Result: %s", json.MarshalPrettyJSON(message))
	c.Assert(len(message.MustArray()), Equals, 3)
}

func (s *TestTargetItSuite) SetUpSuite(c *C) {
	testingDb.InitRdb(c)
}
func (s *TestTargetItSuite) TearDownSuite(c *C) {
	testingDb.ReleaseRdb(c)
}

func (s *TestTargetItSuite) SetUpTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestTargetItSuite.TestGetTargetById":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3224, 'tg-nt-1')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(23091, 'tg-gt-1'),
				(23092, 'tg-gt-2'),
				(23093, 'tg-gt-3')
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id,
				tg_comment
			)
			VALUES(40021, 'tg-by-id-1', 'tg-1.byid.com.cn', 8, 27, 202, 3224, 'Comment of Target 1')
			`,
			`
			INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
			VALUES(40021, 23091), (40021, 23092), (40021, 23093)
			`,
		)
	case "TestTargetItSuite.TestListTargets":
		inTx(
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host, tg_status, tg_available
			)
			VALUES(40901, 'tg-name-1', 'tg-1.fastweb.com', true, true),
				(40902, 'tg-name-2', 'tg-2.fastweb.com', true, true),
				(40903, 'tg-name-3', 'tg-3.fastweb.com', true, true)
			`,
		)
	case "TestTargetItSuite.TestModifyTarget":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(10101, 'tg-nt-1')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(30301, "blue-utg-1"), (30302, "blue-utg-2")
			`,
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host, tg_status, tg_probed_by_all,
				tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id, tg_comment
			)
			VALUES(
				39347, 'Old-Target-1', 'ab01.targets.com.cn', true, true,
				2, 18, 41, 10101, 'ABC Old Comment'
			)
			`,
			`
			INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
			VALUES(39347, 30301),(39347, 30302)
			`,
		)
	}
}
func (s *TestTargetItSuite) TearDownTest(c *C) {
	inTx := rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestTargetItSuite.TestGetTargetById":
		inTx(
			`
			DELETE FROM nqm_target
			WHERE tg_id = 40021
			`,
			`
			DELETE FROM owl_name_tag
			WHERE nt_id = 3224
			`,
			`
			DELETE FROM owl_group_tag
			WHERE gt_id >= 23091 AND
				gt_id <= 23093
			`,
		)
	case "TestTargetItSuite.TestListTargets":
		inTx(
			"DELETE FROM nqm_target WHERE tg_id >= 40901 AND tg_id <= 40903",
		)
	case "TestTargetItSuite.TestAddNewTarget":
		inTx(
			"DELETE FROM nqm_target WHERE tg_name LIKE 'new-tg-%'",
			"DELETE FROM owl_name_tag where nt_value = 'tg-nt-1'",
			"DELETE FROM owl_group_tag where gt_name LIKE 'tg-rest-tag-%'",
		)
	case "TestTargetItSuite.TestModifyTarget":
		inTx(
			"DELETE FROM nqm_target WHERE tg_id = 39347",
			"DELETE FROM owl_name_tag WHERE nt_value LIKE 'tg-nt-%'",
			"DELETE FROM owl_group_tag WHERE gt_name LIKE 'blue-utg-%'",
		)
	}
}
