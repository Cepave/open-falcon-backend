package gorm

import (
	"testing"
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/Cepave/open-falcon-backend/common/db"
	. "gopkg.in/check.v1"
	_ "github.com/mattn/go-sqlite3"
)

func Test(t *testing.T) { TestingT(t) }

type TestGormSuite struct{}

var _ = Suite(&TestGormSuite{})

// Tests the error by panic
func (suite *TestGormSuite) TestPanicIfError(c *C) {
	db := buildGormDb(c)
	defer db.Close()

	dbExt := ToGormDbExt(db.Exec(
		"INSERT INTO no_such_table VALUES(?, ?)", 21, "Bob",
	))

	c.Assert(
		func() {
			dbExt.PanicIfError()
		},
		PanicMatches,
		"no such table.+",
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
	dbExt := ToGormDbExt(dbQuery)
	rows := dbExt.Rows()

	defer rows.Close()
	var numberOfRows int = 0
	for rows.Next() {
		numberOfRows++
	}

	c.Assert(numberOfRows, Equals, 2)
}

type TableScanRows struct {
	Id int `gorm:"column:sr_id"`
	Name string `gorm:"column:sr_name"`
}

// Tests the ScanRows() function
func (suite *TestGormSuite) TestScanRows(c *C) {
	db := buildGormDb(c)
	defer db.Close()

	db.Exec("CREATE TABLE test_scan_rows(sr_id INT, sr_name VARCHAR(32))")
	db.Exec("INSERT INTO test_scan_rows VALUES(?, ?)", 21, "Bob")
	db.Exec("INSERT INTO test_scan_rows VALUES(?, ?)", 22, "Joe")

	dbQuery := db.Raw("SELECT * FROM test_scan_rows")
	dbExt := ToGormDbExt(dbQuery)
	rows := dbExt.Rows()

	defer rows.Close()
	var numberOfRows int = 0
	for rows.Next() {
		rowData := TableScanRows{}
		dbExt.ScanRows(rows, &rowData)

		c.Logf("Row Data: %v", rowData)
		c.Assert(rowData.Id, Equals, 21 + numberOfRows)

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
	dbExt := ToGormDbExt(dbQuery)

	var numberOfRows int = 0
	dbExt.IterateRows(db.RowsCallbackFunc(func (rows *sql.Rows) db.IterateControl {
		numberOfRows++
		return db.IterateContinue
	}))

	c.Assert(numberOfRows, Equals, 2)
}

func buildGormDb(c *C) *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:");

	if err != nil {
		c.Fatalf("Open Databae Error. %v", err)
	}

	return db
}
