package changelog

import (
	dbsql "database/sql"
	psql "github.com/Cepave/open-falcon-backend/scripts/mysql/dbpatch/go/sql"
	"fmt"
	"log"
)

// The constants of result of applying patch
const (
	failed = 0
	applying = 1 // The initial status of change log
	success = 2
)

/**
 * The full configuration for running patches
 */
type ChangeLogConfig struct {
	DriverName string
	Dsn string
	ChangeLog string
	PatchFileBase string
	Delimiter string
}

// Represents the result of log for patching
type patchResult struct {
	patchConfig *PatchConfig
	id int
	result int
	message string
}

// The string representation of ChangeLogConfig
func (changeLogConfig *ChangeLogConfig) String() string {
	return fmt.Sprintf(
		"Database: [%v], [%v]. Change log file: [%v]. Patch Base: [%v]. Delimiter: [%v]",
		changeLogConfig.DriverName, changeLogConfig.Dsn,
		changeLogConfig.ChangeLog, changeLogConfig.PatchFileBase,
		changeLogConfig.Delimiter,
	)
}

// Executes the patches in change log
func ExecutePatches(changeLogConfig *ChangeLogConfig) (err error) {
	/**
	 * Loads configuration of patches
	 */
	var loadedPatches []PatchConfig
	if loadedPatches, err = LoadChangeLogFromFile(changeLogConfig.ChangeLog)
		err != nil {
		return
	}
	// :~)

	/**
	 * Connect to database
	 */
	var dbConfig *psql.DatabaseConfig
	if dbConfig, err = psql.NewDatabaseConfig(
		changeLogConfig.DriverName,
		changeLogConfig.Dsn,
	)
		err != nil {
		return
	}

	defer dbConfig.Close()
	// :~)

	/**
	 * Checking the schema for change log
	 */
	if err = dbConfig.Execute(checkChangeLogSchema); err != nil {
		return
	}
	if err = fixCharset(dbConfig); err != nil {
		return
	}
	// :~)

	/**
	 * Iterates each patch and applies it to database
	 */
	var numberOfApplied = 0
	for _, p := range loadedPatches {
		/**
		 * Checks if the patch has been applied
		 */
		if patchApplied, _ := hasPatchApplied(dbConfig, &p)
			patchApplied {
			continue
		}
		// :~)

		log.Printf("Applying patch: [%v](%v)...", p.Id, p.Filename);

		var scripts []string

		/**
		 * Loads scripts from file
		 */
		if scripts, err = p.loadScripts(changeLogConfig.PatchFileBase, changeLogConfig.Delimiter)
			err != nil {
			return fmt.Errorf("Load script file[%v/%v] error: %v", changeLogConfig.PatchFileBase, p.Filename, err)
		}
		// :~)

		/**
		 * Applies patch to database
		 */
		if err = applyPatch(dbConfig, &p, scripts)
			err != nil {

			var patchErr = fmt.Errorf("Patch [%v](%v) has error: %v", p.Id, p.Filename, err)
			log.Println(patchErr)

			return patchErr;
		}
		// :~)

		numberOfApplied++
		log.Printf("Applying patch success. [%v](%v).", p.Id, p.Filename)
	}
	// :~)

	log.Printf("Number of applied patches: %v", numberOfApplied);
	return
}

// 1. Add log to database
// 2. Applies patch to database
// 3. Update result patching on log in database
func applyPatch(
	dbConfig *psql.DatabaseConfig,
	patchConfig *PatchConfig,
	scripts []string,
) (err error) {
	var patchContent = patchResult {
		patchConfig: patchConfig,
	}

	/**
	 * Add log to dtaabase
	 */
	if err = dbConfig.Execute(newChangeLogFunc(&patchContent))
		err != nil {
		return
	}
	// :~)

	/**
	 * Applies scripts to database
	 */
	for _, script := range scripts {
		log.Printf("Applying script:\n%v\n", script)

		if err = dbConfig.Execute(
			func(db *dbsql.DB) (err error) {
				_, err = db.Exec(script)
				return
			},
		)
			err != nil {

			/**
			 * Logs the failed result to database
			 */
			patchContent.result = failed
			patchContent.message = fmt.Sprintf("Error: [%v]\nScript:\n%v\n", err, script)
			if logErr := dbConfig.Execute(updateChangeLogFunc(&patchContent))
				logErr != nil {
				panic(fmt.Errorf("Cannot log failed patching: [%v]. Error: %f.\nScript:\n%v\n", patchConfig.Id, err, script))
			}
			// :~)

			return
		}
	}
	// :~)

	/**
	 * Update result patching on log in database
	 */
	patchContent.result = success
	if err = dbConfig.Execute(updateChangeLogFunc(&patchContent))
		err != nil {
		return
	}
	// :~)

	return
}

// Checks whether or not a patch has been applied on database
func hasPatchApplied(dbConfig *psql.DatabaseConfig, patchConfig *PatchConfig) (result bool, err error) {
	result = false

	err = dbConfig.Execute(
		func(db *dbsql.DB) (err error) {
			var successRow = -1

			/**
			 * Since MySQL would gives you total number of COUNT(*)
			 * even if you use limit express,
			 * we just check the counting of success patching >= 1
			 */
			if err = db.QueryRow(
				`
				SELECT COUNT(*)
				FROM sysdb_change_log
				WHERE dcl_named_id = ?
					AND dcl_result = 2
				ORDER BY dcl_time_update DESC
				`,
				patchConfig.Id,
			).Scan(&successRow)
				err != nil {
				return
			}

			result = successRow >= 1
			return
		},
	)

	return
}

// Builds the function for adding a new change log to database
func newChangeLogFunc(patchContent *patchResult) func(*dbsql.DB) error {
	return func(db *dbsql.DB) (err error) {
		var result dbsql.Result

		if result, err = db.Exec(
			`
			INSERT INTO sysdb_change_log(dcl_named_id, dcl_comment, dcl_file_name)
			VALUES(?, ?, ?)
			`,
			patchContent.patchConfig.Id,
			patchContent.patchConfig.Comment,
			patchContent.patchConfig.Filename,
		)
			err != nil {
			return
		}

		newIdAsInt64, _ := result.LastInsertId()
		patchContent.id = int(newIdAsInt64)
		patchContent.result = applying

		return
	}
}

//  Builds the function for updating change log to database
func updateChangeLogFunc(patchContent *patchResult) func(*dbsql.DB) error {
	return func(db *dbsql.DB) (err error) {
		_, err = db.Exec(
			`
			UPDATE sysdb_change_log
			SET	dcl_result = ?,
				dcl_message = ?,
				dcl_time_update = NOW()
			WHERE dcl_id = ?
			`,
			patchContent.result,
			patchContent.message,
			patchContent.id,
		)

		return
	}
}

// Checks the schema for change log, creates schema if it is not existing
func checkChangeLogSchema(db *dbsql.DB) (err error) {
	var result int

	if err = db.
		QueryRow(
			`
			SELECT COUNT(*)
			FROM INFORMATION_SCHEMA.TABLES
			WHERE TABLE_SCHEMA = DATABASE()
				AND TABLE_NAME = ?
			`,
			"sysdb_change_log",
		).
			Scan(&result)
		err != nil {
		return
	}

	/**
	 * The schema is existing
	 */
	if result == 1 {
		return
	}
	// :~)

	log.Printf("Initialize 'sysdb_change_log' for patching history")

	/**
	 * Creates the schema for change log
	 */
	if _, err = db.Exec(
		`
		CREATE TABLE sysdb_change_log(
			dcl_id INT AUTO_INCREMENT PRIMARY KEY,
			dcl_named_id VARCHAR(128) NOT NULL,
			dcl_file_name VARCHAR(512) NOT NULL,
			dcl_result TINYINT NOT NULL DEFAULT 1,
			dcl_time_creation DATETIME NOT NULL DEFAULT NOW(),
			dcl_time_update DATETIME,
			dcl_message VARCHAR(1024),
			dcl_comment VARCHAR(1024)
		) DEFAULT CHARSET=UTF8
		`,
	)
		err != nil {
		return
	}

	_, err = db.Exec(
		`
		CREATE UNIQUE INDEX ix_sysdb_change_log__result
		ON sysdb_change_log(dcl_named_id, dcl_result DESC, dcl_time_update DESC)
		`,
	)
	// :~)

	return
}

// Fix the charset of sysdb_change_log table
func fixCharset(dbConfig *psql.DatabaseConfig) (err error) {
	defer func() {
		p := recover()
		if p != nil {
			switch p.(type) {
			case error:
				err = p.(error)
			default:
				err = fmt.Errorf("Fix charset(sysdb_change_log) error: %v", p)
			}
		}
	}()

	dbName := dbConfig.GetDatabaseName()

	var nonUtf8Count uint
	dbConfig.SqlxDbCtrl.QueryRowxAndScan(
		`
		SELECT COUNT(column_name)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = ?
			AND TABLE_NAME='sysdb_change_log'
			AND CHARACTER_SET_NAME != 'utf8'
			AND (COLUMN_TYPE LIKE 'CHAR%' OR COLUMN_TYPE LIKE 'VARCHAR%')
		`,
		[]interface{} { dbName },
		&nonUtf8Count,
	)

	if nonUtf8Count == 0 {
		return
	}

	log.Printf("Modify charset of `sysdb_change_log` for database: [%s]", dbName)

	dbConfig.SqlxDb.MustExec(
		`
		ALTER TABLE sysdb_change_log
		CONVERT TO CHARACTER SET UTF8
		`,
	)

	return
}
