package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"

	"github.com/Cepave/open-falcon-backend/common/model"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	testJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"
	. "gopkg.in/check.v1"
)

type TestNqmAgentSuite struct{}

var _ = Suite(&TestNqmAgentSuite{})

var ts *httptest.Server

var tsCnt = 0

var tsRsps = []string{
	`
	{
		"id" : 405001,
		"name": "ag-name-1",
		"connection_id": "ag-rpc-1",
		"hostname": "rpc-1.org",
		"ip_address": "45.65.0.1",
		"isp" : {
			"id": 3,
			"name": "移动"
		},
		"province": {
			"id": 2,
			"name": "山西"
		},
		"city": {
			"id": -1,
			"name": "<UNDEFINED>"
		},
		"name_tag": {
			"id": -1,
			"value": "\u003cUNDEFINED\u003e"
		},
		"group_tags": [
			9081,
			9082
		],
		"status": true,
		"comment": null,
		"last_heartbeat_time": 2347726123,
		"num_of_enabled_pingtasks": 3
	}
	`,
	`
	[
	{
	  "id": 630001,
	  "host": "1.2.3.4",
	  "isp_id": 1,
	  "isp_name": "北京三信时代",
	  "province_id": 4,
	  "province_name": "北京",
	  "ct_id": 1,
	  "ct_name": "北京市",
	  "nt_id": -1,
	  "nt_value": "\u003cUNDEFINED\u003e",
	  "gt_ids": [
	     9081,
	     9082
	  ]
	},
	{
	  "id": 630002,
	  "host": "1.2.3.5",
	  "isp_id": 2,
	  "isp_name": "教育网",
	  "province_id": 5,
	  "province_name": "辽宁",
	  "ct_id": 280,
	  "ct_name": "本溪市",
	  "nt_id": 9901,
	  "nt_value": "tag-1",
	  "gt_ids": [
	     9081
	  ]
	},
	{
	  "id": 630003,
	  "host": "1.2.3.6",
	  "isp_id": 3,
	  "isp_name": "移动",
	  "province_id": 5,
	  "province_name": "辽宁",
	  "ct_id": 285,
	  "ct_name": "葫芦岛市",
	  "nt_id": -1,
	  "nt_value": "\u003cUNDEFINED\u003e",
	  "gt_ids": [
	     9082
	  ]
	}
	]
	`,
	`
	{
		"id" : 405002,
		"name": "ag-name-1",
		"connection_id": "ag-rpc-2",
		"hostname": "rpc-2.org",
		"ip_address": "45.65.0.2",
		"isp" : {
			"id": 3,
			"name": "移动"
		},
		"province": {
			"id": 2,
			"name": "山西"
		},
		"city": {
			"id": -1,
			"name": "<UNDEFINED>"
		},
		"name_tag": {
			"id": -1,
			"value": "\u003cUNDEFINED\u003e"
		},
		"group_tags": [
			9081,
			9082
		],
		"status": false,
		"comment": null,
		"last_heartbeat_time": 2347726123,
		"num_of_enabled_pingtasks": 3
	}
	`,
	`
	{
		"id" : 405001,
		"name": "ag-name-1",
		"connection_id": "ag-rpc-1",
		"hostname": "rpc-1.org",
		"ip_address": "45.65.0.1",
		"isp" : {
			"id": 3,
			"name": "移动"
		},
		"province": {
			"id": 2,
			"name": "山西"
		},
		"city": {
			"id": -1,
			"name": "<UNDEFINED>"
		},
		"name_tag": {
			"id": -1,
			"value": "\u003cUNDEFINED\u003e"
		},
		"group_tags": [
			9081,
			9082
		],
		"status": false,
		"comment": null,
		"last_heartbeat_time": 2347726123,
		"num_of_enabled_pingtasks": 3
	}
	`,
}

// Tests the refreshing and retrieving list of targets for NQM agent
func (suite *TestNqmAgentSuite) TestTask(c *C) {
	if !testJsonRpc.HasJsonRpcServ(c) {
		return
	}

	testCases := []*struct {
		req               model.NqmTaskRequest
		needPing          bool
		expectedTargetIds []int32
	}{
		{
			model.NqmTaskRequest{
				ConnectionId: "ag-rpc-1",
				Hostname:     "rpc-1.org",
				IpAddress:    "45.65.0.1",
			},
			true, []int32{630001, 630002, 630003},
		},
		{
			model.NqmTaskRequest{
				ConnectionId: "ag-rpc-2",
				Hostname:     "rpc-2.org",
				IpAddress:    "45.65.0.2",
			},
			false, []int32{},
		},
		{ // The period is not elapsed yet
			model.NqmTaskRequest{
				ConnectionId: "ag-rpc-1",
				Hostname:     "rpc-1.org",
				IpAddress:    "45.65.0.1",
			},
			false, []int32{},
		},
	}

	for i, testCase := range testCases {
		ts = httptest.NewUnstartedServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, tsRsps[tsCnt])
				tsCnt++
			}),
		)
		l, err := net.Listen("tcp", MOCK_URL)
		fmt.Println(l, err)
		ts.Listener = l
		ts.Start()

		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		var resp model.NqmTaskResponse
		testJsonRpc.OpenClient(c, func(jsonRpcClient *rpc.Client) {
			err := jsonRpcClient.Call(
				"NqmAgent.Task", testCase.req, &resp,
			)

			c.Assert(err, IsNil)
		})

		connectionId := testCase.req.ConnectionId
		c.Logf("[%s] Response Need Ping: %v", connectionId, resp.NeedPing)
		c.Logf("[%s] Response Agent: %#v", connectionId, resp.Agent)
		c.Logf("[%s] Response Measurements: %#v", connectionId, resp.Measurements)

		/**
		 * The case with no needed of PING
		 */
		if !testCase.needPing {
			c.Assert(resp.NeedPing, Equals, false, comment)
			c.Assert(resp.Agent, IsNil)
			ts.Close()
			continue
		}
		// :~)

		c.Assert(resp.NeedPing, Equals, true, comment)
		c.Assert(
			resp.Agent, DeepEquals,
			&model.NqmAgent{
				Id:    405001,
				Name:  "ag-name-1",
				IspId: 3, IspName: "移动",
				ProvinceId: 2, ProvinceName: "山西",
				CityId: -1, CityName: "<UNDEFINED>",
				NameTagId:   -1,
				GroupTagIds: []int32{9081, 9082},
			},
		)

		c.Assert(len(resp.Targets), Equals, 3)
		c.Assert(resp.Measurements["fping"].Command[0], Equals, "fping")

		c.Assert(resp.Targets, HasLen, len(testCase.expectedTargetIds), comment)
		for i, targetId := range testCase.expectedTargetIds {
			c.Logf("\tTarget: %#v", resp.Targets[i])
			c.Assert(int32(resp.Targets[i].Id), Equals, targetId, comment)
		}
		ts.Close()
	}
}
