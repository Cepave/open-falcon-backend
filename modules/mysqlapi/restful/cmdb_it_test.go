package restful

import (
	"net/http"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	oJson "github.com/Cepave/open-falcon-backend/common/json"
	oClient "github.com/Cepave/open-falcon-backend/common/http/client"
	oGko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	gb "github.com/Cepave/open-falcon-backend/common/testing/ginkgo/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("[POST] /api/v1/cmdb/sync", itSkip.PrependBeforeEach(func() {
	setupImportAndAssertion(
		gb.NewGinkgoBuilder("Imports new data of host/host groups"),
		&model.SyncForAdding {
			Hosts: []*model.SyncHost{
				{ Name: "it.h01.gp1", IP: "10.6.51.1", Activate: 1 },
				{ Name: "it.h01.gp2", IP: "10.6.51.2", Activate: 1 },
				{ Name: "it.h01.gp3", IP: "10.6.51.3", Activate: 0 },
			},
			Hostgroups: []*model.SyncHostGroup {
				{ Name: "it-g01", Creator: "it-user-91" },
				{ Name: "it-g02", Creator: "it-user-91" },
			},
			Relations: map[string][]string {
				"it-g01": { "it.h01.gp1", "it.h01.gp2" },
				"it-g02": { "it.h01.gp3" },
			},
		},
		3, 2,
	).ToContext()

	setupImportAndAssertion(
		gb.NewGinkgoBuilder("Imports over exsiting data"),
		&model.SyncForAdding {
			Hosts: []*model.SyncHost{
				{ Name: "it.h01.gp1", IP: "10.6.51.1", Activate: 1 },
				{ Name: "it.h01.gp4", IP: "10.6.51.2", Activate: 1 },
				{ Name: "it.h01.gp5", IP: "10.6.51.3", Activate: 0 },
			},
			Hostgroups: []*model.SyncHostGroup {
				{ Name: "it-g01", Creator: "it-user-91" },
				{ Name: "it-g03", Creator: "it-user-82" },
			},
			Relations: map[string][]string {
				"it-g01": { "it.h01.gp1", "it.h01.gp4" },
				"it-g03": { "it.h01.gp5" },
			},
		},
		5, 3,
	).
	BeforeFirst(func() {
		time.Sleep(1 * time.Second)
	}).
	AfterLast(func() {
		cleanImportData()
	}).ToContext()
}))

func setupImportAndAssertion(
	ginkgoBuilder *gb.GinkgoBuilder,
	importData interface{}, expectedHosts int, expectedHostgroups int,
) *gb.GinkgoBuilder {
	ginkgoBuilder.
		It("The job should be started successfully", func() {
			resp, err := gentlemanClientConfig.NewClient().
				Path("/api/v1/cmdb/sync").Post().
				JSON(importData).
				Send()

			Expect(err).To(Succeed())

			jsonBody := oClient.ToGentlemanResp(resp).MustGetJson()
			GinkgoT().Logf("[/api/v1/cmdb/sync] JSON Result:\n%s", oJson.MarshalPrettyJSON(jsonBody))

			Expect(resp).To(oGko.MatchHttpStatus(http.StatusOK))
			Expect(jsonBody.Get("sync_id").MustString()).To(MatchRegexp("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"))
			Expect(jsonBody.Get("start_time").MustInt64()).To(BeNumerically(">=", time.Now().Add(-1 * time.Minute).Unix()))
		}).
		It("Waiting for finishing on scheduled job", func() {
			Eventually(
				func() int {
					var isReady int

					dbFacade.SqlxDbCtrl.Get(
						&isReady,
						`
						SELECT COUNT(*) FROM owl_schedule
						WHERE sch_name = 'import.imdb'
							AND sch_lock = 0
						`,
					)

					return isReady
				},
				5 * time.Second, 250 * time.Millisecond,
			).Should(Equal(1))
		}).
		It("The number of imported hosts and host groups", func() {
			var countHolder = &struct {
				CountHosts int `db:"count_hosts"`
				CountHostGroups int `db:"count_hostgroups"`
			} {}

			dbFacade.SqlxDbCtrl.Get(
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
					"CountHosts": Equal(expectedHosts),
					"CountHostGroups": Equal(expectedHostgroups),
				}),
			))
		})

	return ginkgoBuilder
}

func cleanImportData() {
	inTx(
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
		`,
		`
		DELETE FROM host
		WHERE hostname like 'it.h01.%'
		`,
	)
}
