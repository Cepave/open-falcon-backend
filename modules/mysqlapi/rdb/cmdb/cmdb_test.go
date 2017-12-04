package cmdb

import (
	cmdbModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[CMDB] Test SyncHost()", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			`
			INSERT INTO host (hostname, ip, maintain_begin, maintain_end)
			VALUES ("cmdb-test-a","69.69.69.99",946684800,4292329420),
				   ("cmdb-test-e","69.69.69.5", 0, 0),
				   ("cmdb-test-f","69.69.69.6", 0, 0)
			`)
	})
	AfterEach(func() {
		inTx(
			`DELETE FROM host WHERE hostname LIKE "cmdb-test-%"`,
		)
	})
	Context("Sync testCase, Select and Check", func() {
		It("Initially insert 4 entries and then sync 4 entries", func() {
			testCase := []*cmdbModel.SyncHost{
				{
					Activate: 0,
					Name:     "cmdb-test-a",
					IP:       "69.69.69.1",
				},
				{
					Activate: 0,
					Name:     "cmdb-test-b",
					IP:       "69.69.69.2",
				},
				{
					Activate: 1,
					Name:     "cmdb-test-c",
					IP:       "69.69.69.3",
				},
				{
					Activate: 1,
					Name:     "cmdb-test-d",
					IP:       "69.69.69.4",
				},
			}
			txProcessor := &syncHostTx{
				hosts: api2tuple(testCase),
			}
			//
			spec := []*hostTuple{
				{
					Hostname:       "cmdb-test-a",
					Ip:             "69.69.69.1",
					Maintain_begin: MAINTAIN_PERIOD_BEGIN,
					Maintain_end:   MAINTAIN_PERIOD_END,
				},
				{
					Hostname:       "cmdb-test-e",
					Ip:             "69.69.69.5",
					Maintain_begin: 0,
					Maintain_end:   0,
				},
				{
					Hostname:       "cmdb-test-f",
					Ip:             "69.69.69.6",
					Maintain_begin: 0,
					Maintain_end:   0,
				},
				{
					Hostname:       "cmdb-test-b",
					Ip:             "69.69.69.2",
					Maintain_begin: MAINTAIN_PERIOD_BEGIN,
					Maintain_end:   MAINTAIN_PERIOD_END,
				},
				{
					Hostname:       "cmdb-test-c",
					Ip:             "69.69.69.3",
					Maintain_begin: 0,
					Maintain_end:   0,
				},
				{
					Hostname:       "cmdb-test-d",
					Ip:             "69.69.69.4",
					Maintain_begin: 0,
					Maintain_end:   0,
				},
			}
			obtain := []*hostTuple{}
			DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
			DbFacade.NewSqlxDbCtrl().Select(&obtain, "SELECT hostname, ip, maintain_begin, maintain_end from host")
			Expect(obtain).To(Equal(spec))
		})
	})
}))

var _ = Describe("[CMDB] syncHostGroupTx", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			`INSERT INTO grp (grp_name, create_user, come_from)
			 VALUES ("cmdb-test-grp-a", "default", 0),
					("cmdb-test-grp-e", "default", 0),
					("cmdb-test-grp-f", "default", 0)
			`)
		testCase := []*cmdbModel.SyncHostGroup{
			{
				Name:    "cmdb-test-grp-a",
				Creator: "root",
			},
			{
				Name:    "cmdb-test-grp-b",
				Creator: "root",
			},
			{
				Name:    "cmdb-test-grp-c",
				Creator: "root",
			},
			{
				Name:    "cmdb-test-grp-d",
				Creator: "root",
			},
		}
		txProcessor := &syncHostGroupTx{
			groups: testCase,
		}
		DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	})
	AfterEach(func() {
		inTx(
			`DELETE FROM grp WHERE grp_name LIKE "cmdb-test-grp-%"`,
		)
	})
	Context("Select count of grp", func() {
		It("count should be 6", func() {
			var count int
			DbFacade.NewSqlxDbCtrl().Get(&count, "SELECT count(*) from grp")
			Expect(count).To(Equal(6))
		})
	})
	Context("With select come_from = 1", func() {
		It("Name should be cmdb-test-grp-[a-d]", func() {
			var names []string
			DbFacade.NewSqlxDbCtrl().Select(&names, "SELECT grp_name FROM grp where come_from = 1 order by grp_name")
			Expect(names).To(Equal([]string{"cmdb-test-grp-a", "cmdb-test-grp-b", "cmdb-test-grp-c", "cmdb-test-grp-d"}))
		})
	})
	Context("With select come_from = 1", func() {
		It("Creator should be root", func() {
			var name string
			DbFacade.NewSqlxDbCtrl().Get(&name, "SELECT create_user FROM grp where come_from = 1 limit 1")
			Expect(name).To(Equal("root"))
		})
	})
	Context("With select come_from = 0", func() {
		It("Name should be cmdb-test-grp-[ef]", func() {
			var names []string
			DbFacade.NewSqlxDbCtrl().Select(&names, "SELECT grp_name FROM grp where come_from = 0 order by grp_name")
			Expect(names).To(Equal([]string{"cmdb-test-grp-e", "cmdb-test-grp-f"}))
		})
	})
	Context("With select come_from = 0", func() {
		It("Creator should be default", func() {
			var name string
			DbFacade.NewSqlxDbCtrl().Get(&name, "SELECT create_user FROM grp where come_from = 0 limit 1")
			Expect(name).To(Equal("default"))
		})
	})
}))

var _ = Describe("[CMDB] syncRelTx", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
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
		// relation data for sync
		testCase := map[string][]string{
			"cmdb-test-grp-b": []string{"cmdb-test-a", "cmdb-test-b"},
			"cmdb-test-grp-c": []string{"cmdb-test-a", "cmdb-test-b"},
		}
		txProcessor := &syncRelTx{
			relations: testCase,
		}
		DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	})
	AfterEach(func() {
		inTx(
			`DELETE FROM grp_host WHERE grp_id IN (SELECT id FROM grp)`,
			`DELETE FROM grp_host WHERE host_id IN (SELECT id FROM host)`,
			`DELETE FROM host WHERE hostname LIKE "cmdb-test-%"`,
			`DELETE FROM grp WHERE grp_name LIKE "cmdb-test-grp-%"`,
		)
	})
	Context("With select only count", func() {
		It("count should be 6.", func() {
			var count int
			DbFacade.NewSqlxDbCtrl().Get(&count, "SELECT count(*) FROM grp_host")
			Expect(count).To(Equal(6))
		})
	})
	Context("With select group id = 10", func() {
		It("Hid should be 1 and 2", func() {
			var hids []int
			DbFacade.NewSqlxDbCtrl().Select(&hids, "SELECT host_id FROM grp_host where grp_id = 10 order by host_id")
			Expect(hids).To(Equal([]int{1, 2}))
		})
	})
	Context("With select group id = 20", func() {
		It("Hid should be 1 and 2", func() {
			var hids []int
			DbFacade.NewSqlxDbCtrl().Select(&hids, "SELECT host_id FROM grp_host where grp_id = 20 order by host_id")
			Expect(hids).To(Equal([]int{1, 2}))
		})
	})
	Context("With select group id = 30", func() {
		It("Hid should be 1 and 2", func() {
			var hids []int
			DbFacade.NewSqlxDbCtrl().Select(&hids, "SELECT host_id FROM grp_host where grp_id = 30 order by host_id")
			Expect(hids).To(Equal([]int{1, 2}))
		})
	})
}))

var _ = Describe("[CMDB] api2tuple()", itSkip.PrependBeforeEach(func() {
	testCase := []*cmdbModel.SyncHost{
		{
			Activate: 0,
			Name:     "cmdb-test-a",
			IP:       "69.69.69.1",
		},
		{
			Activate: 0,
			Name:     "cmdb-test-b",
			IP:       "69.69.69.2",
		},
		{
			Activate: 1,
			Name:     "cmdb-test-c",
			IP:       "69.69.69.3",
		},
		{
			Activate: 1,
			Name:     "cmdb-test-d",
			IP:       "69.69.69.4",
		},
	}
	Context("With activate 0", func() {
		It("maintain_begin should be MAINTAIN_PERIOD_BEGIN", func() {
			Expect(api2tuple(testCase)[0].Maintain_begin).To(Equal(uint32(MAINTAIN_PERIOD_BEGIN))) //  Sat, 01 Jan 2000 00:00:00 GMT
		})
		It("maintain_end should be MAINTAIN_PERIOD_END", func() {
			Expect(api2tuple(testCase)[0].Maintain_end).To(Equal(uint32(MAINTAIN_PERIOD_END))) //  Thu, 07 Jan 2106 17:43:40 GMT
		})
	})
	Context("With activate 1", func() {
		It("maintain_begin should be 0", func() {
			Expect(api2tuple(testCase)[3].Maintain_begin).To(Equal(uint32(0)))
		})
		It("maintain_end should be 0", func() {
			Expect(api2tuple(testCase)[3].Maintain_end).To(Equal(uint32(0)))
		})
	})
	Context("With name cmdb-test-a", func() {
		It("Hostname should be cmdb-test-a", func() {
			Expect(api2tuple(testCase)[0].Hostname).To(Equal("cmdb-test-a"))
		})
	})
	Context("With IP 69.69.69.1", func() {
		It("Ip should be 69.69.69.1", func() {
			Expect(api2tuple(testCase)[0].Ip).To(Equal("69.69.69.1"))
		})
	})
	Context("With testCase length 4", func() {
		It("output should be length 4", func() {
			Expect(len(api2tuple(testCase))).To(Equal(4))
		})
	})
}))
