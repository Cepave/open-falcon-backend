package service

import (
	"time"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/hbs/cache"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/icrowley/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dummyTime = time.Now().Unix()

type fakeHeartbeat struct {
	rowsAffectedCnt int
	alwaysSuccess   bool
}

func (a *fakeHeartbeat) calling(agents []*model.AgentHeartbeat) (int64, int64) {
	count := len(agents)
	a.rowsAffectedCnt += count

	if a.alwaysSuccess {
		return int64(count), 0
	} else {
		return 0, int64(count)
	}
}

func generateRandomRequest() *commonModel.AgentReportRequest {
	return &commonModel.AgentReportRequest{
		Hostname:     fake.DomainName(),
		IP:           fake.IPv4(),
		AgentVersion: fake.Digits(),
	}
}

func eventuallyWithTimeout(valueGetter interface{}, timeout time.Duration) GomegaAsyncAssertion {
	return Eventually(
		valueGetter, timeout, timeout/8,
	)
}

var _ = Describe("Test the behavior of AgentHeartbeat service", func() {
	var (
		agentHeartbeatService *AgentHeartbeatService
		heartbeatImpl         *fakeHeartbeat
	)

	BeforeEach(func() {
		agentHeartbeatService = NewAgentHeartbeatService(
			&commonQueue.Config{Num: 16, Dur: 100 * time.Millisecond},
		)
		heartbeatImpl = &fakeHeartbeat{alwaysSuccess: true}
		agentHeartbeatService.heartbeatCall = heartbeatImpl.calling
	})

	It("should work normally during its life cycle", func() {
		batchSize := agentHeartbeatService.qConfig.Num
		dataNum := batchSize * 2

		agentHeartbeatService.Start()
		for i := 0; i < dataNum; i++ {
			agentHeartbeatService.Put(generateRandomRequest(), dummyTime)
		}
		agentHeartbeatService.Stop()

		Expect(agentHeartbeatService.CumulativeAgentsPut()).To(Equal(int64(dataNum)))
		eventuallyWithTimeout(func() int {
			return heartbeatImpl.rowsAffectedCnt
		}, time.Second).Should(Equal(dataNum))
		eventuallyWithTimeout(func() int {
			return agentHeartbeatService.CurrentSize()
		}, time.Second).Should(BeZero())
	})
})

var _ = Describe("Test Put() of AgentHeartbeat service", func() {
	var (
		agentHeartbeatService *AgentHeartbeatService
		heartbeatImpl         *fakeHeartbeat
		now                   int64
	)

	BeforeEach(func() {
		agentHeartbeatService = NewAgentHeartbeatService(
			&commonQueue.Config{Num: 16},
		)
		heartbeatImpl = &fakeHeartbeat{alwaysSuccess: true}
		agentHeartbeatService.heartbeatCall = heartbeatImpl.calling
		now = time.Now().Unix()
	})

	AfterEach(func() {
		cache.Agents = cache.NewSafeAgents()
	})

	Context("when service is not running", func() {
		It("should not add data", func() {
			data := generateRandomRequest()
			agentHeartbeatService.Put(data, now)
			_, ok := cache.Agents.Get(data.Hostname)

			Expect(ok).To(BeFalse())
			Expect(agentHeartbeatService.CumulativeAgentsPut()).To(BeZero())
			Expect(agentHeartbeatService.CurrentSize()).To(BeZero())
		})
	})

	Context("when service is running", func() {
		It("should add data", func() {
			data := generateRandomRequest()
			agentHeartbeatService.running = true
			agentHeartbeatService.Put(data, now)
			val, ok := cache.Agents.Get(data.Hostname)

			Expect(ok).To(BeTrue())
			Expect(val.ReportRequest.Hostname).To(Equal(data.Hostname))
			Expect(val.LastUpdate).To(Equal(now))
			Expect(agentHeartbeatService.CumulativeAgentsPut()).To(Equal(int64(1)))
			Expect(agentHeartbeatService.CurrentSize()).To(Equal(1))
		})
	})
})

var _ = Describe("Test Start() of AgentHeartbeat service", func() {
	var (
		agentHeartbeatService *AgentHeartbeatService
	)

	BeforeEach(func() {
		agentHeartbeatService = NewAgentHeartbeatService(
			&commonQueue.Config{},
		)
	})

	AfterEach(func() {
		agentHeartbeatService.Stop()
	})

	Context("when service is stopped", func() {
		It("Start() should change the running status", func() {
			agentHeartbeatService.Start()
			Expect(agentHeartbeatService.running).To(BeTrue())
		})
	})

	Context("when service is started", func() {
		It("Start() should not change the running status", func() {
			agentHeartbeatService.running = true
			agentHeartbeatService.Start()
			Expect(agentHeartbeatService.running).To(BeTrue())
		})
	})
})

var _ = Describe("Test Stop() of AgentHeartbeat service", func() {
	var (
		agentHeartbeatService *AgentHeartbeatService
	)

	BeforeEach(func() {
		agentHeartbeatService = NewAgentHeartbeatService(
			&commonQueue.Config{},
		)
	})

	Context("when service is stopped", func() {
		It("Stop() should not change the running status", func() {
			agentHeartbeatService.Stop()
			Expect(agentHeartbeatService.running).To(BeFalse())
		})
	})

	Context("when service is started", func() {
		It("Stop() should change the running status", func() {
			agentHeartbeatService.running = true
			agentHeartbeatService.Stop()
			Expect(agentHeartbeatService.running).To(BeFalse())
		})
	})
})

var _ = Describe("Test consumeHeartbeatQueue() of AgentHeartbeat service", func() {
	var (
		agentHeartbeatService *AgentHeartbeatService
		heartbeatImpl         *fakeHeartbeat
	)

	BeforeEach(func() {
		heartbeatImpl = &fakeHeartbeat{alwaysSuccess: true}
	})

	JustBeforeEach(func() {
		agentHeartbeatService = NewAgentHeartbeatService(
			&commonQueue.Config{Num: 16},
		)
		agentHeartbeatService.heartbeatCall = heartbeatImpl.calling
		agentHeartbeatService.running = true
	})

	Context("when success", func() {
		It("rowsAffectedCnt should be incremented normally", func() {
			agentHeartbeatService.Put(generateRandomRequest(), dummyTime)
			agentHeartbeatService.consumeHeartbeatQueue(false)

			Expect(heartbeatImpl.rowsAffectedCnt).To(Equal(1))
			Expect(agentHeartbeatService.CumulativeRowsAffected()).To(Equal(int64(1)))
			Expect(agentHeartbeatService.CumulativeAgentsDropped()).To(BeZero())
		})
	})

	Context("when failure", func() {
		BeforeEach(func() {
			heartbeatImpl = &fakeHeartbeat{alwaysSuccess: false}
		})

		It("agentsDroppedCnt should be incremented normally", func() {
			agentHeartbeatService.Put(generateRandomRequest(), dummyTime)
			agentHeartbeatService.consumeHeartbeatQueue(false)

			Expect(heartbeatImpl.rowsAffectedCnt).To(Equal(1))
			Expect(agentHeartbeatService.CumulativeRowsAffected()).To(BeZero())
			Expect(agentHeartbeatService.CumulativeAgentsDropped()).To(Equal(int64(1)))
		})
	})

	Context("when non-flushing mode", func() {
		It("should consume an amount of data = batch size", func() {
			batchSize := agentHeartbeatService.qConfig.Num
			dataNum := batchSize*2 - 1
			for i := 0; i < dataNum; i++ {
				agentHeartbeatService.Put(generateRandomRequest(), dummyTime)
			}
			agentHeartbeatService.consumeHeartbeatQueue(false)

			Expect(heartbeatImpl.rowsAffectedCnt).To(Equal(batchSize))
			Expect(agentHeartbeatService.CurrentSize()).To(Equal(batchSize - 1))
		})
	})

	Context("when flushing mode", func() {
		It("should flush data to 0", func() {
			dataNum := agentHeartbeatService.qConfig.Num*2 - 1
			for i := 0; i < dataNum; i++ {
				agentHeartbeatService.Put(generateRandomRequest(), dummyTime)
			}
			agentHeartbeatService.consumeHeartbeatQueue(true)

			Expect(heartbeatImpl.rowsAffectedCnt).To(Equal(dataNum))
			Expect(agentHeartbeatService.CurrentSize()).To(BeZero())
		})
	})
})
