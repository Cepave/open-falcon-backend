package nqm

import (
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
	"reflect"
)

type TestTargetSuite struct{}

var _ = Suite(&TestTargetSuite{})

// Tests the adding of data for target
func (suite *TestTargetSuite) TestAddTarget(c *C) {
	/**
	 * sample targets
	 */
	defaultTarget_1 := nqmModel.NewTargetForAdding()
	defaultTarget_1.Name = "def-target-1"
	defaultTarget_1.Host = "13.41.60.178"

	defaultTarget_2 := nqmModel.NewTargetForAdding()
	defaultTarget_2.Name = "def-target-2"
	defaultTarget_2.Host = "13.41.60.179"
	defaultTarget_2.IspId = 3
	defaultTarget_2.ProvinceId = 20
	defaultTarget_2.CityId = 6
	defaultTarget_2.NameTagValue = "IBM-617"
	defaultTarget_2.GroupTags = []string {
		"KSC-03", "KSC-04", "KSC-05",
	}

	defaultTarget_3 := *defaultTarget_2
	defaultTarget_3.Name = "def-target-3"
	defaultTarget_3.Host = "13.41.60.180"
	defaultTarget_3.CityId = 50
	// :~)

	testCases := []*struct {
		addedTarget *nqmModel.TargetForAdding
		hasError bool
		errorType reflect.Type
	} {
		{ defaultTarget_1, false, nil }, // Use the minimum value
		{ defaultTarget_1, true, reflect.TypeOf(ErrDuplicatedNqmTarget{}) },
		{ defaultTarget_2, false, nil }, // Use every properties
		{ &defaultTarget_3, true, reflect.TypeOf(owlDb.ErrNotInSameHierarchy{}) }, // Duplicated connection id
	}

	for _, testCase := range testCases {
		currentAddedTarget := testCase.addedTarget
		newTarget, err := AddTarget(currentAddedTarget)

		/**
		 * Asserts the occuring error
		 */
		if testCase.hasError {
			c.Assert(newTarget, IsNil)
			c.Assert(err, NotNil)
			c.Logf("Has error: %v", err)

			c.Assert(reflect.TypeOf(err), Equals, testCase.errorType)
			continue
		}
		// :~)

		c.Assert(err, IsNil)
		c.Logf("New Target: %v", newTarget)
		c.Logf("New Target[Group Tags]: %v", newTarget.GroupTags)

		c.Assert(newTarget.Name, Equals, currentAddedTarget.Name)
		c.Assert(newTarget.Host, Equals, currentAddedTarget.Host)
		c.Assert(newTarget.IspId, Equals, currentAddedTarget.IspId)
		c.Assert(newTarget.ProvinceId, Equals, currentAddedTarget.ProvinceId)
		c.Assert(newTarget.CityId, Equals, currentAddedTarget.CityId)
		c.Assert(newTarget.NameTagId, Equals, currentAddedTarget.NameTagId)
		c.Assert(newTarget.GroupTags, HasLen, len(currentAddedTarget.GroupTags))
	}
}

// Tests the updating for data of target
func (suite *TestTargetSuite) TestUpdateTarget(c *C) {
	modifiedTarget := &nqmModel.TargetForAdding {
		Name: "new-tg-1",
		Comment: "new-comment-1",
		Status: false,
		ProvinceId: 26,
		CityId: 194,
		IspId: 2,
		NameTagValue: "tg-nt-2",
		GroupTags: []string{"tg-gt-3", "tg-gt-4", "tg-gt-5"},
	}

	originalTarget := GetTargetById(34091)

	testedTarget, err := UpdateTarget(originalTarget, modifiedTarget)

	c.Assert(err, IsNil)
	c.Assert(testedTarget.Name, Equals, modifiedTarget.Name)
	c.Assert(testedTarget.Comment, Equals, modifiedTarget.Comment)
	c.Assert(testedTarget.Status, Equals, modifiedTarget.Status)
	c.Assert(testedTarget.ProvinceId, Equals, modifiedTarget.ProvinceId)
	c.Assert(testedTarget.CityId, Equals, modifiedTarget.CityId)
	c.Assert(testedTarget.IspId, Equals, modifiedTarget.IspId)
	c.Assert(testedTarget.NameTagId, Equals, modifiedTarget.NameTagId)

	testedTargetForAdding := testedTarget.ToTargetForAdding()
	c.Assert(testedTargetForAdding.AreGroupTagsSame(modifiedTarget), Equals, true)
}

// Tests the retrieving of data for a target by id
func (suite *TestTargetSuite) TestGetTargetById(c *C) {
	testCases := []*struct {
		sampleIdOfTarget int32
		hasFound bool
	} {
		{ 23041, true },
		{ 23051, false },
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
		query *nqmModel.TargetQuery
		pageSize int32
		pagePosition int32
		expectedCountOfCurrentPage int
		expectedCountOfAll int32
	} {
		{ // All data
			&nqmModel.TargetQuery {},
			10, 1, 3, 3,
		},
		{ // 2nd page
			&nqmModel.TargetQuery {},
			2, 2, 1, 3,
		},
		{ // Match nothing for futher page
			&nqmModel.TargetQuery {},
			10, 10, 0, 3,
		},
		{ // Match 1 row by all of the conditions
			&nqmModel.TargetQuery {
				Name: "tg-name-1",
				Host: "tg-host-1",
				HasIspId: true,
				IspId: 3,
				HasStatusCondition: true,
				Status: true,
			}, 10, 1, 1, 1,
		},
		{ // Match 1 row(by ISP id)
			&nqmModel.TargetQuery {
				HasIspId: true,
				IspId: 5,
			}, 10, 1, 1, 1,
		},
		{ // Match nothing
			&nqmModel.TargetQuery {
				Name: "tg-easy-1",
				Host: "tg-host-1",
			}, 10, 1, 0, 0,
		},
	}

	for i, testCase := range testCases {
		paging := commonModel.Paging{
			Size: testCase.pageSize,
			Position: testCase.pagePosition,
			OrderBy: []*commonModel.OrderByEntity {
				&commonModel.OrderByEntity{ "name", commonModel.Ascending },
				&commonModel.OrderByEntity{ "status", commonModel.Ascending },
				&commonModel.OrderByEntity{ "host", commonModel.Ascending },
				&commonModel.OrderByEntity{ "comment", commonModel.Ascending },
				&commonModel.OrderByEntity{ "isp", commonModel.Ascending },
				&commonModel.OrderByEntity{ "province", commonModel.Ascending },
				&commonModel.OrderByEntity{ "city", commonModel.Ascending },
				&commonModel.OrderByEntity{ "creation_time", commonModel.Ascending },
				&commonModel.OrderByEntity{ "name_tag", commonModel.Ascending },
				&commonModel.OrderByEntity{ "group_tag", commonModel.Descending },
			},
		}

		testedResult, newPaging := ListTargets(
			testCase.query, paging,
		)

		c.Logf("[List] Query condition: %v. Number of targets: %d", testCase.query, len(testedResult))

		for _, target := range testedResult {
			c.Logf("[List] Target: %v.", target)
		}
		c.Assert(testedResult, HasLen, testCase.expectedCountOfCurrentPage, Commentf("Test Case: %d", i + 1))
		c.Assert(newPaging.TotalCount, Equals, testCase.expectedCountOfAll, Commentf("Test Case: %d", i + 1))
	}
}

// Tests the getting of a target by id
func (suite *TestTargetSuite) TestGetSimpleTarget1ById(c *C) {
	testCases := []*struct {
		sampleId int32
		hasFound bool
	} {
		{ 6981, true },
		{ 6982, false },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		c.Assert(GetSimpleTarget1ById(testCase.sampleId), ocheck.ViableValue, testCase.hasFound, comment)
	}
}

// Tests the loading of targets by filter
func (suite *TestTargetSuite) TestLoadSimpleTarget1sByFilter(c *C) {
	testCases := []*struct {
		sampleFilter *nqmModel.TargetFilter
		expectedNumber int
	} {
		{ // List all of data
			&nqmModel.TargetFilter {}, 3,
		},
		{ // Matches some of data
			&nqmModel.TargetFilter {
				Name: []string{ "ftg-1", "ftg-1-C01" },
				Host: []string{ "20.45" },
			}, 2,
		},
		{ // Matches nothing
			&nqmModel.TargetFilter {
				Name: []string{ "no-such-tt" },
			}, 0,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

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
