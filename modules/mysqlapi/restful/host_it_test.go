package restful

import (
	"math"
	"net/http"
	"strconv"

	ogko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[Intg] Test Listhosts", ginkgoDb.NeedDb(func() {
	const NumOfTestHost = 4

	BeforeEach(func() {
		inTx(
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
		inTx(
			`DELETE FROM host WHERE hostname LIKE 'listhosts-hostname-%'`,
			`DELETE FROM grp WHERE grp_name LIKE 'listhosts-grpname-%'`,
			`DELETE FROM grp_host WHERE grp_id <= 10`,
		)
	})

	DescribeTable("when page size is",
		func(pageSize int, expectedHostsNums int) {
			result := tHttp.NewResponseResultBySling(
				httpClientConfig.NewSlingByBase().Get("api/v1/hosts").Set("page-size", strconv.Itoa(pageSize)),
			)
			Expect(result).To(ogko.MatchHttpStatus(http.StatusOK))

			jsonBody := result.GetBodyAsJson()
			Expect(jsonBody.MustArray()).To(HaveLen(expectedHostsNums))
		},
		Entry("Zero", 0, 0),
		Entry("Two", 2, 2),
		Entry("Max", math.MaxInt32, NumOfTestHost),
	)
}))
