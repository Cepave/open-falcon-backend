package service

import (
	"fmt"
	"math/rand"
	"time"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/icrowley/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type dbNqmHeartbeatCapture int

func (db *dbNqmHeartbeatCapture) updator(agents []*model.NqmAgentHeartbeatRequest) {
	v := int(*db)
	v += len(agents)

	*db = dbNqmHeartbeatCapture(v)
}
func (db *dbNqmHeartbeatCapture) getNumber() int {
	return int(*db)
}

func mockNqmHeartbeatOnDb(agents []*model.NqmAgentHeartbeatRequest) {}

var _ = Describe("Tests Put() function", func() {
	var testedService *nqmAgentUpdateService
	var numberOfSampleHeartbeats = 4

	BeforeEach(func() {
		testedService = newNqmAgentUpdateServiceForTesting(&commonQueue.Config{Num: numberOfSampleHeartbeats, Dur: 0}, mockNqmHeartbeatOnDb)
	})
	AfterEach(func() {
		testedService.Stop()
	})

	It("Running service(effective operation)", func() {
		testedService.Start()
		putRandomHeartbeat(testedService, numberOfSampleHeartbeats)

		eventuallyWithTimeout(
			func() uint64 {
				count := testedService.ConsumedCount()
				GinkgoT().Logf("Current consumed count: %d", count)
				return count
			},
			2*time.Second,
		).Should(Equal(uint64(numberOfSampleHeartbeats)))
	})

	It("Stopped service(un-effective operation)", func() {
		putRandomHeartbeat(testedService, numberOfSampleHeartbeats)
		assertNumbersOfService(testedService, 0, 0)
	})
})
var _ = Describe("Tests syncToDatabase()", func() {
	var testedService *nqmAgentUpdateService
	var numberOfSampleHeartbeats = 6
	var dbCallingCapture *dbNqmHeartbeatCapture

	BeforeEach(func() {
		dbCallingCapture = new(dbNqmHeartbeatCapture)

		testedService = newNqmAgentUpdateServiceForTesting(&commonQueue.Config{Num: 3, Dur: 0}, dbCallingCapture.updator)
		testedService.running = true
		putRandomHeartbeat(testedService, numberOfSampleHeartbeats)
	})

	It("In normal mode(as running service)", func() {
		testedService.syncToDatabase(_DRAIN)
		Expect(testedService.ConsumedCount()).To(Equal(uint64(3)))
		Expect(testedService.PendingLen()).To(Equal(3))
		Expect(dbCallingCapture.getNumber()).To(Equal(3))
	})

	It("In flush mode(as stopping service)", func() {
		testedService.syncToDatabase(_FLUSH)
		Expect(testedService.ConsumedCount()).To(Equal(uint64(6)))
		Expect(testedService.PendingLen()).To(Equal(0))
		Expect(dbCallingCapture.getNumber()).To(Equal(6))
	})
})
var _ = Describe("Tests running flag", func() {
	var testedService *nqmAgentUpdateService

	BeforeEach(func() {
		testedService = newNqmAgentUpdateServiceForTesting(&commonQueue.Config{Num: 3, Dur: 0}, mockNqmHeartbeatOnDb)
	})
	AfterEach(func() {
		testedService.Stop()
	})

	It("By Start()", func() {
		Expect(testedService.running).To(BeFalse())
		testedService.Start()
		Expect(testedService.running).To(BeTrue())
	})

	It("By Stop()", func() {
		testedService.Start()
		Expect(testedService.running).To(BeTrue())
		testedService.Stop()
		Expect(testedService.running).To(BeFalse())
	})
})

var _ = Describe("Tests functions of service on full lifecycle", func() {
	numberOfHeartbeat := 512
	var dbCallingCapture = new(dbNqmHeartbeatCapture)
	testedService := newNqmAgentUpdateServiceForTesting(&commonQueue.Config{Num: 64, Dur: 0}, dbCallingCapture.updator)

	BeforeEach(func() {
		testedService.Start()
		GinkgoT().Logf("Puts %d heartbeats", numberOfHeartbeat)
		putRandomHeartbeat(testedService, numberOfHeartbeat)
	})
	AfterEach(func() {
		testedService.Stop()
	})

	It("Tests the consumer number", func() {
		eventuallyWithTimeout(
			func() int { return int(testedService.ConsumedCount()) },
			3*time.Second,
		).Should(Equal(numberOfHeartbeat))

		Expect(testedService.PendingLen()).To(Equal(0))
		Expect(dbCallingCapture.getNumber()).To(Equal(numberOfHeartbeat))
	})
})

func assertNumbersOfService(testedService *nqmAgentUpdateService, expectedConsumedCount uint64, expectedPendingLen uint64) {
	consumedCount := testedService.ConsumedCount()
	pendingLen := testedService.PendingLen()

	GinkgoT().Logf("Consumed count: %d. Pending len: %d.", consumedCount, pendingLen)

	Ω(consumedCount).Should(Equal(consumedCount))
	Ω(pendingLen).Should(Equal(pendingLen))
}

func eventuallyWithTimeout(valueGetter interface{}, timeout time.Duration) GomegaAsyncAssertion {
	return Eventually(
		valueGetter, timeout, timeout/8,
	)
}

func putRandomHeartbeat(srv *nqmAgentUpdateService, number int) {
	for i := 0; i < number; i++ {
		srv.Put(buildRandomRequest())
	}
}
func buildRandomRequest() *model.NqmAgentHeartbeatRequest {
	hostname := fake.UserName()
	ip := fake.IPv4()

	// 2012-04-04T00:00:00Z
	var startTime int64 = 1333497600
	// 2013-04-04T00:00:00Z
	var endTime int64 = 1365033600

	return &model.NqmAgentHeartbeatRequest{
		ConnectionId: fmt.Sprintf("%s@%s", hostname, ip),
		Hostname:     hostname,
		IpAddress:    ojson.NewIP(ip),
		Timestamp:    ojson.JsonTime(time.Unix(startTime+rand.Int63n(endTime), 0)),
	}
}

func newNqmAgentUpdateServiceForTesting(queueConfig *commonQueue.Config, dbUpdator func([]*model.NqmAgentHeartbeatRequest)) *nqmAgentUpdateService {
	testedService := newNqmAgentUpdateService(queueConfig)
	testedService.updateToDatabase = dbUpdator

	return testedService
}
