package sqlx

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	. "gopkg.in/check.v1"
)

type TestWrapperSuite struct{}

var _ = Suite(&TestWrapperSuite{})

// Tests the query(rowx) and scan(on values) for db controller
func (suite *TestWrapperSuite) TestQueryRowxAndScanOfDbCtrl(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE pip_307(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")
	dbCtrl.SqlxDb().MustExec("INSERT INTO pip_307(pi_id, pi_name) VALUES(61, 'Here-1')")
	dbCtrl.SqlxDb().MustExec("INSERT INTO pip_307(pi_id, pi_name) VALUES(62, 'Here-2')")

	testedData := struct {
		Id   int
		Name string
	}{}
	dbCtrl.QueryRowxAndScan(
		"SELECT * FROM pip_307 WHERE pi_name = 'Here-1'",
		nil, &testedData.Id, &testedData.Name,
	)

	c.Logf("QueryRowxAndScan: %#v", testedData)
	c.Assert(testedData.Id, Equals, 61)
	c.Assert(testedData.Name, Equals, "Here-1")
}

// Tests the query(rowx) and scan(on map) for db controller
func (suite *TestWrapperSuite) TestQueryRowxAndMapScanOfDbCtrl(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE zip_604(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")
	dbCtrl.SqlxDb().MustExec("INSERT INTO zip_604(pi_id, pi_name) VALUES(61, 'Here-1')")
	dbCtrl.SqlxDb().MustExec("INSERT INTO zip_604(pi_id, pi_name) VALUES(62, 'Here-2')")

	testedData := make(map[string]interface{})
	dbCtrl.QueryRowxAndMapScan(
		testedData,
		"SELECT * FROM zip_604 WHERE pi_name = 'Here-1'", nil,
	)

	c.Logf("QueryRowxAndMapScan: %#v", testedData)
	c.Assert(testedData["pi_id"], Equals, int64(61))
	c.Assert(string(testedData["pi_name"].([]uint8)), Equals, "Here-1")
}

// Tests the query(rowx) and scan(on struct) for db controller
func (suite *TestWrapperSuite) TestQueryRowxAndStructScanOfDbCtrl(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE zoo_604(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")
	dbCtrl.SqlxDb().MustExec("INSERT INTO zoo_604(pi_id, pi_name) VALUES(61, 'Here-1')")
	dbCtrl.SqlxDb().MustExec("INSERT INTO zoo_604(pi_id, pi_name) VALUES(62, 'Here-2')")

	testedData := struct {
		Id   int    `db:"pi_id"`
		Name string `db:"pi_name"`
	}{}
	dbCtrl.QueryRowxAndStructScan(
		&testedData,
		"SELECT * FROM zoo_604 WHERE pi_name = 'Here-1'", nil,
	)

	c.Logf("QueryRowxAndStructScan: %#v", testedData)
	c.Assert(testedData.Id, Equals, 61)
	c.Assert(testedData.Name, Equals, "Here-1")
}

// Tests the query(rowx) and scan(on slice) for db controller
func (suite *TestWrapperSuite) TestQueryRowxAndSliceScanOfDbCtrl(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE gta_604(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")
	dbCtrl.SqlxDb().MustExec("INSERT INTO gta_604(pi_id, pi_name) VALUES(61, 'Here-1')")
	dbCtrl.SqlxDb().MustExec("INSERT INTO gta_604(pi_id, pi_name) VALUES(62, 'Here-2')")

	testedData := dbCtrl.QueryRowxAndSliceScan(
		"SELECT * FROM gta_604 WHERE pi_name = 'Here-1'", nil,
	)

	c.Logf("QueryRowxAndSliceScan: %#v", testedData)
	c.Assert(testedData[0], Equals, int64(61))
	c.Assert(string(testedData[1].([]uint8)), Equals, "Here-1")
}

// Tests the getting of data by named statement
func (suite *TestWrapperSuite) TestGetOfNamedStmt(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE cloc_604(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")
	dbCtrl.SqlxDb().MustExec("INSERT INTO cloc_604(pi_id, pi_name) VALUES(61, 'Here-1')")
	dbCtrl.SqlxDb().MustExec("INSERT INTO cloc_604(pi_id, pi_name) VALUES(62, 'Here-2')")

	stmtExt := dbCtrl.PrepareNamedExt(
		"SELECT * FROM cloc_604 WHERE pi_name = :name",
	)

	testedData := struct {
		Id   int    `db:"pi_id"`
		Name string `db:"pi_name"`
	}{}
	stmtExt.Get(&testedData, map[string]interface{}{"name": "Here-1"})

	c.Logf("NamedStmt Get: %#v", testedData)
	c.Assert(testedData.Id, Equals, 61)
	c.Assert(testedData.Name, Equals, "Here-1")
}

// Tests the getting of data by statement
func (suite *TestWrapperSuite) TestGetOfStmt(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE cloc_604(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")
	dbCtrl.SqlxDb().MustExec("INSERT INTO cloc_604(pi_id, pi_name) VALUES(61, 'Here-1')")
	dbCtrl.SqlxDb().MustExec("INSERT INTO cloc_604(pi_id, pi_name) VALUES(62, 'Here-2')")

	stmtExt := dbCtrl.PreparexExt(
		"SELECT * FROM cloc_604 WHERE pi_name = ?",
	)

	testedData := struct {
		Id   int    `db:"pi_id"`
		Name string `db:"pi_name"`
	}{}
	stmtExt.Get(&testedData, "Here-1")

	c.Logf("Stmt Get: %#v", testedData)
	c.Assert(testedData.Id, Equals, 61)
	c.Assert(testedData.Name, Equals, "Here-1")
}

// Tests the scanning(on map) for RowExt upon no row case
func (suite *TestWrapperSuite) TestMapScanOrNoRowOnRowExt(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE kc_8871(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")

	testedRowx := dbCtrl.QueryRowxExt("SELECT * FROM kc_8871 WHERE pi_id = 981")

	c.Assert(testedRowx.MapScanOrNoRow(nil), Equals, false)
}

// Tests the scanning(on value) for RowExt upon no row case
func (suite *TestWrapperSuite) TestScanOrNoRowOnRowExt(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE kcv_8871(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")

	testedRowx := dbCtrl.QueryRowxExt("SELECT * FROM kcv_8871 WHERE pi_id = 981")

	c.Assert(testedRowx.ScanOrNoRow(), Equals, false)
}

// Tests the scanning(on slice) for RowExt upon no row case
func (suite *TestWrapperSuite) TestSliceScanOrNoRowOnRowExt(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE kcs_8871(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")

	testedRowx := dbCtrl.QueryRowxExt("SELECT * FROM kcs_8871 WHERE pi_id = 981")

	_, hasRow := testedRowx.SliceScanOrNoRow()
	c.Assert(hasRow, Equals, false)
}

// Tests the scanning(on struct) for RowExt upon no row case
func (suite *TestWrapperSuite) TestStructScanOrNoRowOnRowExt(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE kcst_8871(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")

	testedRowx := dbCtrl.QueryRowxExt("SELECT * FROM kcst_8871 WHERE pi_id = 981")

	sampleValue := &struct {
		Id   int    `db:"pi_id"`
		Name string `db:"pi_name"`
	}{}
	c.Assert(testedRowx.StructScanOrNoRow(sampleValue), Equals, false)
}

// Tests the InTx() callback for db controller
func (suite *TestWrapperSuite) TestInTxOfDbCtrl(c *C) {
	dbCtrl := buildDbCtrl(c)

	dbCtrl.SqlxDb().MustExec("CREATE TABLE txol_8871(pi_id INT PRIMARY KEY, pi_name VARCHAR(32))")

	dbCtrl.InTx(TxCallbackFunc(func(tx *sqlx.Tx) db.TxFinale {
		tx.MustExec("INSERT INTO txol_8871(pi_id, pi_name) VALUES(81, 'easy-1')")
		tx.MustExec("INSERT INTO txol_8871(pi_id, pi_name) VALUES(82, 'easy-2')")
		return db.TxCommit
	}))

	var numberOfData int
	dbCtrl.QueryRowxAndScan(
		"SELECT COUNT(*) FROM txol_8871", nil, &numberOfData,
	)

	c.Assert(numberOfData, Equals, 2)
}

func buildDbCtrl(c *C) *DbController {
	db, err := sqlx.Connect("sqlite3", ":memory:")

	c.Assert(err, IsNil)

	return NewDbController(db)
}
