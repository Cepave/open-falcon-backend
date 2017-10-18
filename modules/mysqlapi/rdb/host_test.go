package rdb

import (
	cModel "github.com/Cepave/open-falcon-backend/common/model"
	rt "github.com/Cepave/open-falcon-backend/common/reflect/types"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[Intg] Tests ListHosts(...)", itSkip.PrependBeforeEach(func() {
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

	DescribeTable("ListHosts(<paging>)",
		func(pageSize int, order string, expectedHostIDs []int, expectedGroupIDs []string, expectedTotalCount int) {
			page := cModel.Paging{
				Size:    int32(pageSize),
				OrderBy: []*cModel.OrderByEntity{{order, cModel.Ascending}},
			}
			res, paging := ListHosts(page)

			resHostIDs := utils.MakeAbstractArray(res).
				MapTo(
					func(elem interface{}) interface{} {
						return elem.(*model.HostsResult).ID
					},
					rt.TypeOfInt,
				).GetArray()
			resGroupIDs := utils.MakeAbstractArray(res).
				MapTo(
					func(elem interface{}) interface{} {
						return elem.(*model.HostsResult).IdsOfGroups
					},
					rt.TypeOfString,
				).GetArray()

			Expect(resHostIDs).To(Equal(expectedHostIDs))
			Expect(resGroupIDs).To(Equal(expectedGroupIDs))
			Expect(paging.TotalCount).To(Equal(int32(expectedTotalCount)))
		},
		Entry("List all data", 5, "id", []int{1, 2, 3, 4}, []string{"1", "1,2", "2", ""}, 4),
		Entry("List all data", 5, "name", []int{1, 3, 2, 4}, []string{"1", "2", "1,2", ""}, 4),
		Entry("List 1st paging of data", 2, "id", []int{1, 2}, []string{"1", "1,2"}, 4),
	)
}))

var _ = Describe("[Unit] Test buildSortingClauseOfHosts(...)", func() {
	DescribeTable("buildSortingClauseOfHosts(<paging>)",
		func(paging *cModel.Paging, expectedClause string) {
			res := buildSortingClauseOfHosts(paging)
			Expect(res).To(Equal(expectedClause))
		},
		Entry("Undefined",
			&cModel.Paging{},
			"",
		),
		Entry("Sort by id ASC",
			&cModel.Paging{OrderBy: []*cModel.OrderByEntity{{"id", cModel.Ascending}}},
			"id ASC",
		),
		Entry("Sort by name ASC",
			&cModel.Paging{OrderBy: []*cModel.OrderByEntity{{"name", cModel.Ascending}}},
			"hostname ASC",
		),
	)
})

var _ = Describe("[Intg] Tests ListHostgroups", itSkip.PrependBeforeEach(func() {
	const (
		NumOfTestHostgroup = 3
	)

	BeforeEach(func() {
		inTx(
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
		inTx(
			`DELETE FROM grp WHERE grp_name LIKE 'listhostgroups-grpname-%'`,
			`DELETE FROM plugin_dir WHERE dir LIKE 'listhostgroups/plugin-%'`,
		)
	})

	DescribeTable("when set paging to",
		func(pageSize int, order string, expectedHostIDs []int, expectedGroupIDs []string, expectedTotalCount int) {
			page := cModel.Paging{
				Size:    int32(pageSize),
				OrderBy: []*cModel.OrderByEntity{{order, cModel.Ascending}},
			}
			res, paging := ListHostgroups(page)

			resHostIDs := utils.MakeAbstractArray(res).
				MapTo(
					func(elem interface{}) interface{} {
						return elem.(*model.HostgroupsResult).ID
					},
					rt.TypeOfInt,
				).GetArray()
			resGroupIDs := utils.MakeAbstractArray(res).
				MapTo(
					func(elem interface{}) interface{} {
						return elem.(*model.HostgroupsResult).IdsOfGroups
					},
					rt.TypeOfString,
				).GetArray()

			Expect(resHostIDs).To(Equal(expectedHostIDs))
			Expect(resGroupIDs).To(Equal(expectedGroupIDs))
			Expect(paging.TotalCount).To(Equal(int32(expectedTotalCount)))
		},
		Entry("List all data, sort by id", 5, "id", []int{1, 2, 3}, []string{"1,3,4", "2", ""}, NumOfTestHostgroup),
		Entry("List all data, sort by name", 5, "name", []int{3, 2, 1}, []string{"", "2", "1,3,4"}, NumOfTestHostgroup),
		Entry("List 1st paging of data", 2, "id", []int{1, 2}, []string{"1,3,4", "2"}, NumOfTestHostgroup),
	)
}))

var _ = Describe("[Unit] Test buildSortingClauseOfHostgroups", func() {
	DescribeTable("when the order sorted by",
		func(paging *cModel.Paging, expectedClause string) {
			res := buildSortingClauseOfHostgroups(paging)
			Expect(res).To(Equal(expectedClause))
		},
		Entry("Undefined",
			&cModel.Paging{},
			"",
		),
		Entry("id ASC",
			&cModel.Paging{OrderBy: []*cModel.OrderByEntity{{"id", cModel.Ascending}}},
			"id ASC",
		),
		Entry("name ASC",
			&cModel.Paging{OrderBy: []*cModel.OrderByEntity{{"name", cModel.Ascending}}},
			"grp_name ASC",
		),
	)
})
