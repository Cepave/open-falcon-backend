package rpc

import (
	"net/rpc"
	"strconv"

	coModel "github.com/Cepave/open-falcon-backend/common/model"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("[Stress] Test Agent.ReportStatus in HBS", ginkgoJsonRpc.NeedJsonRpc(func() {
	var (
		numberOfGoRoutines int = 50
		pool               *agentPool
		routines           chan bool

		numberOfFakeAgents int = 0
		numberOfSamples    int = 3
	)

	checkSkipCondition := func() {
		if numberOfFakeAgents == 0 {
			Skip("Number of total request is 0. See numberOfFakeAgents in code")
		}
	}

	BeforeEach(func() {
		checkSkipCondition()

		pool = &agentPool{numberOfFakeAgents}
		routines = make(chan bool, numberOfGoRoutines)
		for i := 0; i < numberOfGoRoutines; i++ {
			routines <- true
		}
	})

	Measure("It should serve lots of rpc-clients efficiently", func(b Benchmarker) {
		checkSkipCondition()
		b.Time("runtime", func() {
			for i := 0; i < numberOfFakeAgents; i++ {
				request := pool.getNextRequest(i)
				var resp coModel.SimpleRpcResponse

				<-routines

				go func() {
					defer GinkgoRecover()

					defer func() {
						routines <- true
					}()
					ginkgoJsonRpc.OpenClient(func(jsonRpcClient *rpc.Client) {

						err := jsonRpcClient.Call(
							"Agent.ReportStatus", request, &resp,
						)

						if err != nil || resp.Code == 1 {
							GinkgoT().Errorf("[%s] Has error: %v", request.AgentVersion, err)
						} else {
							GinkgoT().Logf("[%s/%d] Success.", request.PluginVersion, numberOfFakeAgents)
						}
					})
				}()

			}

			for i := 0; i < numberOfGoRoutines; i++ {
				<-routines
			}
		})
	}, numberOfSamples)
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
