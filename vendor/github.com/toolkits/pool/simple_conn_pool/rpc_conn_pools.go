package simple_conn_pool

import (
	"net"
	"net/rpc"
	"sync"
	"time"
)

// TODO 以下的部分, 应该放在一个公共库文件中
// RpcCient, 要实现io.Closer接口
type RpcClient struct {
	cli  *rpc.Client
	name string
}

func (this RpcClient) Name() string {
	return this.name
}

func (this RpcClient) Closed() bool {
	return this.cli == nil
}

func (this RpcClient) Close() error {
	if this.cli != nil {
		err := this.cli.Close()
		this.cli = nil
		return err
	}
	return nil
}

func (this RpcClient) Call(method string, args interface{}, reply interface{}) error {
	return this.cli.Call(method, args, reply)
}

// ConnPools Manager
type SafeRpcConnPools struct {
	sync.RWMutex
	M           map[string]*ConnPool
	MaxConns    int32
	MaxIdle     int32
	ConnTimeout int32
	CallTimeout int32
}

func CreateSafeRpcConnPools(maxConns, maxIdle, connTimeout, callTimeout int32, cluster []string) *SafeRpcConnPools {
	cp := &SafeRpcConnPools{M: make(map[string]*ConnPool), MaxConns: maxConns, MaxIdle: maxIdle,
		ConnTimeout: connTimeout, CallTimeout: callTimeout}

	ct := time.Duration(cp.ConnTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOnePool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

func (this *SafeRpcConnPools) Proc() []string {
	procs := []string{}
	for _, cp := range this.M {
		procs = append(procs, cp.Proc())
	}
	return procs
}

func (this *SafeRpcConnPools) Get(address string) (*ConnPool, bool) {
	this.RLock()
	defer this.RUnlock()
	p, exists := this.M[address]
	return p, exists
}

func (this *SafeRpcConnPools) Destroy() {
	this.Lock()
	defer this.Unlock()
	addresses := make([]string, 0, len(this.M))
	for address := range this.M {
		addresses = append(addresses, address)
	}

	for _, address := range addresses {
		this.M[address].Destroy()
		delete(this.M, address)
	}
}

func createOnePool(name string, address string, connTimeout time.Duration, maxConns, maxIdle int32) *ConnPool {
	p := NewConnPool(name, address, maxConns, maxIdle)
	p.New = func(connName string) (NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connTimeout)
		if err != nil {
			return nil, err
		}

		return RpcClient{cli: rpc.NewClient(conn), name: connName}, nil
	}

	return p
}
