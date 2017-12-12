package restful

import (
	"math"
	"net/http"
	"strconv"

	json "github.com/Cepave/open-falcon-backend/common/json"
	ogko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[Intg] Test listhosts", itSkipOnPortal.PrependBeforeEach(func() {
	const NumOfTestHost = 4

	BeforeEach(func() {
		inPortalTx(
			`INSERT INTO host(id, hostname)
				VALUES
					(1, 'listhosts-hostname-a1'),
					(2, 'listhosts-hostname-b2'),
					(3, 'listhosts-hostname-a3'),
					(4, 'listhosts-hostname-b4')`,
			`INSERT INTO grp(id, grp_name)
				VALUES
					(1, 'listhosts-grpname-1'),
					(2, 'listhosts-grpname-2')`,
			`INSERT INTO grp_host(grp_id, host_id)
				VALUES
					(1, 1),
					(1, 2),
					(2, 2),
					(2, 3)`,
		)
	})
	AfterEach(func() {
		inPortalTx(
			`DELETE FROM host WHERE hostname LIKE 'listhosts-hostname-%'`,
			`DELETE FROM grp WHERE grp_name LIKE 'listhosts-grpname-%'`,
			`DELETE FROM grp_host WHERE grp_id <= 10`,
		)
	})

	DescribeTable("when page size is",
		func(pageSize int, expectedHostsNums int) {
			result := tHttp.NewResponseResultBySling(
				httpClientConfig.NewClient().Get("api/v1/hosts").Set("page-size", strconv.Itoa(pageSize)),
			)
			Expect(result).To(ogko.MatchHttpStatus(http.StatusOK))

			jsonBody := result.GetBodyAsJson()
			GinkgoT().Logf("[/hosts] JSON Body: %s", json.MarshalPrettyJSON(jsonBody))
			Expect(jsonBody.MustArray()).To(HaveLen(expectedHostsNums))
		},
		Entry("Zero", 0, 0),
		Entry("Two", 2, 2),
		Entry("Max", math.MaxInt32, NumOfTestHost),
	)
}))

var _ = Describe("[Intg] Test listHostgroups", itSkipOnPortal.PrependBeforeEach(func() {
	const NumOfTestHostgroup = 3

	BeforeEach(func() {
		inPortalTx(
			`INSERT INTO grp(id, grp_name)
				VALUES
					(1, 'listhostgroups-grpname-c1'),
					(2, 'listhostgroups-grpname-b2'),
					(3, 'listhostgroups-grpname-a3')`,
			`INSERT INTO plugin_dir(id, grp_id, dir)
				VALUES
					(1, 1, 'listhostgroups/plugin-1'),
					(2, 2, 'listhostgroups/plugin-2'),
					(3, 1, 'listhostgroups/plugin-3'),
					(4, 1, 'listhostgroups/plugin-4')`,
		)
	})
	AfterEach(func() {
		inPortalTx(
			`DELETE FROM grp WHERE grp_name LIKE 'listhostgroups-grpname-%'`,
			`DELETE FROM plugin_dir WHERE dir LIKE 'listhostgroups/plugin-%'`,
		)
	})

	DescribeTable("when page size is",
		func(pageSize int, expectedHostsNums int) {
			result := tHttp.NewResponseResultBySling(
				httpClientConfig.NewClient().Get("api/v1/hostgroups").Set("page-size", strconv.Itoa(pageSize)),
			)
			Expect(result).To(ogko.MatchHttpStatus(http.StatusOK))

			jsonBody := result.GetBodyAsJson()
			GinkgoT().Logf("[/hostgroups] JSON Body: %s", json.MarshalPrettyJSON(jsonBody))
			Expect(jsonBody.MustArray()).To(HaveLen(expectedHostsNums))
		},
		Entry("Zero", 0, 0),
		Entry("Two", 2, 2),
		Entry("Max", math.MaxInt32, NumOfTestHostgroup),
	)
}))
