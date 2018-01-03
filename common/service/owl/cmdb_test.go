package owl

import (
	"net/http"
	"time"

	oJson "github.com/Cepave/open-falcon-backend/common/json"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"

	"github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("[RESTful Client] CMDB synchronization", func() {
	mysqlApiConfig := gockConfig.NewRestfulClientConfig()

	testedSrv := NewCmdbService(
		CmdbServiceConfig{mysqlApiConfig},
	)

	var sampleUuid uuid.UUID
	BeforeEach(func() {
		sampleUuid = uuid.NewV4()
	})

	Context("Start a task", func() {
		Context("Start a task successfully", func() {
			BeforeEach(func() {
				gockConfig.New().
					Post("/api/v1/cmdb/sync").
					Reply(http.StatusOK).
					JSON(map[string]interface{}{
						"sync_id":    sampleUuid.String(),
						"start_time": time.Now().Unix(),
					})
			})

			It("Uuid should be valid", func() {
				testedUuid, err := testedSrv.StartSyncJob()

				Expect(err).To(Succeed())
				Expect(testedUuid).To(Equal(sampleUuid))
			})
		})

		Context("Cannot start a task because of locking", func() {
			BeforeEach(func() {
				gockConfig.New().
					Post("/api/v1/cmdb/sync").
					Reply(http.StatusBadRequest).
					JSON(map[string]interface{}{
						"error_code":    1,
						"error_message": "Has error",
						"last_sync_id":  sampleUuid.String(),
					})
			})

			It("\"error_code == 1\"", func() {
				_, err := testedSrv.StartSyncJob()

				Expect(err).To(BeAssignableToTypeOf(&model.UnableToLockCmdbSyncJob{}))
				Expect(err.(*model.UnableToLockCmdbSyncJob).Uuid).To(BeEquivalentTo(sampleUuid))
			})
		})
	})

	Context("Get info of a task", func() {
		BeforeEach(func() {
			gockConfig.New().
				Get("/api/v1/cmdb/sync/" + sampleUuid.String()).
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"status":     1,
					"start_time": time.Now().Add(-time.Minute).Unix(),
					"end_time":   nil,
					"timeout":    500,
				})
		})

		It("Job's info should match expected one", func() {
			testedInfo := testedSrv.GetJobStatus(sampleUuid)

			Expect(testedInfo).To(PointTo(
				MatchAllFields(Fields{
					"Status": BeEquivalentTo(model.JobRunning),
					"StartTime": WithTransform(
						func(v oJson.JsonTime) time.Time {
							return time.Time(v)
						},
						BeTemporally("<=", time.Now()),
					),
					"EndTime": Equal(oJson.JsonTime(time.Time{})),
					"Timeout": Equal(500),
				}),
			))
		})
	})
})
