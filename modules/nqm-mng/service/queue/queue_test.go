package queue

import (
	"flag"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	// /. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
)

func init() {
	flag.Parse()
}

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var _ = Describe("Start(): Start the queue service", ginkgoDb.NeedDb(func() {
	It("can't be put elements without calling Start() in advance", func() {
		testedQueue := New(&commonQueue.Config{Num: 1, Dur: 0})

		testedQueue.Put(&model.NqmAgentHeartbeatRequest{
			ConnectionId: "test1-hostname@1.2.3.4",
			Hostname:     "test1-hostname",
			IpAddress:    "1.2.3.4",
			Timestamp:    ojson.JsonTime(time.Now()),
		})
		Expect(testedQueue.Count()).To(Equal(uint64(0)))
		Expect(testedQueue.Len()).To(Equal(0))
	})

	It("can be put elements by calling Start() in advance", func() {
		testedQueue := New(&commonQueue.Config{Num: 1, Dur: 0})

		testedQueue.Start()

		testedQueue.Put(&model.NqmAgentHeartbeatRequest{
			ConnectionId: "test1-hostname@1.2.3.4",
			Hostname:     "test1-hostname",
			IpAddress:    "1.2.3.4",
			Timestamp:    ojson.JsonTime(time.Now()),
		})

		testedQueue.Put(&model.NqmAgentHeartbeatRequest{
			ConnectionId: "test2-hostname@1.2.3.4",
			Hostname:     "test2-hostname",
			IpAddress:    "1.2.3.4",
			Timestamp:    ojson.JsonTime(time.Now()),
		})

		testedQueue.Stop()
		Expect(testedQueue.Count()).To(Equal(uint64(2)))
	})
}))

var _ = Describe("Stop(): Stop the queue service", ginkgoDb.NeedDb(func() {
	It("can't be put elements after being stopped", func() {
		testedQueue := New(&commonQueue.Config{Num: 1, Dur: 0})

		testedQueue.Start()
		testedQueue.Stop()

		testedQueue.Put(&model.NqmAgentHeartbeatRequest{
			ConnectionId: "test1-hostname@1.2.3.4",
			Hostname:     "test1-hostname",
			IpAddress:    "1.2.3.4",
			Timestamp:    ojson.JsonTime(time.Now()),
		})
		Expect(testedQueue.Count()).To(Equal(uint64(0)))
		Expect(testedQueue.Len()).To(Equal(0))
	})

	It("doesn't flush elements until Stop() is called", func() {
		testedQueue := New(&commonQueue.Config{Num: 10, Dur: 1 * time.Second})

		testedQueue.Start()

		for i := 0; i < 999; i++ {
			testedQueue.Put(&model.NqmAgentHeartbeatRequest{
				ConnectionId: "test1-hostname@1.2.3.4",
				Hostname:     "test1-hostname",
				IpAddress:    "1.2.3.4",
				Timestamp:    ojson.JsonTime(time.Now()),
			})
		}
		Expect(testedQueue.Len()).NotTo(Equal(0))

		testedQueue.Stop()
		Expect(testedQueue.Count()).To(Equal(uint64(999)))
		Expect(testedQueue.Len()).To(Equal(0))
	})

	It("has no elements after Stop() is called", func() {
		testedQueue := New(&commonQueue.Config{Num: 10, Dur: 1 * time.Second})

		testedQueue.Start()

		go func() {
			for i := 0; i < 999; i++ {
				testedQueue.Put(&model.NqmAgentHeartbeatRequest{
					ConnectionId: "test1-hostname@1.2.3.4",
					Hostname:     "test1-hostname",
					IpAddress:    "1.2.3.4",
					Timestamp:    ojson.JsonTime(time.Now()),
				})
			}
		}()

		testedQueue.Stop()
		Expect(testedQueue.Len()).To(Equal(0))
	})
}))

var ginkgoDb = &dbTest.GinkgoDb{}

var _ = BeforeSuite(func() {
	rdb.DbFacade = ginkgoDb.InitDbFacade()
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(rdb.DbFacade)
})
