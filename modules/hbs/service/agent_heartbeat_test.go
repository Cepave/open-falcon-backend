package service

import (
	"strconv"
	"time"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeHeartbeat struct {
	rowsAffectedCnt int
	alwaysFail bool
}
func (a *fakeHeartbeat) calling(agents []*model.AgentHeartbeat) (int64, int64) {
	a.rowsAffectedCnt += len(agents)

	if a.alwaysFail {
		return int64(len(agents)), 0
	} else {
		return 0, int64(len(agents))
	}
}

var _ = Describe("Test put() of AgentHeartbeat service", func() {
	var (
		agentHeartbeatService *AgentHeartbeatService
		heartbeatImpl *fakeHeartbeat
	)

	BeforeEach(func() {
		agentHeartbeatService = NewAgentHeartbeatService(
			&commonQueue.Config{ Num: 16 },
		)

		heartbeatImpl = &fakeHeartbeat{ alwaysFail: false }
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

func generateRandomHeartbeat() *commonModel.AgentReportRequest {
	return &commonModel.AgentReportRequest{
		Hostname:      "agentHeartbeatService-",
		IP:            "127.0.0.",
		AgentVersion:  "0.0.",
		PluginVersion: "12345abcde",
	}
}

var _ = Describe("Original Test AgentHeartbeat service", func() {
	var (
		caseNumber      int   = 0
		putNumber       int64 = 0
		data            *commonModel.AgentReportRequest
		heartbeatConfig *commonQueue.Config = &commonQueue.Config{
			Num: 16,
			Dur: 100 * time.Millisecond,
		}
		sleepTime time.Duration = time.Second / 2
	)
	agentHeartbeatService := NewAgentHeartbeatService(heartbeatConfig)
	agentHeartbeatService.heartbeatCall = func(agents []*model.AgentHeartbeat) (rowsAffectedCnt int64, agentsDroppedCnt int64) {
		rowsAffectedCnt = int64(len(agents))
		return rowsAffectedCnt, 0
	}

	BeforeEach(func() {
		sampleNumber := strconv.Itoa(caseNumber)
		data = &commonModel.AgentReportRequest{
			Hostname:      "agentHeartbeatService-" + sampleNumber,
			IP:            "127.0.0." + sampleNumber,
			AgentVersion:  "0.0." + sampleNumber,
			PluginVersion: "12345abcde" + sampleNumber,
		}
		caseNumber++
	})

	AfterEach(func() {
		Expect(agentHeartbeatService.CumulativeAgentsPut()).To(Equal(putNumber))
		Expect(agentHeartbeatService.CumulativeRowsAffected()).To(Equal(putNumber))
	})

	Context("before service started", func() {
		It("data should not be processed by Put() method", func() {
			agentHeartbeatService.Put(data)
		})
	})

	Context("after service started", func() {
		It("data should be processed by Put() method", func() {
			anotherData := &commonModel.AgentReportRequest{}
			*anotherData = *data
			agentHeartbeatService.Start()
			agentHeartbeatService.Put(data)
			agentHeartbeatService.Put(anotherData)
			putNumber += 2
			time.Sleep(sleepTime)
		})
	})

	Context("after starting service again", func() {
		It("data should be processed by Put() method normally", func() {
			agentHeartbeatService.Start()
			agentHeartbeatService.Put(data)
			putNumber++
			time.Sleep(sleepTime)
		})
	})

	Context("Stop service after putting elements(# > batch size)", func() {
		It("data should be flushed", func() {
			anotherData := &commonModel.AgentReportRequest{}
			for num := 1; num < heartbeatConfig.Num*2; num++ {
				*anotherData = *data
				agentHeartbeatService.Put(anotherData)
				putNumber++
			}
			agentHeartbeatService.Stop()
			time.Sleep(sleepTime)
		})
	})

	Context("after service stopped", func() {
		It("data should not be processed by Put() method", func() {
			agentHeartbeatService.Put(data)
		})
	})

	Context("after stopping service again", func() {
		It("data should not be processed by Put() method", func() {
			agentHeartbeatService.Stop()
			agentHeartbeatService.Put(data)
		})
	})
})
