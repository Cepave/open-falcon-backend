package rpc

import (
	"net/rpc"
	"sync"

	"github.com/Cepave/open-falcon-backend/common/model"

	rd "github.com/Pallinder/go-randomdata"

	testJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"

	. "gopkg.in/check.v1"
)

type NqmAgentStressSuite struct{}

var _ = Suite(&NqmAgentStressSuite{})

var benchamrkTaskGoRoutines = 4
// Tests the NQM agent HBS in stress condition
func (suite *NqmAgentStressSuite) TestTask(c *C) {
	idsPool := &connIdsPool {
		ids: connectionIds,
		len: len(connectionIds),
		index: 0,
		mutex: &sync.Mutex{},
	}

	if idsPool.len == 0 {
		c.Skip("Connection Id is empty. See \"idsPool.len\" in code")
		return
	}

	for i := 0; i < idsPool.len; i++ {
		connectionId := idsPool.getNextConnId()
		ipAddress := rd.IpV4Address()

		request := &model.NqmTaskRequest{
			ConnectionId: connectionId,
			Hostname:     ipAddress,
			IpAddress:    ipAddress,
		}

		var resp model.NqmTaskResponse
		testJsonRpc.OpenClient(c, func(jsonRpcClient *rpc.Client) {
			err := jsonRpcClient.Call(
				"NqmAgent.Task", request, &resp,
			)

			if err != nil {
				c.Errorf("[%s] Has error: %v", connectionId, err)
			} else {
				c.Logf("[%s] Number of match targets: [%d].", connectionId, len(resp.Targets))
			}
		})
	}
}

type connIdsPool struct {
	ids []string
	len int
	index int
	mutex *sync.Mutex
}
func (p *connIdsPool) getNextConnId() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	connectionId := p.ids[p.index]
	p.index++
	p.index = p.index % p.len

	return connectionId
}

var connectionIds = []string{}
