package pool

import (
	"fmt"
	"io"
	"sync"
)

var ErrMaxConn = fmt.Errorf("maximum connections reached")

// ConnPool manages the life cycle of connections
type ConnPool struct {
	sync.RWMutex

	// New is used to create a new connection when necessary.
	New func() (io.Closer, error)

	// Ping is use to check the conn fetched from pool
	Ping func(io.Closer) error

	Address  string
	MaxConns int
	MaxIdle  int

	TestOnBorrow bool

	active int
	free   []io.Closer
}

func Create(address string, maxConns int, maxIdle int) *ConnPool {
	return &ConnPool{
		Address:  address,
		MaxConns: maxConns,
		MaxIdle:  maxIdle,
	}
}

func (this *ConnPool) Get() (conn io.Closer, err error) {
	conn = this.tryFree()
	if conn != nil {
		if this.TestOnBorrow {
			err = this.Ping(conn)
			if err != nil {
				this.decreActive() //bug fix: 如果不加decrease, 当 free不为空 + server端重启时 就会导致active大于实际情况不符
				conn.Close()
				conn = this.tryFree()
				err = nil
			}
		}

		if conn != nil {
			return
		}
	}

	if this.reachedMax() {
		return nil, ErrMaxConn
	}

	conn, err = this.New()
	if err != nil {
		return
	}

	if this.TestOnBorrow {
		err = this.Ping(conn)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	this.increActive()
	return
}

func (this *ConnPool) Release(conn io.Closer) {
	if this.overMaxIdle() {
		this.decreActive()
		if conn != nil {
			conn.Close()
		}
	} else {
		this.Lock()
		defer this.Unlock()
		this.free = append(this.free, conn)
	}
}

func (this *ConnPool) ForceClose(conn io.Closer) {
	this.decreActive()
	if conn != nil {
		conn.Close()
	}
}

func (this *ConnPool) Destroy() {
	this.Lock()
	defer this.Unlock()

	for _, conn := range this.free {
		if conn != nil {
			conn.Close()
		}
	}
}

func (this *ConnPool) tryFree() io.Closer {
	this.Lock()
	defer this.Unlock()

	if len(this.free) == 0 {
		return nil
	}

	conn := this.free[0]
	this.free = this.free[1:]

	return conn
}

func (this *ConnPool) reachedMax() bool {
	this.RLock()
	defer this.RUnlock()
	return this.active >= this.MaxConns
}

func (this *ConnPool) increActive() {
	this.Lock()
	defer this.Unlock()
	this.active += 1
}

func (this *ConnPool) decreActive() {
	this.Lock()
	defer this.Unlock()
	this.active -= 1
}

func (this *ConnPool) overMaxIdle() bool {
	this.RLock()
	defer this.RUnlock()
	return len(this.free) >= this.MaxIdle
}
