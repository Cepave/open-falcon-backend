package graph

import (
	"net/http"

	oClient "github.com/Cepave/open-falcon-backend/common/http/client"
	json "github.com/Cepave/open-falcon-backend/common/json"
	tg "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	tc "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[POST] /api/v1/graph/endpoint-index/vacuum", itSkip.PrependBeforeEach(func() {
	client := (&tc.GentlemanClientConf{tc.NewHttpClientConfigByFlag()}).NewClient()

	BeforeEach(func() {
		inTx(
			`
			INSERT INTO endpoint(id, endpoint, ts, t_create)
			VALUES
				(30501, 'cmb-js-183-213-022-038', 1001, NOW()),
				(30502, 'cmb-sd-223-099-243-148', 1001, NOW()),
				(30503, 'cnc-gz-058-016-043-037', UNIX_TIMESTAMP(), NOW())
			`,
			`
			INSERT INTO endpoint_counter(
				id, endpoint_id, counter, step, type, ts, t_create
			)
			VALUES
				(20711, 30501, 'disk.io.msec_write/device=sdm', 60, 'GAUGE', 1001, NOW()),
				(20712, 30501, 'disk.io.await/device=sdd', 60, 'DERIVE', 1001, NOW()),
				(20713, 30502, 'net.if.in.multicast/iface=eth4', 60, 'GAUGE', 1001, NOW()),
				(20714, 30503, 'disk.io.read_merged/device=sde', 60, 'GAUGE', UNIX_TIMESTAMP(), NOW()),
				(20715, 30503, 'disk.io.read_sectors/device=sdf', 60, 'DERIVE', UNIX_TIMESTAMP(), NOW())
			`,
			`
			INSERT INTO tag_endpoint(
				id, endpoint_id, tag, ts, t_create
			)
			VALUES
				(23031, 30501, 'device=sds', 1001, NOW()),
				(23032, 30501, 'device=sdl', 1001, NOW()),
				(23033, 30502, 'iface=eth1', 1001, NOW()),
				(23034, 30503, 'iface=eth1', UNIX_TIMESTAMP(), NOW()),
				(23035, 30503, 'iface=eth3', UNIX_TIMESTAMP(), NOW())
			`,
		)
	})
	AfterEach(func() {
		inTx(
			"DELETE FROM tag_endpoint WHERE id >= 23031 AND id <= 23035",
			"DELETE FROM endpoint_counter WHERE id >= 20711 AND id <= 20715",
			"DELETE FROM endpoint WHERE id >= 30501 AND id <= 30503",
		)
	})

	sendVacuumRequest := func(
		expectedAffectedEndpoints, expectedAffectedCounters,
		expectedAffectedTags int,
	) {
		resp, err := client.Post().AddPath("/api/v1/graph/endpoint-index/vacuum").
			AddQuery("for_days", "1000").
			Send()
		Expect(err).To(Succeed())

		Expect(resp).To(tg.MatchHttpStatus(http.StatusOK))

		jsonBody := oClient.ToGentlemanResp(resp).MustGetJson()
		GinkgoT().Logf("[/graph/endpoint-index/vacuum] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))

		Expect(jsonBody.GetPath("affected_rows", "endpoints").MustInt()).To(Equal(expectedAffectedEndpoints))
		Expect(jsonBody.GetPath("affected_rows", "counters").MustInt()).To(Equal(expectedAffectedCounters))
		Expect(jsonBody.GetPath("affected_rows", "tags").MustInt()).To(Equal(expectedAffectedTags))
	}

	It("Send vacuum request", func() {
		By("1st vacuum(something should be vacuumed)")
		sendVacuumRequest(2, 3, 3)

		By("2nd vacuum(nothing to be vacuumed)")
		sendVacuumRequest(0, 0, 0)
	})
}))
