package changelog

import (
	dbsql "database/sql"
	. "gopkg.in/check.v1"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * Tests the creation of database schema
 */
func (suite *ChangeLogSuite) TestCheckChangeLogSchema(c *C) {
	var err error

	/**
	 * First time creation for schema of change log
	 */
	err = dbConfig.Execute(checkChangeLogSchema)
	c.Assert(err, IsNil)
	assertChangeLogSchema(c)
	// :~)

	/**
	 * Tests the existing schema of change log
	 */
	err = dbConfig.Execute(checkChangeLogSchema)
	c.Assert(err, IsNil)
	assertChangeLogSchema(c)
	// :~)
}

func assertChangeLogSchema(c *C) {
	var result int = -1
	err := dbConfig.Execute(func (db *dbsql.DB) (err error) {
		err = db.
			QueryRow(`
			SELECT COUNT(*)
			FROM INFORMATION_SCHEMA.TABLES
			WHERE TABLE_NAME = ?
				AND TABLE_SCHEMA = DATABASE()
			`,
			"sysdb_change_log",
		).
			Scan(&result)
		return
	})

	c.Assert(err, IsNil)
	c.Assert(result, Equals, 1)
}

/**
 * Tests the applying for patches
 */
func (suite *ChangeLogSuite) TestApplyPatch(c *C) {
	var scripts = []string{
		`
		CREATE TABLE sample_patch_1(
			id INT
		)
		`,
		`
		INSERT INTO sample_patch_1 VALUES(1)
		`,
	}

	err := applyPatch(dbConfig, &PatchConfig{ Id: "apply-ph-1", Filename: "apply-ph-1.sql" }, scripts)
	c.Assert(err, IsNil)

	/**
	 * Asserts the result data
	 */
	var testedResult int = -1
	err = dbConfig.Execute(
		func(db *dbsql.DB) (err error) {
			err = db.QueryRow(
				`
				SELECT COUNT(*)
				FROM sample_patch_1
				`,
			).Scan(&testedResult)
			return
		},
	)
	c.Assert(err, IsNil)
	c.Assert(testedResult, Equals, 1)
	// :~)

	/**
	 * Asserts the data of change log
	 */
	var testedChangeLog int = -1

	err = dbConfig.Execute(func(db *dbsql.DB) (err error) {
		err = db.QueryRow(
			`
			SELECT COUNT(dcl_id)
			FROM sysdb_change_log
			WHERE dcl_named_id = ?
				AND dcl_result = 2
			`,
			"apply-ph-1",
		).
			Scan(&testedChangeLog)
		return
	})

	c.Assert(err, IsNil)
	c.Assert(testedChangeLog, Equals, 1)
	// :~)
}

/**
 * Tests the applying for patches with error, which the log of failed should be updated
 */
func (suite *ChangeLogSuite) TestApplyPatchWithError(c *C) {
	var scripts = []string{
		`
		CREATE NO_SUCH_TABLE
		`,
	}

	/**
	 * Asserts the error
	 */
	err := applyPatch(dbConfig, &PatchConfig{ Id: "err-ph-1", Filename: "apply-ph-1.sql" }, scripts)
	c.Assert(err, NotNil)
	// :~)

	/**
	 * Asserts the data of change log
	 */
	var testedChangeLog int = -1

	err = dbConfig.Execute(func(db *dbsql.DB) (err error) {
		err = db.QueryRow(
			`
			SELECT COUNT(dcl_id)
			FROM sysdb_change_log
			WHERE dcl_named_id = ?
				AND dcl_result = 0
			`,
			"err-ph-1",
		).
			Scan(&testedChangeLog)
		return
	})

	c.Assert(err, IsNil)
	c.Assert(testedChangeLog, Equals, 1)
	// :~)
}

/**
 * Tests the checking for whether or not a patch has been applied on database
 */
func (suite *ChangeLogSuite) TestHasPatchApplied(c *C) {
	var testCases = []struct {
		patchConfig PatchConfig
		expectedResult bool
	} {
		{ PatchConfig { Id: "hp-sample-1" }, true }, // Successful
		{ PatchConfig { Id: "hp-sample-2" }, false }, // Failed
		{ PatchConfig { Id: "hp-sample-3" }, false }, // Not applied
	}

	for _, testCase := range testCases {
		testedResult, err := hasPatchApplied(dbConfig, &testCase.patchConfig)
		c.Assert(err, IsNil)
		c.Assert(testedResult, Equals, testCase.expectedResult)
	}
}

/**
 * Tests the updating of change log
 */
func (suite *ChangeLogSuite) TestUpdateChangeLogFunc(c *C) {
	var samplePatchContent = patchResult{
		id: 501,
		result: 2,
		message: "message-1",
	}

	err := dbConfig.Execute(updateChangeLogFunc(&samplePatchContent))
	c.Assert(err, IsNil)

	var testedResult int = -1
	var testedMessage string = ""
	var testedUpdateTime dbsql.NullString
	err = dbConfig.Execute(
		func(db *dbsql.DB) (err error) {

			err = db.QueryRow(
				`
				SELECT dcl_result, dcl_message, dcl_time_update
				FROM sysdb_change_log
				WHERE dcl_id = 501
				`,
			).Scan(&testedResult, &testedMessage, &testedUpdateTime)

			return
		},
	)

	c.Assert(err, IsNil)
	c.Assert(testedResult, Equals, samplePatchContent.result)
	c.Assert(testedMessage, Equals, samplePatchContent.message)
	c.Assert(testedUpdateTime.Valid, Equals, true)
}

/**
 * Tests the creation of change log
 */
func (suite *ChangeLogSuite) TestNewChangeLogFunc(c *C) {
	var testedLogContent = patchResult{
		patchConfig: &PatchConfig{
			Id: "bob-1",
			Comment: "comment-1",
			Filename: "bob-1.sql",
		},
	}

	err := dbConfig.Execute(newChangeLogFunc(&testedLogContent))
	c.Assert(err, IsNil)

	/**
	 * Asserts the inserted data
	 */
	c.Assert(testedLogContent.id > 0, Equals, true)

	var testedNamedId string = ""
	var testedComment string = ""
	var testedFilename string = ""
	var testedResult int = -1

	err = dbConfig.Execute(func(db *dbsql.DB) (err error) {
		err = db.QueryRow(
			`
			SELECT dcl_named_id, dcl_result, dcl_comment, dcl_file_name
			FROM sysdb_change_log
			WHERE dcl_id = ?
			`,
			testedLogContent.id,
		).
			Scan(&testedNamedId, &testedResult, &testedComment, &testedFilename)
		return
	})

	c.Assert(testedNamedId, Equals, "bob-1")
	c.Assert(testedResult, Equals, 1)
	c.Assert(testedComment, Equals, "comment-1")
	c.Assert(testedFilename, Equals, "bob-1.sql")
	c.Assert(err, IsNil)
	// :~)
}
