package service

import (
	"strconv"
	"time"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/dghubble/sling"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test AgentHeartbeat service", func() {
	var (
		caseNumber      int   = 0
		putNumber       int64 = 0
		data            *commonModel.AgentReportRequest
		heartbeatConfig *commonQueue.Config = &commonQueue.Config{
			Num: 16,
			Dur: 3,
		}
	)
	mysqlApiSling = sling.New().Base("dummyString")
	agentHeartbeatService := NewAgentHeartbeatService(heartbeatConfig)
	agentHeartbeatService.heartbeatCall = func(agents []*model.AgentHeartbeat, slingAPI *sling.Sling) (rowsAffectedCnt int64, agentsDroppedCnt int64) {
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
			time.Sleep(timeWaitForQueue)
		})
	})

	Context("after starting service again", func() {
		It("data should be processed by Put() method normally", func() {
			agentHeartbeatService.Start()
			agentHeartbeatService.Put(data)
			putNumber++
			time.Sleep(timeWaitForQueue + timeWaitForInput)
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
			time.Sleep(timeWaitForQueue + timeWaitForInput)
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
