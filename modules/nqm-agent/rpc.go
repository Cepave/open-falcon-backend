package main

import (
	"math"
	"net/rpc"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/toolkits/net"
)

var HbsRespTime time.Time

var (
	req        model.NqmTaskRequest
	rpcClient  *rpc.Client
	rpcServer  string
	rpcTimeout time.Duration
)
var retryCnt = int(1)

func closeConn() {
	if err := rpcClient.Close(); err != nil {
		log.Errorln("[ hbs ]", err)
	}
	rpcClient = nil
}

func wait4Retry() {
	if retryCnt > 6 {
		retryCnt = 1
	}
	time.Sleep(time.Duration(math.Pow(2.0, float64(retryCnt))) * time.Second)
	retryCnt++
}

func hasConn() bool {
	if rpcClient != nil {
		return true
	}
	return false
}

func initConn(server string, timeout time.Duration) *rpc.Client {
	for {
		client, err := net.JsonRpcClient("tcp", server, timeout)
		if err == nil {
			return client
		}
		log.Println("[ hbs ] Error on query:", err)

		dur := int64(time.Since(HbsRespTime).Minutes())
		var nilTime time.Time
		if HbsRespTime.Add(6*time.Minute).Before(time.Now()) && HbsRespTime != nilTime {
			log.Infof("[ hbs ] Last response: %dm ago\n", dur)
		}

		wait4Retry()
	}
}

func RPCCall(method string, args interface{}, reply interface{}) error {
	if !hasConn() {
		initConn(rpcServer, rpcTimeout)
	}

	done := make(chan error, 1)
	go func() {
		done <- rpcClient.Call(method, args, reply)
	}()

	callTimeout := time.Duration(50 * time.Second)
	select {
	case <-time.After(callTimeout):
		log.Warnf("rpc call timeout %v => %v", rpcClient, rpcServer)
		closeConn()
	case err := <-done:
		if err != nil {
			closeConn()
			return err
		}
	}

	return nil
}

func InitRPC() {
	rpcServer = GetGeneralConfig().Hbs.RPCServer
	req = model.NqmTaskRequest{
		Hostname:     GetGeneralConfig().Hostname,
		IpAddress:    GetGeneralConfig().IPAddress,
		ConnectionId: GetGeneralConfig().ConnectionID,
	}
	rpcClient = initConn(rpcServer, rpcTimeout)
}
