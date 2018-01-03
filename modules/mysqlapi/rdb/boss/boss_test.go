package boss

import (
	"database/sql"

	model "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[Boss] GetSyncData", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			`
			INSERT INTO hosts (hostname, ip, isp, activate, platform, exist)
			    VALUES ("boss-test-a", "69.69.69.1", "ctl", 1, "c01.i01", 1),
				       ("boss-test-b", "69.69.69.2", "ctl", NULL, NULL, 1),
					   ("boss-test-c", "69.69.69.3", "ctl", 1, "c01.i01", 0)
			`)
	})
	AfterEach(func() {
		inTx(
			`DELETE FROM hosts WHERE hostname LIKE "boss-test-%"`,
		)
	})

	Context("With testcase with 2 case existed", func() {
		It("result should be length 2", func() {
			result := GetSyncData()

			Expect(result).To(ConsistOf(
				[]*model.BossHost{
					{
						Hostname: "boss-test-a", Ip: "69.69.69.1",
						Activate: sql.NullInt64{Int64: 1, Valid: true},
						Platform: sql.NullString{String: "c01.i01", Valid: true},
					},
					{
						Hostname: "boss-test-b", Ip: "69.69.69.2",
						Activate: sql.NullInt64{Int64: 0, Valid: false},
						Platform: sql.NullString{String: "", Valid: false},
					},
				},
			))
		})
	})
}))
