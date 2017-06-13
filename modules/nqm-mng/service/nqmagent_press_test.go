package service

import (
	"time"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb/test"
	"github.com/icrowley/fake"

	. "github.com/onsi/ginkgo"
)

var heartbeatRequests = make(map[int]*model.NqmAgentHeartbeatRequest)

var ginkgoDb = &dbTest.GinkgoDb{}

var numOfReqs int = 3000
var batchSize int = 0

var _ = Describe("Pressure Test", ginkgoDb.NeedDb(func() {
	BeforeEach(func() {
		if batchSize == 0 {
			return
		}
		rdb.DbFacade = ginkgoDb.InitDbFacade()
		for i := 0; i < numOfReqs; i++ {
			r := randomHeartbeatRequest()
			heartbeatRequests[i] = r
			insert(r)
		}
	})

	AfterEach(func() {
		if batchSize == 0 {
			return
		}
		inTx(
			`DELETE FROM nqm_agent`,
			`DELETE FROM host`,
			test.ResetAutoIncForNqmAgent,
			test.ResetAutoIncForHost,
		)
		ginkgoDb.ReleaseDbFacade(rdb.DbFacade)
	})

	Measure("performance of the batch size", func(b Benchmarker) {
		if batchSize == 0 {
			GinkgoT().Logf(" == Skip Pressure Testing ==")
			return
		}

		GinkgoT().Logf("Number of requests: %d, Batch size: %d", numOfReqs, batchSize)
		InitNqmHeartbeat(&commonQueue.Config{Num: batchSize, Dur: 1})
		NqmQueue.Start()
		b.Time("runtime", func() {
			for _, req := range heartbeatRequests {
				func(r *model.NqmAgentHeartbeatRequest) {
					NqmQueue.Put(req)
				}(req)
			}
			CloseNqmHeartbeat()
		})
	}, 3)

}))

func inTx(sql ...string) {
	rdb.DbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

func randomHeartbeatRequest() *model.NqmAgentHeartbeatRequest {
	hostname := fake.CharactersN(20)
	ipAddr := fake.IPv4()
	connID := hostname + "@" + ipAddr
	return &model.NqmAgentHeartbeatRequest{
		Hostname:     hostname,
		IpAddress:    ojson.NewIP(ipAddr),
		ConnectionId: connID,
		Timestamp:    ojson.JsonTime(time.Now()),
	}
}

func insert(r *model.NqmAgentHeartbeatRequest) {
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
