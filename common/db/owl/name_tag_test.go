package owl

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	rt "github.com/Cepave/open-falcon-backend/common/reflect/types"

	"github.com/Cepave/open-falcon-backend/common/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
)

func inTx(sqls ...string) func() {
	return func() {
		DbFacade.SqlDbCtrl.ExecQueriesInTx(sqls...)
	}
}

var _ = Describe("Tests GetNameTagById(...)", ginkgoDb.NeedDb(func() {
	/**
	 * Prepares data
	 */
	BeforeEach(inTx(
		`
		INSERT INTO owl_name_tag(nt_id, nt_value)
		VALUES(2901, 'here-we-1')
		`,
	))
	AfterEach(inTx(
		`
		DELETE FROM owl_name_tag WHERE nt_id = 2901
		`,
	))
	// :~)

	DescribeTable("GetNameTagById(<id>)",
		func(sampleId int, matcher GomegaMatcher) {
			Expect(GetNameTagById(int16(sampleId))).To(matcher)
		},
		Entry("Existing id", 2901, Not(BeNil())),
		Entry("Not existing id", 2902, BeNil()),
	)
}))

var _ = Describe("Tests ListNameTags(...)", ginkgoDb.NeedDb(func() {
	/**
	 * Prepares data
	 */
	BeforeEach(inTx(
		`
		INSERT INTO owl_name_tag(nt_id, nt_value)
		VALUES(3701, 'pg-tg-car-1'),
			(3702, 'pg-tg-car-2'),
			(3703, 'pg-tg-bird-3'),
			(3704, 'pg-tg-bird-4')
		`,
	))
	AfterEach(inTx(
		`
		DELETE FROM owl_name_tag WHERE nt_id >= 3701 AND nt_id <= 3704
		`,
	))
	// :~)

	DescribeTable("ListNameTags(<value>, <paging>)",
		func(value string, pageSize int, expectedIds []int16, expectedTotalCount int) {
			paging := &model.Paging{
				Size: int32(pageSize),
			}
			testedResult := ListNameTags(
				value, paging,
			)

			testedIds := utils.MakeAbstractArray(testedResult).
				MapTo(
					func(elem interface{}) interface{} {
						return elem.(*owlModel.NameTag).Id
					},
					rt.TypeOfInt16,
				).GetArray()

			Expect(testedIds).To(Equal(expectedIds))
			Expect(paging.TotalCount).To(Equal(int32(expectedTotalCount)))
		},
		Entry("All all of the data", "", 5, []int16{3703, 3704, 3701, 3702}, 4),
		Entry("List 1st paging of data", "", 2, []int16{3703, 3704}, 4),
		Entry("List certain value of name tag", "pg-tg-bird", 5, []int16{3703, 3704}, 2),
	)
}))
