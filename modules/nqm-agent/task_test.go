package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
)

type NqmAgent int

func (t *NqmAgent) Task(request model.NqmTaskRequest, response *model.NqmTaskResponse) (err error) {
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
		{Id: 11, Host: "11.11.11.11", IspId: 12, ProvinceId: 13, CityId: 14},
		{Id: 21, Host: "22.22.22.22", IspId: 22, ProvinceId: 23, CityId: 24},
		{Id: 31, Host: "33.33.33.33", IspId: 32, ProvinceId: 33, CityId: 34},
	}
	response.Targets = targets

	response.Measurements = map[string]model.MeasurementsProperty{
		"fping":   {true, []string{"fping", "-p", "20", "-i", "10", "-C", "4", "-q", "-a"}, 300},
		"tcpping": {false, []string{"tcpping", "-i", "0.01", "-c", "4"}, 300},
		"tcpconn": {false, []string{"tcpconn"}, 300},
	}
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
	rpcServer = srvAddr
	req = model.NqmTaskRequest{
		Hostname:     "nqm-agent",
		IpAddress:    "1.2.3.4",
		ConnectionId: "arg-arg-arg",
	}
	rpcClient = initConn(rpcServer, rpcTimeout)
}

func TestTask(t *testing.T) {
	GetGeneralConfig().hbsResp.Store(model.NqmTaskResponse{})
	initJsonRpcServer("127.0.0.1:65534")
	initJsonRpcClient("127.0.0.1:65534")

	query()
	hbsResp := GetGeneralConfig().hbsResp.Load().(model.NqmTaskResponse)
	cmd, targets, agent, _, err := Task(new(Fping))
	if err != nil {
		t.Error(err)
	}

	expecteds := []model.NqmTarget{
		{Id: 11, Host: "11.11.11.11", IspId: 12, ProvinceId: 13, CityId: 14},
		{Id: 21, Host: "22.22.22.22", IspId: 22, ProvinceId: 23, CityId: 24},
		{Id: 31, Host: "33.33.33.33", IspId: 32, ProvinceId: 33, CityId: 34},
	}

	expectedAgent := *hbsResp.Agent
	expectedCommand := []string{
		"fping", "-p", "20", "-i", "10", "-C", "4", "-q", "-a",
		"11.11.11.11", "22.22.22.22", "33.33.33.33",
	}

	for i, v := range targets {
		if !reflect.DeepEqual(expecteds[i], v) {
			t.Error(expecteds[i], v)
		}
		t.Log(expecteds[i], v)
	}
	if !reflect.DeepEqual(expectedAgent, agent) {
		t.Error(expectedAgent, agent)
	}
	t.Log(agent)

	if &expectedAgent == &agent {
		t.Errorf("%p == %p", &expectedAgent, &agent)
	}
	t.Logf("%p != %p", &agent, &expectedAgent)

	if !reflect.DeepEqual(expectedCommand, cmd) {
		t.Error(expectedCommand, cmd)
	}
	t.Log(cmd)
}
