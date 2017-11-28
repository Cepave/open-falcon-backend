package cmdb

import (
	_ "net"
	_ "reflect"
	"testing"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	cmdbModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestCmdbSuite struct{}

var _ = Suite(&TestCmdbSuite{})

func (suite *TestCmdbSuite) TestApi2tuple(c *C) {
	testCase := []cmdbModel.SyncHost{
		cmdbModel.SyncHost{
			Activate: 0,
			Name:     "cmdb-test-a",
			IP:       "69.69.69.1"},
		cmdbModel.SyncHost{
			Activate: 0,
			Name:     "cmdb-test-b",
			IP:       "69.69.69.2"},
		cmdbModel.SyncHost{
			Activate: 1,
			Name:     "cmdb-test-c",
			IP:       "69.69.69.3"},
		cmdbModel.SyncHost{
			Activate: 1,
			Name:     "cmdb-test-d",
			IP:       "69.69.69.4"},
	}
	spec := []hostTuple{
		hostTuple{
			Hostname:       "cmdb-test-a",
			Ip:             "69.69.69.1",
			Maintain_begin: 946684800,
			Maintain_end:   4292329420},
		hostTuple{
			Hostname:       "cmdb-test-b",
			Ip:             "69.69.69.2",
			Maintain_begin: 946684800,
			Maintain_end:   4292329420},
		hostTuple{
			Hostname:       "cmdb-test-c",
			Ip:             "69.69.69.3",
			Maintain_begin: 0,
			Maintain_end:   0},
		hostTuple{
			Hostname:       "cmdb-test-d",
			Ip:             "69.69.69.4",
			Maintain_begin: 0,
			Maintain_end:   0},
	}
	ocheck.LogTestCase(c, testCase)
	obtain := api2tuple(testCase)
	c.Assert(obtain, DeepEquals, spec)
}

type groupTuple struct {
	Name      string
	Creator   string
	Come_from int8
}

type relTuple struct {
	GrpID  int
	HostID int
}

func (suite *TestCmdbSuite) TestSyncRel(c *C) {
	testCase := map[string][]string{
		"cmdb-test-grp-b": []string{"cmdb-test-a", "cmdb-test-b"},
		"cmdb-test-grp-c": []string{"cmdb-test-a", "cmdb-test-b"},
	}
	txProcessor := &syncRelTx{
		relations: testCase,
	}
	ocheck.LogTestCase(c, testCase)
	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	//
	spec := []relTuple{
		relTuple{
			GrpID:  10,
			HostID: 1},
		relTuple{
			GrpID:  10,
			HostID: 2},
		relTuple{
			GrpID:  20,
			HostID: 1},
		relTuple{
			GrpID:  20,
			HostID: 2},
		relTuple{
			GrpID:  30,
			HostID: 1},
		relTuple{
			GrpID:  30,
			HostID: 2},
	}
	// compare table grp_host after sync
	rows, err := DbFacade.SqlxDb.Query("SELECT grp_id, host_id FROM grp_host order by grp_id, host_id")
	c.Assert(err, IsNil)
	index := 0
	for rows.Next() {
		var gid int
		var hid int
		err := rows.Scan(&gid, &hid)
		c.Assert(err, IsNil)
		c.Assert(gid, Equals, spec[index].GrpID)
		c.Assert(hid, Equals, spec[index].HostID)
		index = index + 1
	}
}

func (suite *TestCmdbSuite) TestSyncGrp(c *C) {
	testCase := []cmdbModel.SyncHostGroup{
		cmdbModel.SyncHostGroup{
			Name:    "cmdb-test-grp-a",
			Creator: "root"},
		cmdbModel.SyncHostGroup{
			Name:    "cmdb-test-grp-b",
			Creator: "root"},
		cmdbModel.SyncHostGroup{
			Name:    "cmdb-test-grp-c",
			Creator: "root"},
		cmdbModel.SyncHostGroup{
			Name:    "cmdb-test-grp-d",
			Creator: "root"},
	}
	txProcessor := &syncHostGroupTx{
		groups: testCase,
	}
	ocheck.LogTestCase(c, testCase)
	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	//
	spec := []groupTuple{
		groupTuple{
			Name:      "cmdb-test-grp-a",
			Creator:   "root",
			Come_from: 1},
		groupTuple{
			Name:      "cmdb-test-grp-e",
			Creator:   "default",
			Come_from: 0},
		groupTuple{
			Name:      "cmdb-test-grp-f",
			Creator:   "default",
			Come_from: 0},
		groupTuple{
			Name:      "cmdb-test-grp-b",
			Creator:   "root",
			Come_from: 1},
		groupTuple{
			Name:      "cmdb-test-grp-c",
			Creator:   "root",
			Come_from: 1},
		groupTuple{
			Name:      "cmdb-test-grp-d",
			Creator:   "root",
			Come_from: 1},
	}
	// compare host table after sync
	rows, err := DbFacade.SqlxDb.Query("SELECT grp_name, create_user, come_from FROM grp")
	c.Assert(err, IsNil)
	index := 0
	for rows.Next() {
		var name string
		var creator string
		var from int8
		err := rows.Scan(&name, &creator, &from)
		c.Assert(err, IsNil)
		c.Assert(name, Equals, spec[index].Name)
		c.Assert(creator, Equals, spec[index].Creator)
		c.Assert(from, Equals, spec[index].Come_from)
		index = index + 1
	}
}

func (suite *TestCmdbSuite) TestSyncHost(c *C) {
	testCase := []cmdbModel.SyncHost{
		cmdbModel.SyncHost{
			Activate: 0,
			Name:     "cmdb-test-a",
			IP:       "69.69.69.1"},
		cmdbModel.SyncHost{
			Activate: 0,
			Name:     "cmdb-test-b",
			IP:       "69.69.69.2"},
		cmdbModel.SyncHost{
			Activate: 1,
			Name:     "cmdb-test-c",
			IP:       "69.69.69.3"},
		cmdbModel.SyncHost{
			Activate: 1,
			Name:     "cmdb-test-d",
			IP:       "69.69.69.4"},
	}
	txProcessor := &syncHostTx{
		hosts: api2tuple(testCase),
	}
	ocheck.LogTestCase(c, testCase)
	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	//
	spec := []hostTuple{
		hostTuple{
			Hostname:       "cmdb-test-a",
			Ip:             "69.69.69.1",
			Maintain_begin: 946684800,
			Maintain_end:   4292329420},
		hostTuple{
			Hostname:       "cmdb-test-e",
			Ip:             "69.69.69.5",
			Maintain_begin: 0,
			Maintain_end:   0},
		hostTuple{
			Hostname:       "cmdb-test-f",
			Ip:             "69.69.69.6",
			Maintain_begin: 0,
			Maintain_end:   0},
		hostTuple{
			Hostname:       "cmdb-test-b",
			Ip:             "69.69.69.2",
			Maintain_begin: 946684800,
			Maintain_end:   4292329420},
		hostTuple{
			Hostname:       "cmdb-test-c",
			Ip:             "69.69.69.3",
			Maintain_begin: 0,
			Maintain_end:   0},
		hostTuple{
			Hostname:       "cmdb-test-d",
			Ip:             "69.69.69.4",
			Maintain_begin: 0,
			Maintain_end:   0},
	}
	// compare table host after sync
	rows, err := DbFacade.SqlxDb.Query("SELECT hostname, ip, maintain_begin, maintain_end FROM host")
	c.Assert(err, IsNil)
	index := 0
	for rows.Next() {
		var name string
		var ip string
		var m_b uint32
		var m_e uint32
		err := rows.Scan(&name, &ip, &m_b, &m_e)
		c.Assert(err, IsNil)
		c.Assert(name, Equals, spec[index].Hostname)
		c.Assert(ip, Equals, spec[index].Ip)
		c.Assert(m_b, Equals, spec[index].Maintain_begin)
		c.Assert(m_e, Equals, spec[index].Maintain_end)
		index = index + 1
	}
}

func (s *TestCmdbSuite) SetUpTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestCmdbSuite.TestSyncHost":
		inTx(
			`
			INSERT INTO host (hostname, ip, maintain_begin, maintain_end)
			VALUES ("cmdb-test-a","69.69.69.99",946684800,4292329420),
			       ("cmdb-test-e","69.69.69.5", 0, 0),
				   ("cmdb-test-f","69.69.69.6", 0, 0)
			`)
	case "TestCmdbSuite.TestSyncGrp":
		inTx(
			`INSERT INTO grp (grp_name, create_user, come_from)
			 VALUES ("cmdb-test-grp-a", "default", 0),
					("cmdb-test-grp-e", "default", 0),
					("cmdb-test-grp-f", "default", 0)
			`)
	case "TestCmdbSuite.TestSyncRel":
		inTx(
			`
			INSERT INTO host (id, hostname)
			VALUES (1, "cmdb-test-a"),
			       (2, "cmdb-test-b"),
				   (3, "cmdb-test-c")
			`,
			`
			INSERT INTO grp (id, grp_name, come_from)
			VALUES (10, "cmdb-test-grp-a", 0),
				   (20, "cmdb-test-grp-b", 1),
				   (30, "cmdb-test-grp-c", 1)
			`,
			`
			INSERT INTO grp_host(grp_id, host_id)
			VALUES (10, 1),
			       (10, 2),
				   (20, 2),
				   (30, 3)
			`)
	}
}

func (s *TestCmdbSuite) TearDownTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestCmdbSuite.TestSyncHost":
		inTx(
			`DELETE FROM host WHERE hostname LIKE "cmdb-test-%"`,
		)
	case "TestCmdbSuite.TestSyncGrp":
		inTx(
			`DELETE FROM grp WHERE grp_name LIKE "cmdb-test-grp-%"`,
		)
	case "TestCmdbSuite.TestSyncRel":
		inTx(
			`DELETE FROM grp_host WHERE grp_id IN (SELECT id FROM grp)`,
			`DELETE FROM grp_host WHERE host_id IN (SELECT id FROM host)`,
			`DELETE FROM host WHERE hostname LIKE "cmdb-test-%"`,
			`DELETE FROM grp WHERE grp_name LIKE "cmdb-test-grp-%"`,
		)
	}
}

func (s *TestCmdbSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}

func (s *TestCmdbSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
