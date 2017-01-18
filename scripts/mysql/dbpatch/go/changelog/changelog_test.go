package changelog

import (
	"testing"
	dbsql "database/sql"
	patchsql "github.com/Cepave/open-falcon-backend/scripts/mysql/dbpatch/go/sql"
	"flag"
	. "gopkg.in/check.v1"
)

var (
	/**
	 * Defines the flags for teste depend on real connection of database
	 */
	driverName = flag.String("driverName", "", "Name of driver of database")
	dsn = flag.String("dsn", "", "Data source name for database")
	// :~)
)

func Test(t *testing.T) {
	TestingT(t)
}

type ChangeLogSuite struct{}

var _ = Suite(&ChangeLogSuite{})

const SAMPLE_YAML =
`
- {
	id: "bob-1",
	filename: "bob-1.sql",
	comment: "comment-1"
}
- {
	id: "bob-2",
	filename: "bob-2.sql"
}
`

/**
 * Tests the loading from YAML to struct "PatchConfig"
 */
func (suite *ChangeLogSuite) TestLoadChangeLog(c *C) {
	var testedConfigOfPatches, err = LoadChangeLog([]byte(SAMPLE_YAML))

	c.Assert(err, IsNil)
	c.Assert(len(testedConfigOfPatches), Equals, 2)

	/**
	 * Asserts the loaded properties of patch
	 */
	c.Assert(testedConfigOfPatches[0].Id, Equals, "bob-1")
	c.Assert(testedConfigOfPatches[0].Filename, Equals, "bob-1.sql")
	c.Assert(testedConfigOfPatches[0].Comment, Equals, "comment-1")
	c.Assert(testedConfigOfPatches[1].Id, Equals, "bob-2")
	c.Assert(testedConfigOfPatches[1].Comment, Equals, "")
	// :~)
}

/**
 * Tests the loading of change log from file
 */
func (suite *ChangeLogSuite) TestLoadChangeLogFromFile(c *C) {
	var testedConfigOfPatches, err = LoadChangeLogFromFile("../test/TestLoadChangeLogFromFile.yaml")

	c.Assert(err, IsNil)
	c.Assert(len(testedConfigOfPatches), Equals, 2)
}

/**
 * Tests the loading of scripts
 */
func (suite *ChangeLogSuite) TestLoadScripts(c *C) {
	var samplePatchConfig = PatchConfig{
		Id: "hello-1",
		Filename: "../test/hello-1.sql",
	}

	var testScripts, err = samplePatchConfig.loadScripts(".", ";")
	c.Assert(err, IsNil)
	c.Assert(testScripts, HasLen, 3)
	c.Assert(testScripts[0], Equals, "CREATE TABLE tab_1(id INT)")
	c.Assert(testScripts[1], Equals, "CREATE TABLE tab_2(id INT)")
	c.Assert(testScripts[2], Matches, "CREATE PROCEDURE proc().+UPDATE tab_1 SET id = 20;.+")
}

// Global usage(packge level)
var dbConfig *patchsql.DatabaseConfig = nil

// 1. Setup connection of database
func (s *ChangeLogSuite) SetUpSuite(c *C) {
	flag.Parse()

	if *driverName == "" {
		c.Log("No database assigned, some tests would be skipped. -driverName=<> -dsn=<>")
		return
	}

	c.Logf("Connect to database. Driver Name: [%v]. DSN: [%v]", *driverName, *dsn)

	var err error
	if dbConfig, err = patchsql.NewDatabaseConfig(*driverName, *dsn)
		err != nil {

		c.Fatalf("Connect database error: %v", err)
	}
}

// 1. Release connection of database
func (s *ChangeLogSuite) TearDownSuite(c *C) {
	if dbConfig != nil {
		c.Log("Release database connection")
		dbConfig.Close()
	}
}

func (s *ChangeLogSuite) SetUpTest(c *C) {
	switch c.TestName() {
	/**
	 * Remove the schema of tables created by test
	 */
	case "ChangeLogSuite.TestApplyPatchWithError":
		if dbConfig == nil {
			c.Skip("Without database connection")
			return
		}
		dbConfig.Execute(checkChangeLogSchema)
	case "ChangeLogSuite.TestHasPatchApplied":
		if dbConfig == nil {
			c.Skip("Without database connection")
			return
		}
		dbConfig.Execute(checkChangeLogSchema)
		dbConfig.Execute(
			func(db *dbsql.DB) (err error) {
				_, err = db.Exec(
					`
					INSERT INTO sysdb_change_log(dcl_id, dcl_named_id, dcl_result, dcl_time_creation, dcl_time_update, dcl_file_name)
					VALUES(505, 'hp-sample-1', 2, NOW(), NOW(), 'hp-sample-1.sql')
					`,
				)
				_, err = db.Exec(
					`
					INSERT INTO sysdb_change_log(dcl_id, dcl_named_id, dcl_result, dcl_time_creation, dcl_time_update, dcl_file_name)
					VALUES(506, 'hp-sample-2', 0, NOW(), NOW(), 'hp-sample-2.sql')
					`,
				)
				return
			},
		)
	case "ChangeLogSuite.TestApplyPatch":
		if dbConfig == nil {
			c.Skip("Without database connection")
			return
		}
		dbConfig.Execute(checkChangeLogSchema)
	case "ChangeLogSuite.TestCheckChangeLogSchema":
		if dbConfig == nil {
			c.Skip("Without database connection")
			return
		}
	case "ChangeLogSuite.TestNewChangeLogFunc":
		if dbConfig == nil {
			c.Skip("Without database connection")
			return
		}
		dbConfig.Execute(checkChangeLogSchema)
	case "ChangeLogSuite.TestUpdateChangeLogFunc":
		if dbConfig == nil {
			c.Skip("Without database connection")
			return
		}
		dbConfig.Execute(checkChangeLogSchema)
		if err := dbConfig.Execute(
			/**
			 * Prepares the data for updating of log
			 */
			func(db *dbsql.DB) (err error) {
				_, err = db.Exec(
					`
					INSERT INTO sysdb_change_log(dcl_id, dcl_named_id, dcl_result, dcl_time_creation, dcl_file_name)
					VALUES(501, 'update-sample-1', 1, NOW(), 'us-1.sql')
					`,
				)
				return
			},
			// :~)
		)
			err != nil {
			c.Fatalf("Cannot prepare data: %v", err)
		}
	// :~)
	}
}

func (s *ChangeLogSuite) TearDownTest(c *C) {
	switch c.TestName() {
	/**
	 * Remove the schema of tables created by test
	 */
	case "ChangeLogSuite.TestApplyPatch":
		dbConfig.Execute(func (db *dbsql.DB) (err error) {
			_, err = db.Exec(
				"DROP TABLE IF EXISTS sample_patch_1",
			)
			return
		})
		fallthrough
	case "ChangeLogSuite.TestHasPatchApplied",
		"ChangeLogSuite.TestCheckChangeLogSchema",
		"ChangeLogSuite.TestNewChangeLogFunc",
		"ChangeLogSuite.TestApplyPatchWithError",
		"ChangeLogSuite.TestUpdateChangeLogFunc":
		dbConfig.Execute(func (db *dbsql.DB) (err error) {
			_, err = db.Exec(
				"DROP TABLE IF EXISTS sysdb_change_log",
			)
			return
		})
	// :~)
	}
}
