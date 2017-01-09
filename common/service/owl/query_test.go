package owl

import (
	"encoding/hex"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	db "github.com/Cepave/open-falcon-backend/common/db"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	t "github.com/Cepave/open-falcon-backend/common/testing"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
	"time"
)

type TestQuerySuite struct{}

var _ = Suite(&TestQuerySuite{})

var testedQueryService = NewQueryService(
	QueryServiceConfig{
		"qservice-q-1", 4, 5 * time.Second,
	},
)

// Tests the loading of query by uuid
func (suite *TestQuerySuite) TestLoadQueryByUuid(c *C) {
	testCases := []*struct {
		sampleUuid string
		expectedMd5Content string
		inCache bool
	} {
		{ "890818a7d438495bb79a1ac0abf1eae2", "8c0a58a7d458430bb79812c0abf1eae7", false }, // Load data into cache
		{ "120818a7d438495bb79a1ac0abf1eae2", "", false }, // Nothing found
		{ "120818a7d438495bb79a1ac0abf1eae2", "", false }, // Nothing found
		{ "890818a7d438495bb79a1ac0abf1eae2", "8c0a58a7d458430bb79812c0abf1eae7", true }, // Update access time
	}

	var lastAccessTime time.Time = time.Now()
	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		sampleUuid := t.ParseUuid(c, testCase.sampleUuid)

		/**
		 * Asserts the existence itme in cache
		 */
		if testCase.inCache {
			c.Assert(testedQueryService.cache.Get(sampleUuid.String()), NotNil, comment)
		} else {
			c.Assert(testedQueryService.cache.Get(sampleUuid.String()), IsNil, comment)
		}
		// :~)

		time.Sleep(1 * time.Second)

		testedQuery := testedQueryService.LoadQueryByUuid(sampleUuid)

		if testCase.expectedMd5Content == "" {
			c.Assert(testedQuery, IsNil)
			continue
		}

		/**
		 * Asserts the loaded data
		 */
		hexTestedMd5Content := hex.EncodeToString(testedQuery.Md5Content[:])
		c.Assert(hexTestedMd5Content, Equals, testCase.expectedMd5Content, comment)
		// :~)

		/**
		 * Asserts the update of access time
		 */
		accessTime := getAccessTimeByUuid(testedQuery.Uuid)
		c.Logf("Access time: %v", accessTime)
		c.Assert(accessTime, ocheck.TimeAfter, lastAccessTime, comment)
		lastAccessTime = accessTime
		// :~)
	}
}

// Tests the creation of loading of query
func (suite *TestQuerySuite) TestCreateOrLoadQuery(c *C) {
	var md5Value db.Bytes16
	(&md5Value).Scan("810512c76a1c44ddb0d6097ef4ef156e")

	sampleQuery := &owlModel.Query{
		NamedId: "gp-query-1",
		Md5Content: md5Value,
		Content: md5Value[:],
	}

	testedQueryService.CreateOrLoadQuery(sampleQuery)

	c.Assert(sampleQuery.Uuid.IsNil(), Equals, false, Commentf("UUID is nil value"))
	// Asserts the value is in cache
	c.Assert(
		testedQueryService.cache.Get(KeyByDbUuid(sampleQuery.Uuid)),
		NotNil,
	)

	/**
	 * Trying to loads the same query with the same md5 value
	 */
	sampleQuery_2 := &owlModel.Query{
		NamedId: "gp-query-1",
		Md5Content: md5Value,
		Content: md5Value[:],
	}
	testedQueryService.CreateOrLoadQuery(sampleQuery_2)
	c.Assert(sampleQuery_2.Uuid, DeepEquals, sampleQuery.Uuid)
	// :~)
}

func getAccessTimeByUuid(uuid db.DbUuid) time.Time {
	var timeValue = time.Time{}
	owlDb.DbFacade.SqlxDbCtrl.QueryRowxAndScan(
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
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestQuerySuite.TestLoadQueryByUuid":
		inTx(
			`
			INSERT INTO owl_query(
				qr_uuid, qr_named_id, qr_md5_content,
				qr_content, qr_time_access, qr_time_creation
			)
			VALUES(
				x'890818a7d438495bb79a1ac0abf1eae2', 'qservice-q-1', x'8c0a58a7d458430bb79812c0abf1eae7',
				x'e2f2384948f94c9c8dbd87278d2222eb', '2012-05-06T20:14:43', '2012-05-05T08:23:03'
			)
			`,
		)
	}
}

func (s *TestQuerySuite) TearDownTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestQuerySuite.TestLoadQueryByUuid":
		inTx("DELETE FROM owl_query WHERE qr_named_id = 'qservice-q-1'")
	case "TestQuerySuite.TestCreateOrLoadQuery":
		inTx("DELETE FROM owl_query WHERE qr_named_id = 'gp-query-1'")
	}
}

func (s *TestQuerySuite) SetUpSuite(c *C) {
	owlDb.DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestQuerySuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, owlDb.DbFacade)
}
