package rpc

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"sync"

	"github.com/Cepave/open-falcon-backend/common/model"

	testJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"

	. "gopkg.in/check.v1"
)

type NqmAgentStressSuite struct{}

var _ = Suite(&NqmAgentStressSuite{})

var numberOfRoutines = 16

// Tests the NQM agent HBS in stress condition
func (suite *NqmAgentStressSuite) TestTask(c *C) {
	idsPool := &connIdsPool{
		ids:   connectionIds,
		len:   len(connectionIds),
		index: 0,
		mutex: &sync.Mutex{},
	}

	if idsPool.len == 0 {
		c.Skip("Connection Id is empty. See \"idsPool.len\" in code")
		return
	}

	routines := make(chan bool, numberOfRoutines)
	for i := 0; i < numberOfRoutines; i++ {
		routines <- true
	}

	for i := 0; i < idsPool.len; i++ {
		connectionId := idsPool.getNextConnId()
		ipAddress := idsPool.getNextIpAddress()

		c.Logf("Random IPAddress[%s] for Connection id: [%s]", ipAddress, connectionId)

		request := &model.NqmTaskRequest{
			ConnectionId: connectionId,
			Hostname:     ipAddress,
			IpAddress:    ipAddress,
		}

		var resp model.NqmTaskResponse

		<-routines

		go testJsonRpc.OpenClient(c, func(jsonRpcClient *rpc.Client) {
			defer func() {
				routines <- true
			}()

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

	for i := 0; i < numberOfRoutines; i++ {
		<-routines
	}
}

type connIdsPool struct {
	ids   []string
	len   int
	index int
	mutex *sync.Mutex

	ip1 int
}

func (p *connIdsPool) getNextConnId() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	connectionId := p.ids[p.index]
	p.index++
	p.index = p.index % p.len

	return connectionId
}
func (p *connIdsPool) getNextIpAddress() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.ip1++
	p.ip1 = p.ip1%254 + 1

	return fmt.Sprintf(
		"%d.%d.%d.%d",
		p.ip1,
		rand.Int31n(255)+1,
		rand.Int31n(255)+1,
		rand.Int31n(255)+1,
	)
}

var connectionIds = []string{}
