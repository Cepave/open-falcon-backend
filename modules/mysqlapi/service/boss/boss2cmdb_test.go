package boss

import (
	"database/sql"
	"fmt"
	model "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"testing"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

// how NullString initialized ?
// how NullInt64 initialzied?
var _ = Describe("Tests boss2cmdb", func() {
	testCase := []*model.BossHost{
		{
			Hostname: "boss-test-a",
			Ip:       "69.69.69.1",
			Activate: sql.NullInt64{Int64: 1, Valid: true},
			Platform: sql.NullString{String: "boss-platform-1", Valid: true},
		},
		{
			Hostname: "boss-test-b",
			Ip:       "69.69.69.2",
			Activate: sql.NullInt64{Int64: 1, Valid: true},
			Platform: sql.NullString{String: "boss-platform-1", Valid: true},
		},
		{
			Hostname: "boss-test-c",
			Ip:       "69.69.69.3",
			Activate: sql.NullInt64{Int64: 1, Valid: true},
			Platform: sql.NullString{String: "boss-platform-2", Valid: true},
		},
		{
			Hostname: "boss-test-d",
			Ip:       "69.69.69.4",
			Activate: sql.NullInt64{Int64: 1, Valid: true},
			Platform: sql.NullString{String: "boss-platform-1", Valid: true},
		},
		{
			Hostname: "boss-test-e",
			Ip:       "69.69.69.5",
			Activate: sql.NullInt64{Int64: 0, Valid: false},
			Platform: sql.NullString{String: "", Valid: false},
		},
		{
			Hostname: "boss-test-f",
			Ip:       "69.69.69.6",
			Activate: sql.NullInt64{Int64: 0, Valid: false},
			Platform: sql.NullString{String: "", Valid: false},
		},
	}
	fmt.Println("b")
	result := Boss2cmdb(testCase)
	fmt.Println("a")
	Context("With testCase has 6 valid Hosts", func() {
		It("result.Hosts should be length 6", func() {
			Expect(len(result.Hosts)).To(Equal(6))
		})
	})
	Context("With testCase has 2 valid HostGroups", func() {
		It("result.Hostgroups should be length 3", func() {
			Expect(len(result.Hostgroups)).To(Equal(3))
		})
	})
	Context("With testCase has 4 valid Relations", func() {
		It("result.Relations should be length 3, which means 3 groups", func() {
			Expect(len(result.Relations)).To(Equal(3))
		})
	})
	Context("With name as 'boss-test-a', ip as '69.69.69.1', Activate is 1", func() {
		It("Hostname should be cmdb-test-a, Ip should be 69.69.69.1, Activate is 1", func() {
			Expect(result.Hosts[0]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":     Equal("boss-test-a"),
				"IP":       Equal("69.69.69.1"),
				"Activate": Equal(1),
			})))
		})
	})
	Context("With group name as 'boss-platform-1'", func() {
		It("Name should be 'boss-platform-1', Creator should be 'root'.", func() {
			Expect(result.Hostgroups[0]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":    Equal("boss-platform-1"),
				"Creator": Equal("root"),
			})))
		})
	})
	Context("With relation of DEFAULT_GRP, boss-platform-1, boss-platform-2", func() {
		It("Should be 4 hosts in DEFAULT_GRP", func() {
			Expect(len(result.Relations[DEFAULT_GRP])).To(Equal(6))
		})
		It("Should be 3 hosts in boss-platform-1", func() {
			Expect(len(result.Relations["boss-platform-1"])).To(Equal(3))
		})
		It("Should be 1 hosts in boss-platform-2", func() {
			Expect(len(result.Relations["boss-platform-2"])).To(Equal(1))
		})
	})
})
