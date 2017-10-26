package graph

import (
	"time"

	gb "github.com/Cepave/open-falcon-backend/common/testing/ginkgo/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("VacuumEndpointIndex()", itSkip.PrependBeforeEach(func() {
	sampleBeforeTime := time.Unix(1003, 0)

	gb.NewGinkgoBuilder("Vacuum endpoint data").
		BeforeFirst(func() {
			inTx(
				`
				INSERT INTO endpoint(id, endpoint, ts, t_create)
				VALUES
					(10031, 'cmb-js-183-213-022-038', 1001, NOW()),
					(10032, 'cmb-sd-223-099-243-148', 1001, NOW()),
					(10033, 'cnc-gz-058-016-043-037', 1005, NOW())
				`,
				`
				INSERT INTO endpoint_counter(
					id, endpoint_id, counter, step, type, ts, t_create
				)
				VALUES
					(14031, 10031, 'disk.io.msec_write/device=sdm', 60, 'GAUGE', 1001, NOW()),
					(14032, 10031, 'disk.io.await/device=sdd', 60, 'DERIVE', 1001, NOW()),
					(14033, 10032, 'net.if.in.multicast/iface=eth4', 60, 'GAUGE', 1001, NOW()),
					(14034, 10033, 'disk.io.read_merged/device=sde', 60, 'GAUGE', 1005, NOW()),
					(14035, 10033, 'disk.io.read_sectors/device=sdf', 60, 'DERIVE', 1005, NOW())
				`,
				`
				INSERT INTO tag_endpoint(
					id, endpoint_id, tag, ts, t_create
				)
				VALUES
					(15031, 10031, 'device=sds', 1001, NOW()),
					(15032, 10031, 'device=sdl', 1001, NOW()),
					(15033, 10032, 'iface=eth1', 1001, NOW()),
					(15034, 10033, 'iface=eth1', 1005, NOW()),
					(15035, 10033, 'iface=eth3', 1005, NOW())
				`,
			)
		}).
		AfterLast(func() {
			inTx(
				"DELETE FROM tag_endpoint WHERE id >= 15031 AND id <= 15035",
				"DELETE FROM endpoint_counter WHERE id >= 14031 AND id <= 14035",
				"DELETE FROM endpoint WHERE id >= 10031 AND id <= 10033",
			)
		}).
		It("Vacuum data", func() {
			By("1st Vacuum, some rows are removed")
			result := VacuumEndpointIndex(sampleBeforeTime)
			Expect(result).To(PointTo(MatchAllFields(Fields{
				"CountOfVacuumedEndpoints": BeEquivalentTo(2),
				"CountOfVacuumedCounters":  BeEquivalentTo(3),
				"CountOfVacuumedTags":      BeEquivalentTo(3),
			})))

			By("2nd Vacuum, no row is removed")
			result = VacuumEndpointIndex(sampleBeforeTime)
			Expect(result).To(PointTo(MatchAllFields(Fields{
				"CountOfVacuumedEndpoints": BeEquivalentTo(0),
				"CountOfVacuumedCounters":  BeEquivalentTo(0),
				"CountOfVacuumedTags":      BeEquivalentTo(0),
			})))
		}).
		ToContext()
}))
