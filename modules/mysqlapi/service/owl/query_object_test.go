package owl

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Tests the loading of query by UUID", ginkgoDb.NeedDb(func() {
	var testedSrv *queryObjectService

	BeforeEach(func() {
		testedSrv = newQueryObjectService(
			&QueryObjectServiceConfig{
				CacheSize:     8,
				CacheDuration: 5 * time.Minute,
			},
		)

		inTx(
			`
			INSERT INTO owl_query(
				qr_uuid, qr_named_id, qr_content, qr_md5_content, qr_time_creation, qr_time_access
			)
			VALUES(
				0x425dd1a3be3848f288d4bbc5db109fd9, 'test.1.feature',
				0xd94e3ed687cc9bb111b70c9341caaa4f, 0x547fc856e55394d8d4d0d1cbafe2d45d,
				NOW(), NOW()
			)
			`,
		)
	})

	AfterEach(func() {
		testedSrv.cache.Stop()
		testedSrv = nil
		inTx(
			`
			DELETE FROM owl_query
			WHERE qr_named_id = 'test.1.feature'
			`,
		)
	})

	DescribeTable("Loading of query object by UUID",
		func(sampleUuid string, expectedNameId string) {
			uuidValue, _ := uuid.FromString(sampleUuid)

			testedQuery := testedSrv.LoadQueryByUuid(uuidValue)

			if expectedNameId == "" {
				Expect(testedQuery).To(BeNil())
				Expect(testedSrv.cache.Get(sampleUuid)).To(BeNil())
			} else {
				Expect(testedQuery).To(Not(BeNil()))
				Expect(testedQuery.NamedId).To(Equal(expectedNameId))
				Expect(testedSrv.cache.Get(sampleUuid)).To(Not(BeNil()))
			}
		},
		Entry("Viable query object", "425dd1a3-be38-48f2-88d4-bbc5db109fd9", "test.1.feature"),
		Entry("UUID is not existing in database", "5a87da8a-4e37-4461-916d-0dd4d486684f", ""),
	)
}))
