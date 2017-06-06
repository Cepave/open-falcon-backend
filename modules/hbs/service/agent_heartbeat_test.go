package service

import (
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/icrowley/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeHeartbeat struct {
	rowsAffectedCnt int
	alwaysFail      bool
}

func (a *fakeHeartbeat) calling(agents []*model.AgentHeartbeat) (int64, int64) {
	a.rowsAffectedCnt += len(agents)

	if a.alwaysFail {
		return int64(len(agents)), 0
	} else {
		return 0, int64(len(agents))
	}
}

func generateRandomHeartbeat() *commonModel.AgentReportRequest {
	return &commonModel.AgentReportRequest{
		Hostname:     fake.DomainName(),
		IP:           fake.IPv4(),
		AgentVersion: fake.Digits(),
	}
}

var _ = Describe("Test put() of AgentHeartbeat service", func() {
	var (
		agentHeartbeatService *AgentHeartbeatService
		heartbeatImpl         *fakeHeartbeat
	)

	BeforeEach(func() {
		agentHeartbeatService = NewAgentHeartbeatService(
			&commonQueue.Config{Num: 16},
		)

		heartbeatImpl = &fakeHeartbeat{alwaysFail: false}
		agentHeartbeatService.heartbeatCall = heartbeatImpl.calling
	})
	AfterEach(func() {
		agentHeartbeatService.Stop()
	})

	Context("when service is not running", func() {
		It("should not add data", func() {
			agentHeartbeatService.Put(generateRandomHeartbeat())
			Expect(agentHeartbeatService.CumulativeAgentsPut()).To(Equal(int64(0)))
			Expect(agentHeartbeatService.CurrentSize()).To(Equal(0))
		})
	})

	Context("when service is running", func() {
		It("should add data", func() {
			agentHeartbeatService.Start()
			agentHeartbeatService.Put(generateRandomHeartbeat())

			Expect(agentHeartbeatService.CumulativeAgentsPut()).To(Equal(int64(1)))
		})
	})
})
