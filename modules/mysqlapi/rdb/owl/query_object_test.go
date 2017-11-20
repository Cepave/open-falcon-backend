package owl

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/Cepave/open-falcon-backend/common/db"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	t "github.com/Cepave/open-falcon-backend/common/testing"
	"github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query object", itSkip.PrependBeforeEach(func() {
	Context("Add by content(md5)", func() {
		assertAddOrRefresh := func(accessTime string) *owlModel.Query {
			timeValue := t.ParseTimeByGinkgo(accessTime)

			sampleQuery := &owlModel.Query{
				NamedId:    "test.query.g1",
				Content:    []byte{29, 87, 61, 4, 5, 78, 91, 11, 91},
				Md5Content: md5.Sum([]byte("This is test - 1")),
			}

			AddOrRefreshQuery(sampleQuery, timeValue)

			Expect(sampleQuery.Uuid.IsNil()).To(BeFalse())
			Expect(sampleQuery.AccessTime.Unix()).To(Equal(timeValue.Unix()))

			return sampleQuery
		}

		Context("Adds query object", func() {
			It("The time of creation and access should be same", func() {
				testedQuery := assertAddOrRefresh("2015-07-01T07:36:55+08:00")
				Expect(testedQuery.CreationTime).To(Equal(testedQuery.AccessTime))
			})
		})

		Context("Refresh existing one", func() {
			AfterEach(func() {
				inTx(
					`
					DELETE FROM owl_query
					WHERE qr_named_id = 'test.query.g1'
					`,
				)
			})

			It("The time of creation and access should not be same", func() {
				testedQuery := assertAddOrRefresh("2015-07-01T07:46:55+08:00")
				Expect(testedQuery.CreationTime).To(Not(Equal(testedQuery.AccessTime)))
			})
		})
	})

	Context("Load by UUID", func() {
		Context("viable UUID(existing)", func() {
			BeforeEach(func() {
				inTx(
					`
					INSERT INTO owl_query(
						qr_uuid, qr_named_id,
						qr_content, qr_md5_content,
						qr_time_access, qr_time_creation
					)
					VALUES(
						x'209f18f4f89b42568e1e5270987c057d', 'test.load.uu2',
						x'7011e902d4a848c184e242e8d71aa961', x'0dfaa9f3df2d4071b2f48b359440e0fc',
						'2012-05-06T20:14:43', '2012-05-05T08:23:03'
					)
					`,
				)
			})

			AfterEach(func() {
				inTx(
					`
					DELETE FROM owl_query
					WHERE qr_named_id = 'test.load.uu2'
					`,
				)
			})

			It("Check loaded query access time and content", func() {
				sampleUuid, _ := uuid.FromString("209f18f4-f89b-4256-8e1e-5270987c057d")
				sampleTime := t.ParseTimeByGinkgo("2013-07-08T10:20:36+08:00")
				testedQuery := LoadQueryByUuidAndUpdateAccessTime(
					sampleUuid, sampleTime,
				)

				Expect(testedQuery).To(Not(BeNil()))
				Expect(testedQuery.NamedId).To(Equal("test.load.uu2"))
				Expect(testedQuery.AccessTime.Unix()).To(Equal(sampleTime.Unix()))

				expectedContent, _ := hex.DecodeString("7011e902d4a848c184e242e8d71aa961")
				Expect(testedQuery.Content).To(Equal(expectedContent))
				expectedMd5Content, _ := hex.DecodeString("0dfaa9f3df2d4071b2f48b359440e0fc")
				Expect(testedQuery.Md5Content[:]).To(Equal(expectedMd5Content))

				By("Use older access time")
				testedQuery = LoadQueryByUuidAndUpdateAccessTime(
					sampleUuid, t.ParseTimeByGinkgo("2013-07-01T10:20:36+08:00"),
				)
				Expect(testedQuery.AccessTime.Unix()).To(Equal(sampleTime.Unix()))
			})
		})

		Context("Not-existing UUID", func() {
			It("Check nil result of \"*Query\" object", func() {
				sampleUuid, _ := uuid.FromString("739f18f4-f89b-4a56-8ece-5070987c057d")
				sampleTime := t.ParseTimeByGinkgo("2013-07-08T10:20:36+08:00")
				testedQuery := LoadQueryByUuidAndUpdateAccessTime(
					sampleUuid, sampleTime,
				)

				Expect(testedQuery).To(BeNil())
			})
		})
	})

	Context("Update access time or creating new one", func() {
		Context("Update viable", func() {
			BeforeEach(func() {
				inTx(
					`
					INSERT INTO owl_query(
						qr_uuid, qr_named_id,
						qr_content, qr_md5_content,
						qr_time_access, qr_time_creation
					)
					VALUES(
						x'890858a7d458435bb7981ac0abf1eae2', 'test.query.uu1',
						x'890858a7d458435bb7981ac0abf1eae2', x'e2f2384748f14c9c8dbd8f276d2222eb',
						'2012-05-06T20:14:43', '2012-05-05T08:23:03'
					)
					`,
				)
			})

			AfterEach(func() {
				inTx(
					`
					DELETE FROM owl_query
					WHERE qr_named_id = 'test.query.uu1'
					`,
				)
			})

			It("Access time should be updated", func() {
				sampleUuid, _ := uuid.FromString("890858a7-d458-435b-b798-1ac0abf1eae2")
				sampleTime := t.ParseTimeByGinkgo("2012-06-22T10:20:31+08:00")

				sampleQuery := &owlModel.Query{
					Uuid:    db.DbUuid(sampleUuid),
					NamedId: "test.query.uu1",
				}
				UpdateAccessTimeOrAddNewOne(sampleQuery, sampleTime)

				Expect(time.Time(getAccessTimeByUuid(db.DbUuid(sampleUuid))).Unix()).
					To(Equal(sampleTime.Unix()))
			})
		})

		Context("Add not-existing one", func() {
			AfterEach(func() {
				inTx(
					`
					DELETE FROM owl_query
					WHERE qr_named_id = 'test.query.dc1'
					`,
				)
			})

			It("Check the creation time and access time(should be same)", func() {
				sampleUuid, _ := uuid.FromString("970856a7-d428-4e5b-b798-1ac0abb10ae2")
				sampleTime := t.ParseTimeByGinkgo("2013-02-12T13:30:07+08:00")

				sampleQuery := &owlModel.Query{
					Uuid:    db.DbUuid(sampleUuid),
					NamedId: "test.query.dc1",
					Content: []byte{0x23, 0x81, 0x91, 0x88, 0x14, 0x82},
					Md5Content: db.Bytes16{
						0xf9, 0xc0, 0x02, 0x08, 0xd3, 0xcf, 0xa9, 0x59,
						0x03, 0x20, 0xd6, 0x19, 0xbf, 0xe0, 0x5c, 0x72,
					},
				}
				UpdateAccessTimeOrAddNewOne(sampleQuery, sampleTime)

				Expect(sampleQuery.CreationTime.Unix()).To(Equal(sampleTime.Unix()))
				Expect(sampleQuery.AccessTime.Unix()).To(Equal(sampleTime.Unix()))
			})
		})
	})

	Context("Remove older query objects", func() {
		BeforeEach(func() {
			inTx(`
				INSERT INTO owl_query(
					qr_uuid, qr_named_id,
					qr_content, qr_md5_content,
					qr_time_access, qr_time_creation
				)
				VALUES
				(
					x'209f18f4f89b42568e1e5270987c057d', 'del.f1',
					x'7011e902d4a848c184e242e8d71aa961', x'bdfa89f3df204071b2f48b359440e0fc',
					'2014-05-06T20:14:43', '2014-05-05T08:23:03'
				),
				(
					x'349f18f4f89b42568e1e5270987c057d', 'del.f1',
					x'7011e902d4a848c184e242e8d71aa961', x'1afaa9f3df2d4071b2f48b359440e0fc',
					'2014-06-06T20:14:43', '2014-05-05T08:23:03'
				),
				(
					x'109f1cf4f89b42568e1e5270987c057d', 'del.f2',
					x'7011e902d4a848c184e242e8d71aa961', x'12faa9f3df2d4071b9f48b359440e0fc',
					'2014-05-06T20:14:43', '2014-05-05T08:23:03'
				),
				(
					x'a09f18f4f09b42568e1e5270987c057d', 'del.f2',
					x'7011e902d4a848c184e242e8d71aa961', x'6dfaa9f39f2d4071b2f43b359440e0fc',
					'2014-06-06T20:14:43', '2014-05-05T08:23:03'
				)
			`)
		})

		AfterEach(func() {
			inTx(`
				DELETE FROM owl_query
				WHERE qr_named_id IN ('del.f1', 'del.f2')
			`)
		})

		It("Check the number of removed rows", func() {
			sampleTime := t.ParseTimeByGinkgo("2014-06-01T00:00:00+08:00")
			By("Remove 2 rows")

			Expect(RemoveOldQueryObject(sampleTime)).To(BeEquivalentTo(2))

			By("Nothing to be removed")

			Expect(RemoveOldQueryObject(sampleTime)).To(BeEquivalentTo(0))
		})
	})
}))

func getAccessTimeByUuid(uuid db.DbUuid) time.Time {
	var timeValue = time.Time{}
	DbFacade.SqlxDbCtrl.QueryRowxAndScan(
		`
		SELECT qr_time_access
		FROM owl_query
		WHERE qr_uuid = ?
		`,
		[]interface{}{uuid},
		&timeValue,
	)

	return timeValue
}
