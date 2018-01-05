// Provides various utilities to construct listener
package net

import (
	"net"
	"sync"
)

// Initializes TCP listener
//
// error object would be viable if anything get failed
func InitTcpListener(address string) (net.Listener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	logger.Infof("Init TCP Address[ %s ]", address)

	return listener, nil
}

// [MUST] Initializes TCP listener
//
// Panic would raise if anything get failed
func MustInitTcpListener(address string) net.Listener {
	listener, err := InitTcpListener(address)

	if err != nil {
		logger.Fatalf("Cannot initialize listener: %v", err)
		panic(err.Error())
	}

	return listener
}

// Constructs the controller for usage on "net.Listener"
func NewListenerController(listener net.Listener) *ListenerController {
	return &ListenerController{
		Listener: listener,
		working:  false,
		lock:     &sync.Mutex{},
	}
}

type ListenerController struct {
	net.Listener

	working bool
	lock    *sync.Mutex
}

// This method would keep accepting message of socket
//
// This method would use go routine to call your conn handler.
func (c *ListenerController) AcceptLoop(connHandler func(conn net.Conn)) {
	c.lock.Lock()
	if c.working {
		return
	}
	c.working = true
	c.lock.Unlock()

	for {
		conn, err := c.Accept()
		if !c.working {
			break
		}

		if err != nil {
			logger.Errorf("Accept message has error: %v", err)
			continue
		}

		go func(lambdaConn net.Conn) {
			defer func() {
				p := recover()
				if p != nil {
					logger.Panicf("Connection handler has error: %v", p)
				}
			}()

			connHandler(lambdaConn)
		}(conn)
	}
}

// Close this controller if any looping is running
func (c *ListenerController) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.working {
		return
	}

	c.working = false

	err := c.Listener.Close()
	if err != nil {
		logger.Errorf("Close listener[%v] has error: %v", c.Listener.Addr(), err)
	}
}
