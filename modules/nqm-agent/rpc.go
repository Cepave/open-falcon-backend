package main

import (
	"fmt"
	"math"
	"net/rpc"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/toolkits/net"
)

const ConnTimeout = 50 * time.Second

var HbsRespTime time.Time

var (
	req       model.NqmTaskRequest
	rpcServer string
)
var retryCnt = int(1)

func closeConn(c *rpc.Client) {
	if err := c.Close(); err != nil {
		log.Errorln("[ hbs ] Error on closing RPC connection", err)
	}
}

func wait4Retry() {
	if retryCnt > 6 {
		retryCnt = 1
	}
	time.Sleep(time.Duration(math.Pow(2.0, float64(retryCnt))) * time.Second)
	retryCnt++
}

func initConn(server string) *rpc.Client {
	for {
		client, err := net.JsonRpcClient("tcp", server, ConnTimeout)
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
	currentRpcClient := initConn(rpcServer)
	defer closeConn(currentRpcClient)

	done := make(chan error, 1)
	go func() {
		done <- currentRpcClient.Call(method, args, reply)
	}()

	select {
	case <-time.After(ConnTimeout):
		return fmt.Errorf("Call to server <%s> timed out (%d seconds)", rpcServer, ConnTimeout)
	case err := <-done:
		return err
	}
}

func InitRPC() {
	rpcServer = JSONConfig().Hbs.RPCServer
	req = model.NqmTaskRequest{
		Hostname:     GetGeneralConfig().Hostname,
		IpAddress:    GetGeneralConfig().IPAddress,
		ConnectionId: GetGeneralConfig().ConnectionID,
	}
}
