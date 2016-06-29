package conn_pool

import (
	"fmt"
	jsonrpc2 "github.com/Cepave/rpc-codec/jsonrpc2"
	"log"
	"time"
)

type JsonRpcV2Client struct {
	cli  *jsonrpc2.Client
	name string
}

func (this JsonRpcV2Client) Name() string {
	return this.name
}

func (this JsonRpcV2Client) Closed() bool {
	return this.cli == nil
}

func (this JsonRpcV2Client) Close() error {
	if this.cli != nil {
		err := this.cli.Close()
		this.cli = nil
		return err
	}
	return nil
}

func (this JsonRpcV2Client) Call(method string, args interface{}, reply interface{}) error {
	return this.cli.Call(method, args, reply)
}

type NqmRpcConnPoolHelper struct {
	p           *ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

func newNqmConnPool(address string, maxConns int, maxIdle int) *ConnPool {
	pool := NewConnPool("httpRpcClient", address, maxConns, maxIdle)

	pool.New = func(connName string) (NConn, error) {
		url := "http://" + address
		log.Println("Invoke rpcClient2 by url: " + url)
		return JsonRpcV2Client{cli: jsonrpc2.NewHTTPClient(url), name: connName}, nil
	}

	return pool
}

func NewNqmRpcConnPoolHelper(address string, maxConns int, maxIdle int, connTimeout int, callTimeout int) *NqmRpcConnPoolHelper {
	return &NqmRpcConnPoolHelper{
		p:           newNqmConnPool(address, maxConns, maxIdle),
		maxConns:    maxConns,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

func (this *NqmRpcConnPoolHelper) Destroy() {
	if this.p != nil {
		this.p.Destroy()
	}
}

func (this *NqmRpcConnPoolHelper) Call(method string, args interface{}, resp interface{}) error {
	conn, err := this.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, this.p.Proc())
	}

	rpcClient := conn.(JsonRpcV2Client)

	done := make(chan error)
	go func() {
		done <- rpcClient.Call(method, args, resp)
	}()

	select {
	case <-time.After(time.Duration(this.callTimeout) * time.Millisecond):
		this.p.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", this.address)
	case err = <-done:
		if err != nil {
			this.p.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", this.address, err, this.p.Proc())
		} else {
			this.p.Release(conn)
		}
		return err
	}
}
