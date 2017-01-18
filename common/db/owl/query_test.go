package owl

import (
	"time"
	"encoding/hex"
	"crypto/md5"
	"github.com/Cepave/open-falcon-backend/common/db"
	"github.com/satori/go.uuid"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	owlCheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	t "github.com/Cepave/open-falcon-backend/common/testing"
	. "gopkg.in/check.v1"
)

type TestQuerySuite struct{}

var _ = Suite(&TestQuerySuite{})

// Tests the refreshing of query object
func (suite *TestQuerySuite) TestAddOrRefreshQuery(c *C) {
	query1 := &owlModel.Query{
		NamedId: "test.query.g1",
		Content: []byte { 29, 87, 61, 4, 5, 78, 91, 11, 91 },
		Md5Content: md5.Sum([]byte("This is test - 1")),
	}

	testCases := []*struct {
		sampleQuery *owlModel.Query
		sampleTime string
	} {
		{ // Adds a new onw
			query1, "2015-07-01T07:36:55+08:00",
		},
		{ // Updates an existing one
			query1, "2015-07-01T08:36:55+08:00",
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		sampleQuery := testCase.sampleQuery
		sampleTime := t.ParseTime(c, testCase.sampleTime)
		AddOrRefreshQuery(sampleQuery, sampleTime)

		var testedTime = getAccessTimeByUuid(sampleQuery.Uuid)

		c.Assert(testedTime, owlCheck.TimeEquals, sampleTime, comment)

		c.Assert(sampleQuery.Md5Content, DeepEquals, query1.Md5Content)
		c.Assert(sampleQuery.Content, DeepEquals, query1.Content)
	}
}

// Tests the updateing of access time by UUID(or adding new one)
func (suite *TestQuerySuite) TestUpdateAccessTimeOrAddNewOne(c *C) {
	testCases := []*struct {
		sampleUuid string
		sampleTime string
	} {
		{ "890858a7-d458-435b-b798-1ac0abf1eae2", "2014-03-23T10:32:45Z" },
		{ "e2f23847-48f1-4c9c-8dbd-8f276d2222eb", "2014-04-17T22:40:47Z" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		sampleQuery := &owlModel.Query{
			Uuid: db.DbUuid(t.ParseUuid(c, testCase.sampleUuid)),
			NamedId: "test.query.uu1",
			Content: []byte { },
			Md5Content: md5.Sum([]byte(testCase.sampleUuid)),
		}
		sampleTime := t.ParseTime(c, testCase.sampleTime)

		UpdateAccessTimeOrAddNewOne(sampleQuery, sampleTime)

		testedTime := getAccessTimeByUuid(sampleQuery.Uuid)

		c.Assert(testedTime, owlCheck.TimeEquals, sampleTime, comment)
	}
}

// Tests the load ing query by UUID
func (suite *TestQuerySuite) TestLoadQueryByUuid(c *C) {
	testCases := []*struct {
		sampleUuid string
		expectedMd5Content string
		expectedContent string
	} {
		{
			"209f18f4f89b42568e1e5270987c057d",
			"7011e902d4a848c184e242e8d71aa961",
			"0dfaa9f3df2d4071b2f48b359440e0fc",
		},
		{
			"219f18f4f89b42568e1e5270987c057d",
			"", "",
		},
	}

	sampleTime := t.ParseTime(c, "2013-07-08T10:20:36+08:00")
	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		uuidValue := t.ParseUuid(c, testCase.sampleUuid)
		testedQuery := LoadQueryByUuidAndUpdateAccessTime(
			"test.load.uu2", uuidValue, sampleTime,
		)

		if testCase.expectedMd5Content == "" {
			c.Assert(testedQuery, IsNil, comment)
			continue
		}


		c.Assert(testedQuery.NamedId, Equals, "test.load.uu2", comment)
		c.Assert(uuid.UUID(testedQuery.Uuid), DeepEquals, uuidValue, comment)

		hexTestedContent := hex.EncodeToString(testedQuery.Content)
		c.Assert(hexTestedContent, Equals, testCase.expectedContent, comment)

		hexTestedMd5Content := hex.EncodeToString(testedQuery.Md5Content[:])
		c.Assert(hexTestedMd5Content, Equals, testCase.expectedMd5Content, comment)

		testedTime := getAccessTimeByUuid(testedQuery.Uuid)
		c.Assert(testedTime, owlCheck.TimeEquals, sampleTime, comment)
	}
}

func getAccessTimeByUuid(uuid db.DbUuid) time.Time {
	var timeValue = time.Time{}
	DbFacade.SqlxDbCtrl.QueryRowxAndScan(
		`
		SELECT qr_time_access
		FROM owl_query
		WHERE qr_uuid = ?
		`,
		[]interface{} { uuid },
		&timeValue,
	)

	return timeValue
}

func (s *TestQuerySuite) SetUpTest(c *C) {
	inTx := DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestQuerySuite.TestUpdateAccessTimeOrAddNewOne":
		inTx(
			`
			INSERT INTO owl_query(
				qr_uuid, qr_named_id, qr_md5_content,
				qr_content, qr_time_access, qr_time_creation
			)
			VALUES(
				x'890858a7d458435bb7981ac0abf1eae2', 'test.query.uu1', x'890858a7d458435bb7981ac0abf1eae2',
				x'e2f2384748f14c9c8dbd8f276d2222eb', '2012-05-06T20:14:43', '2012-05-05T08:23:03'
			)
			`,
		)
	case "TestQuerySuite.TestLoadQueryByUuid":
		inTx(
			`
			INSERT INTO owl_query(
				qr_uuid, qr_named_id, qr_md5_content,
				qr_content, qr_time_access, qr_time_creation
			)
			VALUES(
				x'209f18f4f89b42568e1e5270987c057d', 'test.load.uu2', x'7011e902d4a848c184e242e8d71aa961',
				x'0dfaa9f3df2d4071b2f48b359440e0fc', '2012-05-06T20:14:43', '2012-05-05T08:23:03'
			)
			`,
		)
	}
}
func (s *TestQuerySuite) TearDownTest(c *C) {
	inTx := DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestQuerySuite.TestAddOrRefreshQuery":
		inTx(
			`
			DELETE FROM owl_query
			WHERE qr_named_id = 'test.query.g1'
			`,
		)
	case "TestQuerySuite.TestUpdateAccessTimeOrAddNewOne":
		inTx(
			`
			DELETE FROM owl_query
			WHERE qr_named_id = 'test.query.uu1'
			`,
		)
	case "TestQuerySuite.TestLoadQueryByUuid":
		inTx(
			`
			DELETE FROM owl_query
			WHERE qr_named_id = 'test.load.uu2'
			`,
		)
	}
}

func (s *TestQuerySuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestQuerySuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
