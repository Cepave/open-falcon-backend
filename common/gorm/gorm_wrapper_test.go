package gorm

import (
	"database/sql"
	"github.com/Cepave/open-falcon-backend/common/db"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestGormSuite struct{}

var _ = Suite(&TestGormSuite{})

// Tests the error by panic
func (suite *TestGormSuite) TestPanicIfError(c *C) {
	db := buildGormDb(c)
	defer db.Close()

	dbExt := ToDefaultGormDbExt(db.Exec(
		"INSERT INTO no_such_table VALUES(?, ?)", 21, "Bob",
	))

	c.Assert(
		func() {
			dbExt.PanicIfError()
		},
		PanicMatches,
		"(?s:.*no such table.*)",
	)
}

// Tests the Rows() function
func (suite *TestGormSuite) TestRows(c *C) {
	db := buildGormDb(c)
	defer db.Close()

	db.Exec("CREATE TABLE test_rows(tr_id INT, tr_name VARCHAR(32))")
	db.Exec("INSERT INTO test_rows VALUES(?, ?)", 21, "Bob")
	db.Exec("INSERT INTO test_rows VALUES(?, ?)", 22, "Joe")

	dbQuery := db.Raw("SELECT * FROM test_rows")
	dbExt := ToDefaultGormDbExt(dbQuery)
	rows := dbExt.Rows()

	defer rows.Close()
	var numberOfRows int = 0
	for rows.Next() {
		numberOfRows++
	}

	c.Assert(numberOfRows, Equals, 2)
}

// Tests the ScanRows() function
func (suite *TestGormSuite) TestScanRows(c *C) {
	db := buildGormDb(c)
	defer db.Close()

	/**
	 * Prepares data
	 */
	db.Exec("CREATE TABLE test_scan_rows(sr_id INT, sr_name VARCHAR(32))")
	db.Exec("INSERT INTO test_scan_rows VALUES(?, ?)", 21, "Bob")
	db.Exec("INSERT INTO test_scan_rows VALUES(?, ?)", 22, "Joe")
	// :~)

	dbExt := ToDefaultGormDbExt(db.Raw("SELECT * FROM test_scan_rows"))
	rows := dbExt.Rows()

	type sampleRow struct {
		Id   int    `gorm:"column:sr_id"`
		Name string `gorm:"column:sr_name"`
	}

	defer rows.Close()
	var numberOfRows int = 0
	for rows.Next() {
		rowData := sampleRow{}
		dbExt.ScanRows(rows, &rowData)

		c.Logf("Row Data: %#v", rowData)
		numberOfRows++
	}

	c.Assert(numberOfRows, Equals, 2)
}

// Tests the iteration of rows
func (suite *TestGormSuite) TestIterateRows(c *C) {
	ormdb := buildGormDb(c)
	defer ormdb.Close()

	ormdb.Exec("CREATE TABLE test_iter_rows(sr_id INT, sr_name VARCHAR(32))")
	ormdb.Exec("INSERT INTO test_iter_rows VALUES(?, ?)", 21, "Bob")
	ormdb.Exec("INSERT INTO test_iter_rows VALUES(?, ?)", 22, "Joe")

	dbQuery := ormdb.Raw("SELECT * FROM test_iter_rows")
	dbExt := ToDefaultGormDbExt(dbQuery)

	var numberOfRows int = 0
	dbExt.IterateRows(db.RowsCallbackFunc(func(rows *sql.Rows) db.IterateControl {
		numberOfRows++
		return db.IterateContinue
	}))

	c.Assert(numberOfRows, Equals, 2)
}

// Tests the Gorm Ext for transaction
func (suite *TestGormSuite) TestInTx(c *C) {
	ormdb := buildGormDb(c)
	defer ormdb.Close()

	/**
	 * Prepares data
	 */
	ToDefaultGormDbExt(ormdb).InTx(TxCallbackFunc(func(gormDb *gorm.DB) db.TxFinale {
		ToDefaultGormDbExt(gormDb.Exec("CREATE TABLE seed_761(sd_id INT, sd_name VARCHAR(32))")).PanicIfError()
		ToDefaultGormDbExt(gormDb.Exec("INSERT INTO seed_761 VALUES(?, ?)", 11, "Bob")).PanicIfError()
		ToDefaultGormDbExt(gormDb.Exec("INSERT INTO seed_761 VALUES(?, ?)", 12, "Joe")).PanicIfError()
		ToDefaultGormDbExt(gormDb.Exec("INSERT INTO seed_761 VALUES(?, ?)", 13, "JoZZe")).PanicIfError()

		return db.TxCommit
	}))
	// :~)

	dbQuery := ormdb.Raw("SELECT * FROM seed_761")
	gormExtQuery := ToDefaultGormDbExt(dbQuery)

	var numberOfRows int = 0
	gormExtQuery.IterateRows(db.RowsCallbackFunc(func(rows *sql.Rows) db.IterateControl {
		numberOfRows++
		return db.IterateContinue
	}))

	c.Assert(numberOfRows, Equals, 3)
}

// Tests the checking if the record cannot be found
func (suite *TestGormSuite) TestIfRecordNotFound(c *C) {
	ormdb := buildGormDb(c)
	defer ormdb.Close()

	/**
	 * Prepares data
	 */
	ToDefaultGormDbExt(ormdb.Exec("CREATE TABLE king_761(ks_id INT, ks_name VARCHAR(32))")).
		PanicIfError()
	// :~)

	idValue := struct {
		Id int `gorm:"column:ks_id"`
	}{}
	gormExt := ToDefaultGormDbExt(
		ormdb.Table("king_761").Select("ks_id").Where("ks_name = 'oksdfsf'").
			Scan(&idValue),
	)

	c.Assert(gormExt.IfRecordNotFound(32, 77), Equals, 77)
}

func buildGormDb(c *C) *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")

	c.Assert(err, IsNil)

	return db
}
