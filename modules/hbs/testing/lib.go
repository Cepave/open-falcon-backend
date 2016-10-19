package testing

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"flag"
	tknet "github.com/toolkits/net"
	. "gopkg.in/check.v1"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
	check "gopkg.in/check.v1"
)

var dsnMysql = flag.String("dsn_mysql", "", "Dsn of Mysql")

var DbForTest *sql.DB

// This function accepts a checker struct and
// it is executed by testing
type AssertFunc func(*check.C)

type InitDbCallback func(dsn string, maxIdle int) error
type ReleaseDbCallback func()

func DoInitDb(callback InitDbCallback) {
	if *dsnMysql == "" {
		return
	}

	if err := callback(*dsnMysql, 2)
		err != nil {
		panic(err)
	}
}
func DoReleaseDb(callback ReleaseDbCallback) {
	if *dsnMysql == "" {
		return
	}

	callback()
}

// Checks whether or not skipping testing by viable arguments
func HasDbEnvForMysqlOrSkip(c *C) bool {
	var hasMySqlDsn = *dsnMysql != ""

	if !hasMySqlDsn {
		c.Skip("Skip Mysql Test: -dsn_mysql=<dsn>")
	}

	return hasMySqlDsn
}

// IoC execution
func ExecuteOrFail(query string, args ...interface{}) sql.Result {
	result, err := DbForTest.Exec(query, args...)

	if err != nil {
		log.Fatalf("Execute SQL error. %v. SQL:\n%v\n", err, query)
	}

	return result
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
