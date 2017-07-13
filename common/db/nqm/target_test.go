package nqm

import (
	"reflect"

	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	"github.com/Cepave/open-falcon-backend/common/utils"

	. "gopkg.in/check.v1"
)

type TestTargetSuite struct{}

var _ = Suite(&TestTargetSuite{})

// Tests the adding of data for target
func (suite *TestTargetSuite) TestAddTarget(c *C) {
	sPtr := func(v string) *string { return &v }

	/**
	 * sample targets
	 */
	sampleTarget := nqmModel.NewTargetForAdding()
	sampleTarget.Host = "13.41.60.178"
	sampleTarget.IspId = 3
	sampleTarget.ProvinceId = 20
	sampleTarget.CityId = 6
	// :~)

	testCases := []*struct {
		targetName string
		host       string
		nameTag    *string
		groupTags  []string
		cityId     int16
		errorType  interface{}
	}{
		{"tg-1", "1", nil, []string{}, 6, nil},
		{"tg-1", "1", nil, []string{}, 6, ErrDuplicatedNqmTarget{}}, // Duplicated host
		{"tg-2-1", "2", sPtr("IBM-1"), []string{"KSC-1", "KSC-32"}, 6, nil},
		{"tg-2-2", "3", sPtr("IBM-1"), []string{"KSC-1", "KSC-32"}, 112, owlDb.ErrNotInSameHierarchy{}}, // Location is not consistent
	}

	for _, testCase := range testCases {
		sampleTarget.Name = "def-target-" + testCase.targetName
		sampleTarget.Host = "13.41.60.3" + testCase.host
		sampleTarget.NameTagValue = testCase.nameTag
		sampleTarget.GroupTags = testCase.groupTags
		sampleTarget.CityId = testCase.cityId

		newTarget, err := AddTarget(sampleTarget)

		/**
		 * Asserts the occurring error
		 */
		if testCase.errorType != nil {
			c.Assert(newTarget, IsNil)
			c.Assert(err, NotNil)
			c.Logf("Has error: %v", err)

			c.Assert(reflect.TypeOf(err), Equals, reflect.TypeOf(testCase.errorType))
			continue
		}
		// :~)

		c.Assert(err, IsNil)
		c.Logf("New Target: %v", newTarget)
		c.Logf("New Target[Group Tags]: %v", newTarget.GroupTags)

		c.Assert(newTarget.Name, Equals, sampleTarget.Name)
		c.Assert(newTarget.Host, Equals, sampleTarget.Host)
		c.Assert(newTarget.IspId, Equals, sampleTarget.IspId)
		c.Assert(newTarget.ProvinceId, Equals, sampleTarget.ProvinceId)
		c.Assert(newTarget.CityId, Equals, sampleTarget.CityId)

		if sampleTarget.NameTagValue != nil {
			c.Assert(newTarget.NameTagValue, Equals, *sampleTarget.NameTagValue)
		} else {
			c.Assert(newTarget.NameTagId, Equals, int16(-1))
		}

		c.Assert(newTarget.GroupTags, HasLen, len(sampleTarget.GroupTags))
	}
}

// Tests the updating for data of target
func (suite *TestTargetSuite) TestUpdateTarget(c *C) {
	sPtr := func(v string) *string { return &v }
	modifiedTarget := &nqmModel.TargetForAdding{
		Name:       "new-tg-1",
		Comment:    utils.PointerOfCloneString("new-comment-1"),
		Status:     false,
		ProvinceId: 26,
		CityId:     194,
		IspId:      2,
	}

	testCases := []*struct {
		nameTag   *string
		groupTags []string
	}{
		{sPtr("tg-nt-2"), []string{"tg-gt-3", "tg-gt-4", "tg-gt-5"}},
		{nil, []string{}},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		modifiedTarget.NameTagValue = testCase.nameTag
		modifiedTarget.GroupTags = testCase.groupTags

		originalTarget := GetTargetById(34091)
		testedTarget, err := UpdateTarget(originalTarget, modifiedTarget)

		c.Assert(err, IsNil)
		c.Assert(testedTarget.Name, Equals, modifiedTarget.Name, comment)
		c.Assert(testedTarget.Comment, DeepEquals, modifiedTarget.Comment, comment)
		c.Assert(testedTarget.Status, Equals, modifiedTarget.Status, comment)
		c.Assert(testedTarget.ProvinceId, Equals, modifiedTarget.ProvinceId, comment)
		c.Assert(testedTarget.CityId, Equals, modifiedTarget.CityId, comment)
		c.Assert(testedTarget.IspId, Equals, modifiedTarget.IspId, comment)

		if modifiedTarget.NameTagValue != nil {
			c.Assert(testedTarget.NameTagValue, Equals, *modifiedTarget.NameTagValue, comment)
		} else {
			c.Assert(testedTarget.NameTagId, Equals, int16(-1), comment)
		}

		testedTargetForAdding := testedTarget.ToTargetForAdding()
		c.Assert(testedTargetForAdding.AreGroupTagsSame(modifiedTarget), Equals, true, comment)
	}
}

// Tests the retrieving of data for a target by id
func (suite *TestTargetSuite) TestGetTargetById(c *C) {
	testCases := []*struct {
		sampleIdOfTarget int32
		hasFound         bool
	}{
		{23041, true},
		{23051, false},
	}

	for i, testCase := range testCases {
		result := GetTargetById(testCase.sampleIdOfTarget)

		if testCase.hasFound {
			c.Logf("Found target by id: %v", result)
			c.Assert(result, NotNil, Commentf("Test Case: %d", i))
		} else {
			c.Assert(result, IsNil, Commentf("Test Case: %d", i))
		}
	}
}

// Tests the listing of targets
func (suite *TestTargetSuite) TestListTargets(c *C) {
	testCases := []*struct {
		query                      *nqmModel.TargetQuery
		pageSize                   int32
		pagePosition               int32
		expectedCountOfCurrentPage int
		expectedCountOfAll         int32
	}{
		{ // All data
			&nqmModel.TargetQuery{IspId: -2, HasStatusParam: false},
			10, 1, 3, 3,
		},
		{ // 2nd page
			&nqmModel.TargetQuery{IspId: -2, HasStatusParam: false},
			2, 2, 1, 3,
		},
		{ // Match nothing for further page
			&nqmModel.TargetQuery{IspId: -2, HasStatusParam: false},
			10, 10, 0, 3,
		},
		{ // Match 1 row by all of the conditions
			&nqmModel.TargetQuery{
				Name:           "tg-name-1",
				Host:           "tg-host-1",
				IspId:          3,
				HasStatusParam: true,
				Status:         true,
			}, 10, 1, 1, 1,
		},
		{ // Match 1 row(by ISP id)
			&nqmModel.TargetQuery{
				IspId:         5,
				HasIspIdParam: true,
			}, 10, 1, 1, 1,
		},
		{ // Match nothing
			&nqmModel.TargetQuery{
				IspId:          -2,
				HasStatusParam: false,
				Name:           "tg-easy-1",
				Host:           "tg-host-1",
			}, 10, 1, 0, 0,
		},
	}

	for i, testCase := range testCases {
		paging := commonModel.Paging{
			Size:     testCase.pageSize,
			Position: testCase.pagePosition,
			OrderBy: []*commonModel.OrderByEntity{
				{"id", commonModel.Descending},
				{"name", commonModel.Ascending},
				{"status", commonModel.Ascending},
				{"host", commonModel.Ascending},
				{"comment", commonModel.Ascending},
				{"isp", commonModel.Ascending},
				{"province", commonModel.Ascending},
				{"city", commonModel.Ascending},
				{"creation_time", commonModel.Ascending},
				{"name_tag", commonModel.Ascending},
				{"group_tag", commonModel.Descending},
			},
		}

		testedResult, newPaging := ListTargets(
			testCase.query, paging,
		)

		c.Logf("[List] Query condition: %v. Number of targets: %d", testCase.query, len(testedResult))

		for _, target := range testedResult {
			c.Logf("[List] Target: %v.", target)
		}
		c.Assert(testedResult, HasLen, testCase.expectedCountOfCurrentPage, Commentf("Test Case: %d", i+1))
		c.Assert(newPaging.TotalCount, Equals, testCase.expectedCountOfAll, Commentf("Test Case: %d", i+1))
	}
}

// Tests the getting of a target by id
func (suite *TestTargetSuite) TestGetSimpleTarget1ById(c *C) {
	testCases := []*struct {
		sampleId int32
		hasFound bool
	}{
		{6981, true},
		{6982, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		c.Assert(GetSimpleTarget1ById(testCase.sampleId), ocheck.ViableValue, testCase.hasFound, comment)
	}
}

// Tests the loading of targets by filter
func (suite *TestTargetSuite) TestLoadSimpleTarget1sByFilter(c *C) {
	testCases := []*struct {
		sampleFilter   *nqmModel.TargetFilter
		expectedNumber int
	}{
		{ // List all of data
			&nqmModel.TargetFilter{}, 3,
		},
		{ // Matches some of data
			&nqmModel.TargetFilter{
				Name: []string{"ftg-1", "ftg-1-C01"},
				Host: []string{"20.45"},
			}, 2,
		},
		{ // Matches nothing
			&nqmModel.TargetFilter{
				Name: []string{"no-such-tt"},
			}, 0,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := LoadSimpleTarget1sByFilter(testCase.sampleFilter)
		c.Assert(testedResult, HasLen, testCase.expectedNumber, comment)
	}
}

func (s *TestTargetSuite) SetUpTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestTargetSuite.TestGetSimpleTarget1ById":
		inTx(
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host)
			VALUES (6981, 'get-by-id-1', 'get-by-id-1')
			`,
		)
	case "TestTargetSuite.TestLoadSimpleTarget1sByFilter":
		inTx(
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host)
			VALUES (15071, 'ftg-1-C01', '20.45.71.91'),
				(15072, 'ftg-1-C02', '20.45.71.92'),
				(15073, 'ftg-2-C01', '120.33.27.23')
			`,
		)
	case "TestTargetSuite.TestUpdateTarget":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3341, 'tg-nt-1')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(10501, 'tg-gt-1'), (10502, 'tg-gt-2')
			`,
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host, tg_available, tg_status, tg_nt_id)
			VALUES (34091, 'org-tg-1', 'up-tg-host-1', false, true, 3341)
			`,
			`
			INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
			VALUES(34091, 10501), (34091, 10502)
			`,
		)
	case "TestTargetSuite.TestListTargets":
		inTx(
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3081, '美國 IP 群')
			`,
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(23401, '串流一網'),(23402, '串流二 網')
			`,
			`
			INSERT INTO nqm_target(tg_id, tg_name, tg_host, tg_available, tg_status, tg_isp_id, tg_nt_id)
			VALUES
				(40201, 'tg-name-1', 'tg-host-1', false, true, 3, -1),
				(40202, 'tg-name-2', 'tg-host-2', true, true, 4, 3081),
				(40203, 'tg-name-3', 'tg-host-3', true, false, 5, -1)
			`,
			`
			INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
			VALUES(40201, 23401), (40201, 23402)
			`,
		)
	case "TestTargetSuite.TestGetTargetById":
		inTx(
			`
			INSERT INTO nqm_target(
				tg_id, tg_name, tg_host, tg_available, tg_status, tg_comment,
				tg_probed_by_all
			)
			VALUES(23041, 'tg-get-1', '123.45.18.19', true, true, 'sure', true)
			`,
		)
	}
}
func (s *TestTargetSuite) TearDownTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestTargetSuite.TestGetSimpleTarget1ById":
		inTx(
			"DELETE FROM nqm_target WHERE tg_id = 6981",
		)
	case "TestTargetSuite.TestLoadSimpleTarget1sByFilter":
		inTx(
			"DELETE FROM nqm_target WHERE tg_id >= 15071 AND tg_id <= 15073",
		)
	case "TestTargetSuite.TestGetTargetById":
		inTx(
			"DELETE FROM nqm_target WHERE tg_id = 23041",
		)
	case "TestTargetSuite.TestAddTarget":
		inTx(
			"DELETE FROM nqm_target WHERE tg_name LIKE 'def-target-%'",
			"DELETE FROM owl_name_tag WHERE nt_value LIKE 'IBM-%'",
			"DELETE FROM owl_group_tag WHERE gt_name LIKE 'KSC-%'",
		)
	case "TestTargetSuite.TestUpdateTarget":
		inTx(
			"DELETE FROM nqm_target WHERE tg_id = 34091",
			"DELETE FROM owl_name_tag WHERE nt_value LIKE 'tg-nt-%'",
			"DELETE FROM owl_group_tag WHERE gt_name LIKE 'tg-gt-%'",
		)
	case "TestTargetSuite.TestListTargets":
		inTx(
			`
			DELETE FROM nqm_target
			WHERE tg_id >= 40201 AND tg_id <= 40203
			`,
			`
			DELETE FROM owl_name_tag WHERE nt_id = 3081
			`,
			`
			DELETE FROM owl_group_tag WHERE gt_id >= 23401 AND gt_id <= 23402
			`,
		)
	}
}

func (s *TestTargetSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
	owlDb.DbFacade = DbFacade
}
func (s *TestTargetSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
	owlDb.DbFacade = nil
}
