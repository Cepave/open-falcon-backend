package sql

import (
	dbsql "database/sql"
	_ "github.com/mattn/go-sqlite3"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type SqlSuite struct{}

var _ = Suite(&SqlSuite{})

/**
 * Tests the execution to database
 */
func (suite *SqlSuite) TestExecute(c *C) {
	var sampleDbConfig, err = NewDatabaseConfig(
		"sqlite3", ":memory:",
	)

	c.Assert(err, IsNil)

	defer sampleDbConfig.Close()

	err = sampleDbConfig.Execute(
		func(db *dbsql.DB) (err error) {
			if _, err = db.Exec("CREATE TABLE test_execute (te_id INT)"); err != nil {
				return
			}

			if _, err = db.Exec("INSERT INTO test_execute VALUES(1)"); err != nil {
				return
			}

			/**
			 * Loads the result data
			 */
			var rows, _ = db.Query("SELECT COUNT(*) FROM test_execute")
			defer rows.Close()
			rows.Next()
			var result int
			rows.Scan(&result)
			// :~)

			c.Assert(result, Equals, 1)

			return
		},
	)

	c.Assert(err, IsNil)
}
