package cron

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe(`Usage of "AutoJob()" to create Job instance`, func() {
	Context(`Valid "Job" object is constructed`, func() {
		It(`Simple implementation of "Do() error"`, func() {
			sampleJob := new(simpleJob)

			testedJob := AutoJob(sampleJob)
			testedJob.Do()

			Expect(*sampleJob).To(BeEquivalentTo(1))
		})

		It(`Full implementation of "Job"`, func() {
			sampleJob := &fullJob{false, false, false, false}

			testedJob := AutoJob(sampleJob)
			testedJob.BeforeJobStart()
			testedJob.Do()
			testedJob.BeforeJobStop()
			testedJob.AfterJobStopped()

			Expect(sampleJob).To(PointTo(
				MatchAllFields(Fields{
					"DoOk":           BeTrue(),
					"BeforeStartOk":  BeTrue(),
					"BeforeStopOk":   BeTrue(),
					"AfterStoppedOk": BeTrue(),
				}),
			))
		})
	})

	Context(`Panic if there is no "Do() error" method`, func() {
		It("Should panic", func() {
			Expect(func() {
				AutoJob(new(failedJob))
			}).To(Panic())
		})
	})
})

type failedJob int

func (j *failedJob) Do(v int) error { return nil }

type simpleJob int

func (j *simpleJob) Do() error {
	*j = 1
	return nil
}

type fullJob struct {
	DoOk           bool
	BeforeStartOk  bool
	BeforeStopOk   bool
	AfterStoppedOk bool
}

func (j *fullJob) Do() error {
	j.DoOk = true
	return nil
}
func (j *fullJob) BeforeJobStart() {
	j.BeforeStartOk = true
}
func (j *fullJob) BeforeJobStop() {
	j.BeforeStopOk = true
}
func (j *fullJob) AfterJobStopped() {
	j.AfterStoppedOk = true
}
