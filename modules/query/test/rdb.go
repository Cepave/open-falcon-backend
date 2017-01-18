package test

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	tknet "github.com/toolkits/net"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	testDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
	"gopkg.in/check.v1"
)

var DbForTest *sql.DB

func InitDb(c *check.C) {
	dbFacade := testDb.InitDbFacade(c)
	owlDb.DbFacade = dbFacade
	nqmDb.DbFacade = dbFacade

	DbForTest = dbFacade.SqlDb
}

func ReleaseDb(c *check.C) {
	testDb.ReleaseDbFacade(c, owlDb.DbFacade)

	DbForTest = nil
	owlDb.DbFacade = nil
	nqmDb.DbFacade = nil
}

// IoC Callback for rows
// This method would use sql.DB.Query method to retrive data
func QueryForRows(
	rowsCallback func(row *sql.Rows),
	sqlQuery string, args ...interface{},
) {
	rows, err := DbForTest.Query(
		sqlQuery, args...,
	)

	if err != nil {
		log.Fatalf("Query SQL error. %v. SQL:\n%v\n", err, sqlQuery)
	}

	defer rows.Close()
	for rows.Next() {
		rowsCallback(rows)
	}
}

// IoC Callback for row
// This method would use sql.DB.QueryRow method to retrive data
func QueryForRow(
	rowCallback func(row *sql.Row),
	sqlQuery string, args ...interface{},
) {
	row := DbForTest.QueryRow(
		sqlQuery, args...,
	)

	rowCallback(row)
}

// IoC execution(in transaction)
func ExecuteInTx(txCallback func(*sql.Tx)) {
	var tx *sql.Tx
	var err error

	if tx, err = DbForTest.Begin(); err != nil {
		log.Fatalf("Cannot create transaction: %v", err)
	}

	defer tx.Commit()
	txCallback(tx)
}

// IoC execution
func ExecuteOrFail(query string, args ...interface{}) sql.Result {
	result, err := DbForTest.Exec(query, args...)

	if err != nil {
		log.Fatalf("Execute SQL error. %v. SQL:\n%v\n", err, query)
	}

	return result
}

// IoC execution(in transaction)
func ExecuteQueriesOrFailInTx(queries ...string) {
	ExecuteInTx(
		func(tx *sql.Tx) {
			for _, v := range queries {
				if _, err := tx.Exec(v); err != nil {
					log.Fatalf("Execute SQL error. %v. SQL:\n%v\n", err, v)
				}
			}
		},
	)
}

type RpcTestEnv struct {
	Port      int
	RpcClient *rpc.Client

	stop bool
	wait chan bool
}

func DefaultListenAndExecute(
	rcvr interface{}, rpcCallback func(*RpcTestEnv),
) {
	var rpcTestEnvInstance = RpcTestEnv{
		Port: 18080,
		stop: false,
		wait: make(chan bool, 1),
	}

	rpcTestEnvInstance.ListenAndExecute(
		rcvr, rpcCallback,
	)
}
func (rpcTestEnvInstance *RpcTestEnv) ListenAndExecute(
	rcvr interface{}, rpcCallback func(*RpcTestEnv),
) {
	server := rpc.NewServer()
	server.Register(rcvr)

	var address = fmt.Sprintf("localhost:%d", rpcTestEnvInstance.Port)

	/**
	 * Listening RPC
	 */
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Errorf("Cannot listen TCP[%d]. Error: %v", rpcTestEnvInstance.Port, err))
	}

	log.Printf("Listen RPC at port [%d]", rpcTestEnvInstance.Port)
	// :~)

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				if rpcTestEnvInstance.stop {
					log.Printf("Stop RPC server")
					rpcTestEnvInstance.wait <- true
				}

				break
			}

			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()

	/**
	 * Initialize RPC client
	 */
	rpcTestEnvInstance.RpcClient, err = tknet.JsonRpcClient("tcp", address, time.Second*3)
	if err != nil {
		panic(fmt.Errorf("Initialize RPC client error: %v", err))
	}
	log.Printf("Construct RPC Client")
	// :~)

	defer func() {
		rpcTestEnvInstance.RpcClient.Close()
		rpcTestEnvInstance.stop = true
		listener.Close()
		<-rpcTestEnvInstance.wait
	}()

	rpcCallback(rpcTestEnvInstance)
}
