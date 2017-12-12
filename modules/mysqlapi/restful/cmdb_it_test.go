package restful

import (
	"net/http"
	"time"

	oClient "github.com/Cepave/open-falcon-backend/common/http/client"
	oJson "github.com/Cepave/open-falcon-backend/common/json"
	oGko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	gb "github.com/Cepave/open-falcon-backend/common/testing/ginkgo/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("[POST] /api/v1/cmdb/sync", itSkipForCmdb.PrependBeforeEach(func() {
	setupImportAndAssertion(
		gb.NewGinkgoBuilder("Imports new data of host/host groups"),
		3, 2,
	).
		BeforeFirst(func() {
			inBossTx(
				`
				INSERT INTO hosts(hostname, activate, platform, ip, isp, exist)
				VALUES
					('it.h01.gp1', 1, 'it-g01', '10.6.51.1', '', 1),
					('it.h01.gp2', 1, 'it-g01', '10.6.51.2', '', 1),
					('it.h01.gp3', 1, 'it-g02', '10.6.51.2', '', 1)
				`,
			)
		}).
	ToContext()

	setupImportAndAssertion(
		gb.NewGinkgoBuilder("Imports over existing data"),
		5, 3,
	).
		BeforeFirst(func() {
			time.Sleep(1 * time.Second)

			inBossTx(
				`
				INSERT INTO hosts(hostname, activate, platform, ip, isp, exist)
				VALUES
					('it.h01.gp4', 1, 'it-g01', '10.6.51.4', '', 1),
					('it.h01.gp5', 1, 'it-g03', '10.6.51.5', '', 1)
				`,
			)
		}).
		AfterLast(func() {
			cleanImportData()
			cleanBossData()
		}).
	ToContext()
}))

var _ = Describe("[GET] /api/v1/cmdb/sync/:uuid", itSkipOnPortal.PrependBeforeEach(func() {
	Context("Found schedule log", func() {
		BeforeEach(func() {
			inPortalTx(
				`
				INSERT INTO owl_schedule(sch_id, sch_name, sch_lock, sch_modify_time)
				VALUES(312, 'gp.009', 0, NOW())
				`,
				`
				INSERT INTO owl_schedule_log(
					sl_sch_id, sl_uuid, sl_status,
					sl_timeout, sl_start_time, sl_end_time
				)
				VALUES(312, x'7dc535add25693024620dfe4258b9001', 0, 300, '2017-05-05 10:10:20', '2017-05-05 10:22:07'),
					(312, x'ec612430312ed634a5cbc2307780a03f', 1, 300, '2017-05-06 14:32:12', NULL)
				`,
			)
		})
		AfterEach(func() {
			inPortalTx(
				`
				DELETE FROM owl_schedule_log
				WHERE sl_sch_id = 312
				`,
				`
				DELETE FROM owl_schedule
				WHERE sch_id = 312
				`,
			)
		})

		DescribeTable("Content should match expected one",
			func(uuidString string, hasEndtime bool, expectedStatus int) {
				resp, err := gentlemanClientConfig.NewClient().
					Path("/api/v1/cmdb/sync/" + uuidString).Get().
					Send()

				Expect(err).To(Succeed())

				jsonBody := oClient.ToGentlemanResp(resp).MustGetJson()
				GinkgoT().Logf("[/api/v1/cmdb/%s] JSON Result:\n%s", uuidString, oJson.MarshalPrettyJSON(jsonBody))

				Expect(resp).To(oGko.MatchHttpStatus(http.StatusOK))

				Expect(jsonBody.Get("status").MustInt()).To(Equal(expectedStatus))
				Expect(jsonBody.Get("timeout").MustInt()).To(Equal(300))

				matchEndTime := BeNumerically(">=", 1)
				if !hasEndtime {
					matchEndTime = Equal(0)
				}
				Expect(jsonBody.Get("end_time").MustInt()).To(matchEndTime)
			},
			Entry("Has end time", "7dc535ad-d256-9302-4620-dfe4258b9001", true, 0),
			Entry("No end time", "ec612430-312e-d634-a5cb-c2307780a03f", false, 1),
		)
	})

	Context("Not found schedule log", func() {
		It("Should be 404", func() {
			resp, err := gentlemanClientConfig.NewClient().
				Path("/api/v1/cmdb/sync/80915aae-7b75-1eed-2db9-6e0b450172e4").Get().
				Send()

			Expect(err).To(Succeed())

			jsonBody := oClient.ToGentlemanResp(resp).MustGetJson()
			GinkgoT().Logf("[/api/v1/cmdb/sync] JSON Result:\n%s", oJson.MarshalPrettyJSON(jsonBody))

			Expect(resp).To(oGko.MatchHttpStatus(http.StatusNotFound))
		})
	})
}))

func setupImportAndAssertion(
	ginkgoBuilder *gb.GinkgoBuilder,
	expectedHosts int, expectedHostgroups int,
) *gb.GinkgoBuilder {
	ginkgoBuilder.
		It("The job should be started successfully", func() {
			resp, err := gentlemanClientConfig.NewClient().
				Path("/api/v1/cmdb/sync").Post().
				Send()

			Expect(err).To(Succeed())

			jsonBody := oClient.ToGentlemanResp(resp).MustGetJson()
			GinkgoT().Logf("[/api/v1/cmdb/sync] JSON Result:\n%s", oJson.MarshalPrettyJSON(jsonBody))

			Expect(resp).To(oGko.MatchHttpStatus(http.StatusOK))
			Expect(jsonBody.Get("sync_id").MustString()).To(MatchRegexp("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"))
			Expect(jsonBody.Get("start_time").MustInt64()).To(BeNumerically(">=", time.Now().Add(-1*time.Minute).Unix()))
		}).
		It("Waiting for finishing on scheduled job", func() {
			Eventually(
				func() int {
					var isReady int

					portalDbFacade.SqlxDbCtrl.Get(
						&isReady,
						`
						SELECT COUNT(*) FROM owl_schedule
						WHERE sch_name = 'import.imdb'
							AND sch_lock = 0
						`,
					)

					return isReady
				},
				5*time.Second, 250*time.Millisecond,
			).Should(Equal(1))
		}).
		It("The number of imported hosts and host groups", func() {
			var countHolder = &struct {
				CountHosts      int `db:"count_hosts"`
				CountHostGroups int `db:"count_hostgroups"`
			}{}

			portalDbFacade.SqlxDbCtrl.Get(
				countHolder,
				`
				SELECT
					(
						SELECT COUNT(*)
						FROM host
						WHERE hostname LIKE 'it.h01%'
					) AS count_hosts,
					(
						SELECT COUNT(*)
						FROM grp
						WHERE grp_name LIKE 'it-%'
					) AS count_hostgroups
				`,
			)

			Expect(countHolder).To(PointTo(
				MatchAllFields(Fields{
					"CountHosts":      Equal(expectedHosts),
					"CountHostGroups": Equal(expectedHostgroups),
				}),
			))
		})

	return ginkgoBuilder
}

func cleanBossData() {
	inBossTx(
		`
		DELETE FROM hosts
		WHERE hostname LIKe 'it.h01%'
		`,
	)
}

func cleanImportData() {
	inPortalTx(
		`
		DELETE owl_schedule_log
		FROM owl_schedule_log
			INNER JOIN
			owl_schedule
			ON sl_sch_id = sch_id
				AND sch_name = 'import.imdb'
		`,
		`
		DELETE FROM owl_schedule
		WHERE sch_name = 'import.imdb'
		`,
		`
		DELETE gh
		FROM grp_host AS gh
			INNER JOIN
			grp AS gp
			ON gh.grp_id = gp.id
				AND gp.grp_name LIKE 'it-%'
		`,
		`
		DELETE FROM grp
		WHERE grp_name LIKE 'it-%'
			OR
			grp_name = 'Owl_Default_Group'
		`,
		`
		DELETE FROM host
		WHERE hostname like 'it.h01.%'
		`,
	)
}
