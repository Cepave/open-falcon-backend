package boss

import (
	"database/sql"
	"testing"

	model "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var _ = Describe("Tests boss2cmdb", func() {
	sampleData := []*model.BossHost{
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
	result := Boss2cmdb(sampleData)

	Context("Data of hosts", func() {
		It("Hosts should match expected result", func() {
			Expect(result.Hosts).To(ConsistOf(
				[]*model.SyncHost{
					{Name: "boss-test-a", IP: "69.69.69.1", Activate: 1},
					{Name: "boss-test-b", IP: "69.69.69.2", Activate: 1},
					{Name: "boss-test-c", IP: "69.69.69.3", Activate: 1},
					{Name: "boss-test-d", IP: "69.69.69.4", Activate: 1},
					{Name: "boss-test-e", IP: "69.69.69.5", Activate: 0},
					{Name: "boss-test-f", IP: "69.69.69.6", Activate: 0},
				},
			))
		})
	})

	Context("Data of host groups", func() {
		It("Host groups should match expected result", func() {
			Expect(result.Hostgroups).To(ConsistOf(
				[]*model.SyncHostGroup{
					{Name: "boss-platform-1", Creator: "root"},
					{Name: "boss-platform-2", Creator: "root"},
					{Name: DEFAULT_GRP, Creator: "root"},
				},
			))
		})
	})

	Context("Data of relations", func() {
		It("Relations should match expected result", func() {
			Expect(result.Relations).To(And(
				HaveKeyWithValue("boss-platform-1",
					ConsistOf([]string{"boss-test-a", "boss-test-b", "boss-test-d"}),
				),
				HaveKeyWithValue("boss-platform-2",
					ConsistOf([]string{"boss-test-c"}),
				),
				HaveKeyWithValue(DEFAULT_GRP,
					ConsistOf([]string{"boss-test-a", "boss-test-b", "boss-test-c", "boss-test-d", "boss-test-e", "boss-test-f"}),
				),
			))
		})
	})
})
