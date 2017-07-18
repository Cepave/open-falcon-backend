package rpc

import (
	"net/rpc"
	"strconv"

	. "github.com/onsi/ginkgo"
	coModel "github.com/open-falcon/common/model"
)

var _ = Describe("[Stress] Test Agent.ReportStatus in HBS", ginkgoJsonRpc.NeedJsonRpc(func() {
	var (
		numberOfGoRoutines int = 50
		pool               *agentPool
		routines           chan bool

		numberOfFakeAgent int = 0
	)

	checkSkipCondition := func() {
		if numberOfFakeAgent == 0 {
			Skip("Number of total request is 0. See numberOfFakeAgent in code")
		}
	}

	BeforeEach(func() {
		checkSkipCondition()

		pool = &agentPool{numberOfFakeAgent}
		routines = make(chan bool, numberOfGoRoutines)
		for i := 0; i < numberOfGoRoutines; i++ {
			routines <- true
		}
	})

	Measure("It should serve lots rpc-clients efficiently", func(b Benchmarker) {
		checkSkipCondition()
		b.Time("runtime", func() {
			for i := 0; i < numberOfFakeAgent; i++ {
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

					if err != nil || resp.Code == 1 {
						GinkgoT().Errorf("[%s] Has error: %v", request.AgentVersion, err)
					} else {
						GinkgoT().Logf("[%s/%d] Success.", request.PluginVersion, numberOfFakeAgent)
					}
				})
			}

			for i := 0; i < numberOfGoRoutines; i++ {
				<-routines
			}
		})
	}, 3)
}))

type agentPool struct {
	ringSize int
}

func (ap *agentPool) getNextRequest(requestNumber int) *coModel.AgentReportRequest {
	agentIdx := strconv.Itoa((requestNumber % ap.ringSize) + 1)
	return &coModel.AgentReportRequest{
		Hostname:      "stress-reportstatus-" + agentIdx,
		IP:            "127.0.0.56",
		AgentVersion:  agentIdx,
		PluginVersion: strconv.Itoa(requestNumber + 1),
	}
}
