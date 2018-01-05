package service

import (
	"time"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/test"
	"github.com/icrowley/fake"

	. "github.com/onsi/ginkgo"
)

var numberOfSampleRequests = 0 // 2048
var batchSize = 64

var _ = Describe("Pressure Test", itSkip.PrependBeforeEach(func() {
	var heartbeatRequests map[int]*nqmModel.HeartbeatRequest

	BeforeEach(func() {
		heartbeatRequests = make(map[int]*nqmModel.HeartbeatRequest)

		if numberOfSampleRequests == 0 {
			Skip("The variable \"numberOfSampleRequests\" == 0, skip benchmarks")
		}
		for i := 0; i < numberOfSampleRequests; i++ {
			r := randomHeartbeatRequest()
			heartbeatRequests[i] = r
			insert(r)
		}
	})

	AfterEach(func() {
		inTx(
			`DELETE FROM nqm_agent`,
			`DELETE FROM host`,
			test.ResetAutoIncForNqmAgent,
			test.ResetAutoIncForHost,
		)
	})

	Measure("performance of the batch size", func(b Benchmarker) {
		GinkgoT().Logf("Number of requests: %d, Batch size: %d", numberOfSampleRequests, batchSize)
		InitNqmHeartbeat(&commonQueue.Config{Num: batchSize, Dur: 1})
		NqmQueue.Start()
		b.Time("runtime", func() {
			for _, req := range heartbeatRequests {
				reqObject := req
				NqmQueue.Put(reqObject)
			}
			CloseNqmHeartbeat()
		})

		GinkgoT().Logf("Count of consumed request: [%d]", NqmQueue.ConsumedCount())
	}, 3)

}))

func inTx(sql ...string) {
	rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

func randomHeartbeatRequest() *nqmModel.HeartbeatRequest {
	hostname := fake.CharactersN(20)
	ipAddr := fake.IPv4()
	connID := hostname + "@" + ipAddr
	return &nqmModel.HeartbeatRequest{
		Hostname:     hostname,
		IpAddress:    ojson.NewIP(ipAddr),
		ConnectionId: connID,
		Timestamp:    ojson.JsonTime(time.Now()),
	}
}

func insert(r *nqmModel.HeartbeatRequest) {
	rdb.DbFacade.SqlxDb.MustExec(`
			INSERT INTO host(hostname, ip, agent_version, plugin_version)
			VALUES(?, ?, '', '')
			ON DUPLICATE KEY UPDATE
				ip = VALUES(ip)
		`,
		r.Hostname,
		r.IpAddress.String(),
	)
	rdb.DbFacade.SqlxDb.MustExec(`
			INSERT INTO nqm_agent(ag_connection_id, ag_hostname, ag_ip_address, ag_last_heartbeat, ag_hs_id)
			SELECT ?, ?, ?, FROM_UNIXTIME(?), id
			FROM host
			WHERE hostname = ?
			ON DUPLICATE KEY UPDATE
				ag_hostname = VALUES(ag_hostname),
				ag_ip_address = VALUES(ag_ip_address),
				ag_last_heartbeat = VALUES(ag_last_heartbeat)
		`,
		r.ConnectionId,
		r.Hostname,
		r.IpAddress,
		r.Timestamp,
		r.Hostname,
	)
}
