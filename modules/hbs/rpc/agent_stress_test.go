package rpc

import (
	"net/rpc"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	coModel "github.com/open-falcon/common/model"
)

var _ = Describe("[Stress] Test Agent.ReportStatus in HBS", ginkgoJsonRpc.NeedJsonRpc(func() {
	var (
		numberOfGoRoutines int = 20
		pool               *agentPool
		routines           chan bool

		numberOfFakeAgent    int = 500
		numberOfTotalRequest int = 2000
		stepOfHeartbeat      int = 30
	)

	BeforeEach(func() {
		pool = &agentPool{numberOfFakeAgent}
		routines = make(chan bool, numberOfGoRoutines)
		for i := 0; i < numberOfGoRoutines; i++ {
			routines <- true
		}
	})

	It("Should run without error msg", func() {
		for i := 0; i < numberOfTotalRequest; i++ {
			if i > 0 && (i%numberOfFakeAgent == 0) {
				time.Sleep(time.Duration(stepOfHeartbeat) * time.Second)
			}
			request := pool.getNextRequest(i)
			var resp coModel.SimpleRpcResponse

			<-routines

			go ginkgoJsonRpc.OpenClient(func(jsonRpcClient *rpc.Client) {
				defer func() {
					routines <- true
				}()

				err := jsonRpcClient.Call(
					"Agent.ReportStatus", request, &resp,
				)

				if err != nil {
					GinkgoT().Errorf("[%s] Has error: %v", request.AgentVersion, err)
				} else {
					GinkgoT().Logf("[%s/%d] Success.", request.AgentVersion, numberOfTotalRequest)
				}
			})
		}

		for i := 0; i < numberOfGoRoutines; i++ {
			<-routines
		}
	})
}))

type agentPool struct {
	ringSize int
}

func (ap *agentPool) getNextRequest(requestNumber int) *coModel.AgentReportRequest {
	agentIdx := strconv.Itoa(requestNumber % ap.ringSize)
	return &coModel.AgentReportRequest{
		Hostname:     "stress-reportstatus-" + agentIdx,
		IP:           "127.0.0.56",
		AgentVersion: agentIdx,
	}
}
