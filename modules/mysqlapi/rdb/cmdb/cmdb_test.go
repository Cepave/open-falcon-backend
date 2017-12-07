package cmdb

import (
	cmdbModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("[CMDB] syncHostTx", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			`
			INSERT INTO host (hostname, ip, maintain_begin, maintain_end)
			VALUES ("cmdb-test-a","69.69.69.99",946684800,4292329420),
				   ("cmdb-test-b","69.69.69.2", 0, 0),
				   ("cmdb-test-nm-e","69.69.69.5", 0, 0),
				   ("cmdb-test-nm-f","69.69.69.6", 0, 0)
			`)
		testCase := []*cmdbModel.SyncHost{
			{
				Activate: 0,
				Name:     "cmdb-test-a",
				IP:       "69.69.69.1",
			},
			{
				Activate: 1,
				Name:     "cmdb-test-b",
				IP:       "69.69.69.2",
			},
			{
				Activate: 1,
				Name:     "cmdb-test-new-c",
				IP:       "69.69.69.3",
			},
			{
				Activate: 1,
				Name:     "cmdb-test-new-d",
				IP:       "69.69.69.4",
			},
		}
		txProcessor := &syncHostTx{
			hosts: api2tuple(testCase),
		}
		DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	})
	AfterEach(func() {
		inTx(
			`DELETE FROM host WHERE hostname LIKE "cmdb-test-%"`,
		)
	})
	Context("Total number of hosts after importing(4 records before)", func() {
		It("The number should be 6", func() {
			var count int
			DbFacade.NewSqlxDbCtrl().Get(&count, "SELECT COUNT(*) FROM host")
			Expect(count).To(Equal(6))
		})
	})
	Context("Updated and inserted hosts to in-active", func() {
		It("maintain_begin and maintain_end should be pre-define values", func() {
			var mt int
			DbFacade.NewSqlxDbCtrl().Get(
				&mt,
				`
				SELECT COUNT(*) FROM host
				WHERE hostname = 'cmdb-test-a'
					AND ip = '69.69.69.1'
					AND maintain_begin = ? AND maintain_end = ?
				`,
				MAINTAIN_PERIOD_BEGIN, MAINTAIN_PERIOD_END,
			)
			Expect(mt).To(Equal(1))
		})
	})
	Context("Updated and inserted hosts to active", func() {
		It("maintain_begin and maintain_end should be pre-define values", func() {
			var mt int
			DbFacade.NewSqlxDbCtrl().Get(
				&mt,
				`
				SELECT COUNT(*) FROM host
				WHERE hostname = 'cmdb-test-b'
					AND ip = '69.69.69.2'
					AND maintain_begin = ? AND maintain_end = ?
				`,
				0, 0,
			)
			Expect(mt).To(Equal(1))
		})
	})
	Context("Updated and inserted with some hosts left not modified", func() {
		It("intact host number should be 2", func() {
			var mt int
			DbFacade.NewSqlxDbCtrl().Get(
				&mt,
				`SELECT COUNT(*) FROM host
				 WHERE hostname LIKE "cmdb-test-nm-%"
					AND maintain_begin = ? AND maintain_end = ?
				`,
				0, 0)
			Expect(mt).To(Equal(2))
		})
	})
	Context("Updated and inserted with some hosts inserted", func() {
		It("newly inserted host number be 2", func() {
			var mt int
			DbFacade.NewSqlxDbCtrl().Get(
				&mt,
				`SELECT COUNT(*) FROM host
				 WHERE hostname LIKE "cmdb-test-new-%"
					AND maintain_begin = ? AND maintain_end = ?
				`,
				0, 0)
			Expect(mt).To(Equal(2))
		})
	})
}))

var _ = Describe("[CMDB] syncHostGroupTx", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			`INSERT INTO grp (grp_name, create_user, come_from)
			 VALUES ("cmdb-test-grp-a", "default", 0),
					("cmdb-test-grp-nm-e", "default", 0),
					("cmdb-test-grp-nm-f", "default", 0)
			`)
		testCase := []*cmdbModel.SyncHostGroup{
			{
				Name:    "cmdb-test-grp-a",
				Creator: "root",
			},
			{
				Name:    "cmdb-test-grp-new-b",
				Creator: "root",
			},
			{
				Name:    "cmdb-test-grp-new-c",
				Creator: "root",
			},
			{
				Name:    "cmdb-test-grp-new-d",
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
	Context("Total number of groups after importing (4 records before)", func() {
		It("The number should be 6", func() {
			var count int
			DbFacade.NewSqlxDbCtrl().Get(&count, "SELECT COUNT(*) FROM grp")
			Expect(count).To(Equal(6))
		})
	})
	Context("Updated and inserted grp's creator and come_from", func() {
		It("Creator should be root and come_from should be 1", func() {
			var n int
			DbFacade.NewSqlxDbCtrl().Get(
				&n,
				`
				SELECT COUNT(*) FROM grp
				WHERE grp_name = 'cmdb-test-grp-a'
					AND come_from = 1
					AND create_user = 'root'
				`,
			)
			Expect(n).To(Equal(1))
		})
	})
	Context("Updated and inserted with some grps left not modified", func() {
		It("intact grp number be 2", func() {
			var n int
			DbFacade.NewSqlxDbCtrl().Get(
				&n,
				`
				SELECT COUNT(*) FROM grp
				WHERE grp_name LIKE "cmdb-test-grp-nm-%"
					AND come_from = 0
					AND create_user = 'default'
				`,
			)
			Expect(n).To(Equal(2))
		})
	})
	Context("Updated and inserted with some grps inserted", func() {
		It("newly inserted grp number be 3", func() {
			var n int
			DbFacade.NewSqlxDbCtrl().Get(
				&n,
				`
				SELECT COUNT(*) FROM grp
				WHERE grp_name LIKE "cmdb-test-grp-new-%"
					AND come_from = 1
					AND create_user = 'root'
				`,
			)
			Expect(n).To(Equal(3))
		})
	})
}))

var _ = Describe("[CMDB] syncHostgroupContaining", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			`
			INSERT INTO grp (id, grp_name, come_from)
			VALUES (10021, "cmdb-test-grp-a", 0),
				(10022, "cmdb-test-grp-b", 1),
				(10023, "cmdb-test-grp-c", 1),
				(10024, "cmdb-test-grp-d", 1)
			`,
			`
			INSERT INTO host (id, hostname)
			VALUES (91081, "cmdb-test-a"),
				(91082, "cmdb-test-b"),
				(91083, "cmdb-test-c"),
				(91084, "cmdb-test-d")
			`,
			`
			INSERT INTO grp_host(grp_id, host_id)
			VALUES (10021, 91081), (10021, 91082), (10022, 91083), (10023, 91084),
				(10024, 91081), (10024, 91082)
			`,
		)
		// relation data for sync
		testCase := map[string][]string{
			"cmdb-test-grp-b": {"cmdb-test-a", "cmdb-test-c"},
			"cmdb-test-grp-c": {"cmdb-test-a", "cmdb-test-b", "cmdb-test-c"},
		}
		txProcessor := &syncHostgroupContaining{
			relations: testCase,
		}
		DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	})
	AfterEach(func() {
		inTx(
			`DELETE FROM grp_host`,
			`DELETE FROM host WHERE hostname LIKE "cmdb-test-%"`,
			`DELETE FROM grp WHERE grp_name LIKE "cmdb-test-grp-%"`,
		)
	})
	Context("Total number of relations after importing (4 records before)", func() {
		It("The number should be 7.", func() {
			var count int
			DbFacade.NewSqlxDbCtrl().Get(&count, "SELECT COUNT(*) FROM grp_host")
			Expect(count).To(Equal(7))
		})
	})
	Context("For relations of hostgroup which come from UI ('grp_host.come_from' is 0)", func() {
		It("The ralations should be intact", func() {
			var hids []int
			DbFacade.NewSqlxDbCtrl().Select(
				&hids,
				`
				SELECT host_id FROM grp_host
				WHERE grp_id = 10021
				ORDER BY host_id
				`,
			)
			Expect(hids).To(Equal([]int{ 91081, 91082 }))
		})
	})
	Context("For relations of hostgroup which come_from BOSS ('grp_host.come_from' is 1)", func() {
		It("The relations with added hosts", func() {
			var hids []int
			DbFacade.NewSqlxDbCtrl().Select(
				&hids,
				`
				SELECT host_id FROM grp_host
				WHERE grp_id = 10022
				ORDER BY host_id
				`,
			)
			Expect(hids).To(Equal([]int{ 91081, 91083 }))
		})

		It("The relations with different set(totally)", func() {
			var hids []int
			DbFacade.NewSqlxDbCtrl().Select(
				&hids,
				`
				SELECT host_id FROM grp_host
				WHERE grp_id = 10023
				ORDER BY host_id
				`,
			)
			Expect(hids).To(Equal([]int{ 91081, 91082, 91083 }))
		})
	})
	Context("For relations no hosts", func() {
		It("The relations with added hosts", func() {
			var count int
			DbFacade.NewSqlxDbCtrl().Get(
				&count,
				`
				SELECT COUNT(host_id) FROM grp_host
				WHERE grp_id = 10024
				`,
			)
			Expect(count).To(Equal(0))
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
	result := api2tuple(testCase)
	Context("With testCase length 4", func() {
		It("result should be length 4", func() {
			Expect(len(result)).To(Equal(4))
		})
	})
	Context("With activate 0", func() {
		It("(maintain_begin, maintain_end) should be pre-defined values", func() {
			Expect(result[0]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"MaintainBegin": BeEquivalentTo(MAINTAIN_PERIOD_BEGIN), // Sat, 01 Jan 2000 00:00:00 GMT
				"MaintainEnd":   BeEquivalentTo(MAINTAIN_PERIOD_END),   // Thu, 07 Jan 2106 17:43:40 GMT
			})))
		})
	})
	Context("With activate 1", func() {
		It("(maintain_begin, maintain_end) should be pre-defined values", func() {
			Expect(result[3]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"MaintainBegin": BeEquivalentTo(0),
				"MaintainEnd":   BeEquivalentTo(0),
			})))
		})
	})
	Context("With name as 'cmdb-test-a', ip as '69.69.69.1'", func() {
		It("Hostname should be cmdb-test-a, Ip should be 69.69.69.1", func() {
			Expect(result[0]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Hostname": Equal("cmdb-test-a"),
				"Ip":       Equal("69.69.69.1"),
			})))
		})
	})
}))
