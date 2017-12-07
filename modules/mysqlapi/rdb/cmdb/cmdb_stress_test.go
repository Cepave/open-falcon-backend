package cmdb

import (
	"fmt"
	"math/rand"
	"strings"

	rdata "github.com/Pallinder/go-randomdata"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

const (
	hostnamePrefix = "st-cch"
	hostGroupPrefix = "st-ccg"
)

var _ = Describe("Stress testing for importing to \"host, grp, and grp_host\" tables", func() {
	var measureTimes = 0

	BeforeEach(func() {
		if measureTimes == 0 {
			Skip(fmt.Sprintf("The measure times is 0(\"var measureTimes = 0\"), skip the stress test."))
		}
	})

	var runTimes = 0
	AfterEach(func() {
		runTimes++
		if runTimes < measureTimes {
			return
		}

		GinkgoT().Logf("Remove all of generated hosts, host groups, and relations")
		inTx(
			`
			DELETE gh
			FROM grp_host AS gh
				INNER JOIN
				grp AS gp
				ON gh.grp_id = gp.id
					AND gp.grp_name LIKE 'st-ccg-%'
			`,
			`DELETE FROM host WHERE hostname LIKE "st-cch-%"`,
			`DELETE FROM grp WHERE grp_name LIKE "st-ccg-%"`,
		)
	})

	Context("Heavy data tests", func() {
		var (
			numberOfHosts = 20480
			numberOfHostGroups = 256
			hostGroupWeights = []int{ 1, 1, 10, 8, 30, 5, 6, 10, 20, 15, 15, 3 }
		)

		sampleData := generateSourceData(numberOfHosts, numberOfHostGroups, hostGroupWeights)

		GinkgoT().Logf("Number of hosts: [%d]. Number of host groups: [%d]. Group Weights: %v",
			numberOfHosts, numberOfHostGroups, hostGroupWeights)

		Measure("Heavy data tests", func(b Benchmarker) {
			b.Time("runtime", func() {
				SyncForHosts(sampleData)
			})

			assertByNumberOfHosts(numberOfHosts, numberOfHostGroups)
		}, measureTimes)
	})
})

func assertByNumberOfHosts(expectedNumberOfHosts int, expectedNumberOfHostGroups int) {
	var countHolder = &struct {
		CountHosts int `db:"count_hosts"`
		CountHostGroups int `db:"count_host_groups"`
		CountRelations int `db:"count_relations"`
	} {}

	DbFacade.NewSqlxDbCtrl().Get(
		countHolder,
		`
		SELECT
			(
				SELECT COUNT(*)
				FROM host
			) AS count_hosts,
			(
				SELECT COUNT(*)
				FROM grp
			) AS count_host_groups,
			(
				SELECT COUNT(*)
				FROM grp_host
			) AS count_relations
		`,
	)

	Expect(countHolder).To(PointTo(
		MatchAllFields(Fields{
			"CountHosts": Equal(expectedNumberOfHosts),
			"CountHostGroups": Equal(expectedNumberOfHostGroups),
			"CountRelations": Equal(expectedNumberOfHosts),
		}),
	))
}

func generateSourceData(numberOfHosts int, numberOfHostGroups int, hostGroupWeights []int) *model.SyncForAdding {
	sourceData := &model.SyncForAdding {
		Hosts: generatedHosts(numberOfHosts),
		Hostgroups: generatedHostGroups(numberOfHostGroups),
		Relations: make(map[string][]string, 0),
	}

	var currentWeight = 0
	var currentHostGroup = 0
	for i := 0; i < numberOfHosts; {
		for w := 0; w < hostGroupWeights[currentWeight] && i < numberOfHosts; w++ {
			currentNameOfHostGroup := sourceData.Hostgroups[currentHostGroup].Name

			sourceData.Relations[currentNameOfHostGroup] = append(
				sourceData.Relations[currentNameOfHostGroup], sourceData.Hosts[i].Name,
			)

			i++
		}

		currentWeight++
		currentWeight %= len(hostGroupWeights)

		currentHostGroup++
		currentHostGroup %= numberOfHostGroups
	}

	return sourceData
}

func generatedHostGroups(totalNumber int) []*model.SyncHostGroup {
	var result = make([]*model.SyncHostGroup, totalNumber)

	for i := 0; i < totalNumber; i++ {
		result[i] = &model.SyncHostGroup {
			Creator: rdata.FirstName(rdata.Male),
			Name: fmt.Sprintf("%s-%03d", hostGroupPrefix, i + 1),
		}
	}

	return result
}
func generatedHosts(totalNumber int) []*model.SyncHost {
	var result = make([]*model.SyncHost, totalNumber)

	var currentIp uint32 = (rand.Uint32() % 254 + 1) << 24
	for i := 0; i < totalNumber; i++ {
		currentIp++

		if currentIp & 0xFF == 0 {
			currentIp += 1
		}

		ipBytes := []byte {
			byte(currentIp >> 24 & 0xFF), byte(currentIp >> 16 & 0xFF),
			byte(currentIp >> 8 & 0xFF), byte(currentIp & 0xFF),
		}

		result[i] = &model.SyncHost {
			Name: fmt.Sprintf("%s-%s-%03d-%03d-%03d-%03d",
				hostnamePrefix, strings.ToLower(rdata.Letters(3)),
				ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3],
			),
			IP: fmt.Sprintf(
				"%d.%d.%d.%d",
				ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3],
			),
			Activate: rand.Int() % 2,
		}
	}

	return result
}
