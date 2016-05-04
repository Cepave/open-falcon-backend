package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"reflect"
	"testing"

	"github.com/Cepave/common/model"
)

type NqmAgent int

func (t *NqmAgent) PingTask(request model.NqmPingTaskRequest, response *model.NqmPingTaskResponse) (err error) {
	response.NeedPing = true

	nqmAgent := &model.NqmAgent{
		Id:           1,
		Name:         "agent_for_test",
		IspId:        2,
		IspName:      "IspName_for_test",
		ProvinceId:   3,
		ProvinceName: "ProvinceName_for_test",
		CityId:       4,
		CityName:     "CityName_for_test",
	}
	response.Agent = nqmAgent

	targets := []model.NqmTarget{
		{Id: 11, IspId: 12, ProvinceId: 13, CityId: 14},
		{Id: 21, IspId: 22, ProvinceId: 23, CityId: 24},
		{Id: 31, IspId: 32, ProvinceId: 33, CityId: 34},
	}
	response.Targets = targets

	response.Command = []string{"fping", "-p", "20", "-i", "10", "-C", "4", "-q", "-a"}
	return nil
}

func initJsonRpcServer(addr string) {
	rpc.Register(new(NqmAgent))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		panic(fmt.Errorf("Cannot listen to address %v. Error %v", addr, e))
	} else {
		log.Println("RPC server is listening at", addr)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println("listener accept fail:", err)
				break
			}
			go jsonrpc.ServeConn(conn)
		}
	}()
}

func initJsonRpcClient(srvAddr string) {
	rpcClient.RpcServer = srvAddr
	req = model.NqmPingTaskRequest{
		Hostname:     "nqm-agent",
		IpAddress:    "1.2.3.4",
		ConnectionId: "arg-arg-arg",
	}
}

func TestQueryTask(t *testing.T) {
	initJsonRpcServer("127.0.0.1:65534")
	initJsonRpcClient("127.0.0.1:65534")

	cmd, targetList, agent, err := QueryTask()
	if err != nil {
		t.Error(err)
	}

	expecteds := []model.NqmTarget{
		{Id: 11, IspId: 12, ProvinceId: 13, CityId: 14},
		{Id: 21, IspId: 22, ProvinceId: 23, CityId: 24},
		{Id: 31, IspId: 32, ProvinceId: 33, CityId: 34},
	}
	for i, v := range targetList {
		if !reflect.DeepEqual(expecteds[i], v) {
			t.Error(v)
		}
	}

	t.Log(agent)
	t.Log(cmd)
}
