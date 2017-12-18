package cron

import (
	"github.com/satori/go.uuid"
	"time"

	oJson "github.com/Cepave/open-falcon-backend/common/json"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	owlSrv "github.com/Cepave/open-falcon-backend/common/service/owl"

	"github.com/Cepave/open-falcon-backend/modules/task/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
)

var _ = Describe("\"Do()\" of syncCmdbFromBoss", func() {
	var testedJob = &syncCmdbFromBoss{}

	AfterEach(func() {
		database.CmdbService = nil
	})

	DescribeTable("The error should match expected one",
		func(cmdbSrv owlSrv.CmdbService, errorMatcher GomegaMatcher) {
			database.CmdbService = cmdbSrv

			testedErr := testedJob.Do()
			Expect(testedErr).To(errorMatcher)
		},
		Entry(
			"Success(no error)",
			&sampleCmdbService{
				startTime: time.Now(),
				statusInfo: []owlModel.TaskStatus{
					owlModel.JobRunning, owlModel.JobDone,
				},
				timeoutValue: 8,
			},
			Succeed(),
		),
		Entry(
			"Timeout",
			&sampleCmdbService{
				startTime: time.Now(),
				statusInfo: []owlModel.TaskStatus{
					owlModel.JobRunning, owlModel.JobRunning,
				},
				timeoutValue: 1,
			},
			SatisfyAll(
				HaveOccurred(),
				WithTransform(
					func(err error) string {
						return err.Error()
					},
					MatchRegexp("timeout"),
				),
			),
		),
		Entry(
			"Failed",
			&sampleCmdbService{
				startTime: time.Now(),
				statusInfo: []owlModel.TaskStatus{
					owlModel.JobRunning, owlModel.JobFailed,
				},
				timeoutValue: 8,
			},
			SatisfyAll(
				HaveOccurred(),
				WithTransform(
					func(err error) string {
						return err.Error()
					},
					MatchRegexp("failed"),
				),
			),
		),
	)
})

type sampleCmdbService struct {
	statusIndex int

	startTime    time.Time
	statusInfo   []owlModel.TaskStatus
	timeoutValue int
}

func (s *sampleCmdbService) StartSyncJob() (uuid.UUID, error) {
	return uuid.NewV4(), nil
}
func (s *sampleCmdbService) GetJobStatus(syncId uuid.UUID) *owlModel.SyncCmdbJobInfo {
	if s.statusIndex >= len(s.statusInfo) {
		return nil
	}

	currentStatus := s.statusInfo[s.statusIndex]
	s.statusIndex++
	endTime := time.Time{}
	if currentStatus != owlModel.JobRunning {
		endTime = time.Time(s.startTime).Add(16)
	}

	return &owlModel.SyncCmdbJobInfo{
		StartTime: oJson.JsonTime(s.startTime),
		EndTime:   oJson.JsonTime(endTime),
		Status:    currentStatus,
		Timeout:   s.timeoutValue,
	}
}
