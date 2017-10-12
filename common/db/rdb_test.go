package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	. "gopkg.in/check.v1"
)

type TestRdbSuite struct{}

var _ = Suite(&TestRdbSuite{})

// Tests the panic(no panic handler)
func (suite *TestRdbSuite) TestOperateOnDbWithPanic(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	var testedFunc = func() {
		testedCtrl.OperateOnDb(DbCallbackFunc(func(db *sql.DB) {
			panic("Test Panic")
		}))
	}

	defer func() {
		c.Assert(testedFunc, PanicMatches, ".*Test Panic.*")
	}()
}

// Tests the operate on database
func (suite *TestRdbSuite) TestOperateOnDb(c *C) {
	var getCalled bool = false

	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	testedCtrl.OperateOnDb(DbCallbackFunc(func(db *sql.DB) {
		getCalled = true

		_, err := db.Exec("CREATE TABLE test_on_db(tob_id INT PRIMARY KEY)")
		c.Assert(err, IsNil)
	}))
	c.Assert(getCalled, Equals, true)
}

// Tests the releasing for resources on controller
//
// Also tests the re-release(panic)
func (suite *TestRdbSuite) TestRelease(c *C) {
	testedCtrl := buildSampleDbController(c)

	testedCtrl.Release()

	defer func() {
		panicObject := recover()
		c.Assert(panicObject, NotNil)
	}()

	testedCtrl.Release()
}

// Tests the register of handler
func (suite *TestRdbSuite) TestRegisterPanicHandler(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	var getCalled bool = false
	testedCtrl.RegisterPanicHandler(func(panicValue interface{}) {
		c.Logf("Got panic: %v", panicValue)
		getCalled = true
	})

	testedCtrl.Exec("No Such SQL Stmt")

	c.Assert(getCalled, Equals, true)
}

// Tests the execute of SQL
func (suite *TestRdbSuite) TestExec(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	testedCtrl.Exec("CREATE TABLE test_exec_1(te_id int PRIMARY KEY)")

	testedResult := testedCtrl.Exec("INSERT INTO test_exec_1 VALUES(20)")

	rowAffected, err := testedResult.RowsAffected()

	c.Assert(err, IsNil)
	c.Assert(rowAffected, Equals, int64(1))
}

// Tests the query for rows
func (suite *TestRdbSuite) TestQueryForRows(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	testedCtrl.Exec(
		"CREATE TABLE test_rows(tr_id INT PRIMARY KEY, tr_text VARCHAR(64) NOT NULL)",
	)
	testedCtrl.Exec(
		"INSERT INTO test_rows VALUES(1, 'v-1'), (2, 'v-2'), (3, 'v-3')",
	)

	testedNumberOfRows := testedCtrl.QueryForRows(
		RowsCallbackFunc(func(rows *sql.Rows) IterateControl {
			return IterateContinue
		}),
		"SELECT * FROM test_rows",
	)

	c.Assert(testedNumberOfRows, Equals, uint(3))
}

// Tests the query for row
func (suite *TestRdbSuite) TestQueryForRow(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	testedCtrl.Exec("CREATE TABLE test_row(tr_id INT PRIMARY KEY, tr_text VARCHAR(64) NOT NULL)")
	testedCtrl.Exec("INSERT INTO test_row VALUES(1, 'v-1'), (2, 'v-2'), (3, 'v-3')")

	var numberOfRows int
	testedCtrl.QueryForRow(
		RowCallbackFunc(func(row *sql.Row) {
			ToRowExt(row).Scan(&numberOfRows)
		}),
		"SELECT COUNT(*) FROM test_row",
	)

	c.Assert(numberOfRows, Equals, 3)
}

// Tests the executing in transaction
func (suite *TestRdbSuite) TestInTx(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	testedCtrl.Exec(
		"CREATE TABLE test_in_tx(it_id INT PRIMARY KEY, it_text VARCHAR(64) NOT NULL)",
	)

	/**
	 * Builds committed data
	 */
	testedCtrl.InTx(TxCallbackFunc(func(tx *sql.Tx) TxFinale {
		txExt := ToTxExt(tx)

		txExt.Exec("INSERT INTO test_in_tx VALUES(21, 'v-21')")
		txExt.Exec("INSERT INTO test_in_tx VALUES(22, 'v-22')")
		txExt.Exec("INSERT INTO test_in_tx VALUES(23, 'v-23')")

		return TxCommit
	}))
	assertNumberOfDataInTx(c, testedCtrl, 3)
	// :~)

	/**
	 * Builds rollbacked data
	 */
	var testedFunc = func() {
		testedCtrl.InTx(TxCallbackFunc(func(tx *sql.Tx) TxFinale {
			txExt := ToTxExt(tx)

			txExt.Exec("INSERT INTO test_in_tx VALUES(123, 'v-123')")
			txExt.Exec("INSERT INTO test_in_tx VALUES(124, 'v-124')")
			txExt.Exec("INSERT INTO test_in_tx VALUES(21, 'v-21')")

			return TxCommit
		}))
	}

	c.Assert(testedFunc, PanicMatches, ".*UNIQUE.*")
	assertNumberOfDataInTx(c, testedCtrl, 3)
	// :~)
}

// Tests the calling of if callbacks in transaction
func (suite *TestRdbSuite) TestInTxForIf(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	testedTrueSample := &ifSample{true, false}
	testedFalseSample := &ifSample{false, false}

	testedCtrl.InTxForIf(testedTrueSample)
	testedCtrl.InTxForIf(testedFalseSample)

	c.Assert(testedTrueSample.getCalled, Equals, true)
	c.Assert(testedFalseSample.getCalled, Equals, false)
}

// Tests the executing of queries in transaction
func (suite *TestRdbSuite) TestExecQueriesInTx(c *C) {
	testedCtrl := buildSampleDbController(c)
	defer testedCtrl.Release()

	testedCtrl.Exec(
		"CREATE TABLE car_981 ( c_id INT PRIMARY KEY, c_name VARCHAR(64) NOT NULL )",
	)

	testedCtrl.ExecQueriesInTx(
		"INSERT INTO car_981 VALUES(1, 'OK-1')",
		"INSERT INTO car_981 VALUES(2, 'OK-2')",
		"INSERT INTO car_981 VALUES(3, 'OK-3')",
	)

	var numberOfRows int
	testedCtrl.QueryForRow(
		RowCallbackFunc(func(row *sql.Row) {
			ToRowExt(row).Scan(&numberOfRows)
		}),
		"SELECT COUNT(*) FROM car_981",
	)

	c.Assert(numberOfRows, Equals, 3)
}

type ifSample struct {
	ifValue   bool
	getCalled bool
}

func (self *ifSample) BootCallback(tx *sql.Tx) bool {
	return self.ifValue
}
func (self *ifSample) IfTrue(tx *sql.Tx) {
	self.getCalled = true
}

func assertNumberOfDataInTx(
	c *C,
	testedCtrl *DbController, expectedResult int,
) {
	var testedNumber int
	testedCtrl.QueryForRow(
		RowCallbackFunc(func(row *sql.Row) {
			err := row.Scan(&testedNumber)
			c.Assert(err, IsNil)
		}),
		"SELECT COUNT(*) FROM test_in_tx",
	)

	c.Assert(testedNumber, Equals, expectedResult)
}

func buildSampleDbController(c *C) *DbController {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		c.Fatalf("Open Databae Error. %v", err)
	}

	return NewDbController(db)
}
