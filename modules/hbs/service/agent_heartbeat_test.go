package service

import (
	"time"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonModelConfig "github.com/Cepave/open-falcon-backend/common/model/config"
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

func generateRandomHeartbeat() *commonModel.AgentReportRequest {
	return &commonModel.AgentReportRequest{
		Hostname:     fake.DomainName(),
		IP:           fake.IPv4(),
		AgentVersion: fake.Digits(),
	}
}

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

	Context("when service is not running", func() {
		It("should not add data", func() {
			data := generateRandomHeartbeat()
			agentHeartbeatService.Put(data, now)
			_, ok := cache.Agents.Get(data.Hostname)

			Expect(ok).To(Equal(false))
			Expect(agentHeartbeatService.CumulativeAgentsPut()).To(Equal(int64(0)))
			Expect(agentHeartbeatService.CurrentSize()).To(Equal(0))
		})
	})

	Context("when service is running", func() {
		It("should add data", func() {
			data := generateRandomHeartbeat()
			agentHeartbeatService.running = true
			agentHeartbeatService.Put(data, now)
			val, ok := cache.Agents.Get(data.Hostname)

			Expect(ok).To(Equal(true))
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
			Expect(agentHeartbeatService.running).To(Equal(true))
		})
	})

	Context("when service is started", func() {
		It("Start() should not change the running status", func() {
			agentHeartbeatService.running = true
			agentHeartbeatService.Start()
			Expect(agentHeartbeatService.running).To(Equal(true))
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
			Expect(agentHeartbeatService.running).To(Equal(false))
		})
	})

	Context("when service is started", func() {
		It("Stop() should change the running status", func() {
			agentHeartbeatService.running = true
			agentHeartbeatService.Stop()
			Expect(agentHeartbeatService.running).To(Equal(false))
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
			agentHeartbeatService.Put(generateRandomHeartbeat(), dummyTime)
			agentHeartbeatService.consumeHeartbeatQueue(false)

			Expect(heartbeatImpl.rowsAffectedCnt).To(Equal(1))
			Expect(agentHeartbeatService.CumulativeRowsAffected()).To(Equal(int64(1)))
			Expect(agentHeartbeatService.CumulativeAgentsDropped()).To(Equal(int64(0)))
		})
	})

	Context("when failure", func() {
		BeforeEach(func() {
			heartbeatImpl = &fakeHeartbeat{alwaysSuccess: false}
		})

		It("agentsDroppedCnt should be incremented normally", func() {
			agentHeartbeatService.Put(generateRandomHeartbeat(), dummyTime)
			agentHeartbeatService.consumeHeartbeatQueue(false)

			Expect(heartbeatImpl.rowsAffectedCnt).To(Equal(1))
			Expect(agentHeartbeatService.CumulativeRowsAffected()).To(Equal(int64(0)))
			Expect(agentHeartbeatService.CumulativeAgentsDropped()).To(Equal(int64(1)))
		})
	})

	Context("when non-flushing mode", func() {
		It("should consume an amount of data = batch size", func() {
			batchSize := agentHeartbeatService.qConfig.Num
			dataNum := batchSize*2 - 1
			for i := 0; i < dataNum; i++ {
				agentHeartbeatService.Put(generateRandomHeartbeat(), dummyTime)
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
				agentHeartbeatService.Put(generateRandomHeartbeat(), dummyTime)
			}
			agentHeartbeatService.consumeHeartbeatQueue(true)

			Expect(heartbeatImpl.rowsAffectedCnt).To(Equal(dataNum))
			Expect(agentHeartbeatService.CurrentSize()).To(Equal(0))
		})
	})
})

var _ = Describe("Test buildHeartbeatCall() of AgentHeartbeat service", func() {
	Context("when the call fail", func() {
		It("should return correct dropped amount", func() {
			dataNum := 3
			agents := make([]*model.AgentHeartbeat, dataNum)
			InitPackage(&commonModelConfig.MysqlApiConfig{Host: "dummyHost"}, "")
			rowsAffectedCnt, agentsDroppedCnt := agentHeartbeatCall(agents)

			Expect(rowsAffectedCnt).To(Equal(int64(0)))
			Expect(agentsDroppedCnt).To(Equal(int64(dataNum)))
		})
	})
})
