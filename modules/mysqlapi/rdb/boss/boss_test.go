package boss

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
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
			Expect(len(result)).To(Equal(2))
		})
	})
	Context("With hostname as 'boss-test-a', ip as '69.69.69.1', activate as 1, platform as c01.i01", func() {
		It("Hostname should be boss-test-a, Ip should be 69.69.69.1, Activate should be 1, Platform should be c01.i01", func() {
			result := GetSyncData()
			Expect(result[0]).To(PointTo(MatchAllFields(Fields{
				"Hostname": Equal("boss-test-a"),
				"Ip":       Equal("69.69.69.1"),
				"Activate": Equal(sql.NullInt64{Int64: 1, Valid: true}),
				"Platform": Equal(sql.NullString{String: "c01.i01", Valid: true}),
			})))
		})
	})
	Context("With activate as NULL, platform as NULL", func() {
		It("Activate should be Null, Platform should be Null", func() {
			result := GetSyncData()
			Expect(result[1]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Activate": Equal(sql.NullInt64{Int64: 0, Valid: false}),
				"Platform": Equal(sql.NullString{String: "", Valid: false}),
			})))
		})
	})
}))
