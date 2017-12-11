package boss

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("[Boss] GetSyncData", itSkip.PrependBeforeEach(func() {
	BeforeEach(func() {
		inTx(
			`
			INSERT INTO hosts (hostname, ip, isp, activate, platform, exist) VALUES ("boss-test-a", "69.69.69.1", "ctl", 1, "c01.i01", 1), ("boss-test-b", "69.69.69.2", "ctl", 0, "c01.i01", 1), ("boss-test-c", "69.69.69.3", "ctl", 1, "c01.i01", 0)
			`)
	})
	AfterEach(func() {
		inTx(
			`DELETE FROM hosts WHERE hostname LIKE "boss-test-%"`,
		)
	})

	Context("With testcase with 2 case existed", func() {
		result := GetSyncData()
		It("result should be length 2", func() {
			Expect(len(result)).To(Equal(2))
		})
	})
	Context("With hostname as 'boss-test-a', ip as '69.69.69.1', activate as 1, platform as ctl", func() {
		result := GetSyncData()
		It("Hostname should be cmdb-test-a, Ip should be 69.69.69.1, Activate should be 1, Platform should be ctl", func() {
			Expect(result[0]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Hostname": Equal("cmdb-boss-a"),
				"Ip":       Equal("69.69.69.1"),
				"Activate": Equal(1),
				"Platform": Equal("ctl"),
			})))
		})
	})
}))
