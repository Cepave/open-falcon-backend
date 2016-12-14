package owl

import (
	"crypto/md5"
	"github.com/Cepave/open-falcon-backend/common/db"
	"github.com/satori/go.uuid"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	owlCheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
	"time"
)

type TestQuerySuite struct{}

var _ = Suite(&TestQuerySuite{})

// Tests the refreshing of query object
func (suite *TestQuerySuite) TestAddOrRefreshQuery(c *C) {
	testCases := []struct {
		sampleQuery *owlModel.Query
		sampleTime string
	} {
		{ // Adds a new onw
			&owlModel.Query{
				NamedId: "test.query.g1",
				Content: []byte { 29, 87, 61, 4, 5, 78, 91, 11, 91 },
				Md5Content: md5.Sum([]byte("This is test - 1")),
			},
			"2015-07-01T07:36:55+08:00",
		},
		{ // Updates an existing one
			&owlModel.Query{
				NamedId: "test.query.g1",
				Content: []byte { 29, 87, 61, 4, 5, 78, 91, 11, 91 },
				Md5Content: md5.Sum([]byte("This is test - 1")),
			},
			"2015-07-01T08:36:55+08:00",
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		sampleQuery := testCase.sampleQuery
		sampleTime := parseTime(c, testCase.sampleTime)
		AddOrRefreshQuery(sampleQuery, sampleTime)

		var testedTime = getAccessTimeByUuid(sampleQuery.Uuid)
		c.Assert(testedTime, owlCheck.TimeEquals, sampleTime, comment)
	}
}

// Tests the updateing of access time by UUID(or adding new one)
func (suite *TestQuerySuite) TestUpdateAccessTimeOrAddNewOne(c *C) {
	testCases := []struct {
		sampleUuid string
		sampleTime string
	} {
		{ "890858a7-d458-435b-b798-1ac0abf1eae2", "2014-03-23T10:32:45Z" },
		{ "e2f23847-48f1-4c9c-8dbd-8f276d2222eb", "2014-04-17T22:40:47Z" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		sampleQuery := &owlModel.Query{
			Uuid: db.DbUuid(parseUuid(c, testCase.sampleUuid)),
			NamedId: "test.query.uu1",
			Content: []byte { },
			Md5Content: md5.Sum([]byte(testCase.sampleUuid)),
		}
		sampleTime := parseTime(c, testCase.sampleTime)

		UpdateAccessTimeOrAddNewOne(sampleQuery, sampleTime)

		testedTime := getAccessTimeByUuid(sampleQuery.Uuid)

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

func parseTime(c *C, timeAsString string) time.Time {
	timeValue, err := time.Parse(time.RFC3339, timeAsString)
	c.Assert(err, IsNil)

	return timeValue
}

func parseUuid(c *C, uuidString string) uuid.UUID {
	uuidValue, err := uuid.FromString(uuidString)
	c.Assert(err, IsNil)

	return uuidValue
}

func (s *TestQuerySuite) SetUpTest(c *C) {
	switch c.TestName() {
	case "TestQuerySuite.TestUpdateAccessTimeOrAddNewOne":
		prepareUpdateAccessTimeOrAddNewOne()
	}
}
func (s *TestQuerySuite) TearDownTest(c *C) {
	switch c.TestName() {
	case "TestQuerySuite.TestAddOrRefreshQuery":
		cleanAddOrRefreshQuery()
	case "TestQuerySuite.TestUpdateAccessTimeOrAddNewOne":
		cleanUpdateAccessTimeOrAddNewOne()
	}
}

func prepareUpdateAccessTimeOrAddNewOne() {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(
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
}
func cleanUpdateAccessTimeOrAddNewOne() {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(
		`
		DELETE FROM owl_query
		WHERE qr_named_id = 'test.query.uu1'
		`,
	)
}

func cleanAddOrRefreshQuery() {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(
		`
		DELETE FROM owl_query
		WHERE qr_named_id = 'test.query.g1'
		`,
	)
}

func (s *TestQuerySuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestQuerySuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
