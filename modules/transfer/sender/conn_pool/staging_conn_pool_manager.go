package conn_pool

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"time"
)

func newStagingConnPool(address string, maxConns int, maxIdle int, connTimeout int) *ConnPool {
	connectionTimeout := time.Duration(connTimeout) * time.Millisecond
	p := NewConnPool("staging", address, maxConns, maxIdle)

	p.New = func(connName string) (NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connectionTimeout)
		if err != nil {
			return nil, err
		}

		return RpcClient{cli: jsonrpc.NewClient(conn), name: connName}, nil
	}

	return p
}

type StagingConnPoolHelper struct {
	p           *ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

func NewStagingConnPoolHelper(address string, maxConns, maxIdle, connTimeout, callTimeout int) *StagingConnPoolHelper {
	return &StagingConnPoolHelper{
		p:           newStagingConnPool(address, maxConns, maxIdle, connTimeout),
		maxConns:    maxConns,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

// A synchronous call; return if completed or time-out
func (this *StagingConnPoolHelper) Call(method string, args interface{}, resp interface{}) error {
	conn, err := this.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, this.p.Proc())
	}

	rpcClient := conn.(RpcClient)
	callTimeout := time.Duration(this.callTimeout) * time.Millisecond

	done := make(chan error)
	go func() {
		done <- rpcClient.Call(method, args, resp)
	}()

	select {
	case <-time.After(callTimeout):
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

func (this *StagingConnPoolHelper) Destroy() {
	if this.p != nil {
		this.p.Destroy()
	}
}
