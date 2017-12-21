package cron

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe(`Testing on "IntervalService"`, func() {
	Context("Start() / Stop() functions", func() {
		var (
			testedService *IntervalService
		)

		JustBeforeEach(func() {
			testedService.Start()
		})

		AfterEach(func() {
			testedService.Stop()
			testedService = nil
		})

		Context("Multiple jobs", func() {
			var (
				job1 = 0
				job2 = 0

				jobFunc1 = ProcJob(func() error {
					job1++
					return nil
				})
				jobFunc2 = ProcJob(func() error {
					job2++
					return nil
				})
			)

			BeforeEach(func() {
				testedService = NewIntervalService()

				testedService.Add("job-1", &IntervalConfig{FixedDelay: 100 * time.Millisecond}, jobFunc1)
				testedService.Add("job-2", &IntervalConfig{FixedDelay: 100 * time.Millisecond}, jobFunc2)
			})

			It("Both of Job-1 and Job-2 execute 3 times", func() {
				Eventually(
					func() []int { return []int{job1, job2} },
					2*time.Second, 150*time.Millisecond,
				).Should(ConsistOf(BeNumerically(">=", 3), BeNumerically(">=", 3)))

				testedService.Stop()

				currentJob1 := job1
				currentJob2 := job2

				Consistently(
					func() []int { return []int{job1, job2} },
					time.Second, 200*time.Millisecond,
				).Should(
					ConsistOf(Equal(currentJob1), Equal(currentJob2)),
					"The job should not be executed after the service is stopped",
				)
			})
		})

		Context("Multiple calling for Start()/Stop() function", func() {
			BeforeEach(func() {
				testedService = NewIntervalService()
			})

			Context("Start() function", func() {
				It("The flag should be true and nothing failed", func() {
					testedService.Start()
					testedService.Start()
					testedService.Start()

					Expect(testedService.started).To(BeTrue())
				})
			})

			Context("Stop() function", func() {
				It("The flag should be false and nothing failed", func() {
					Expect(testedService.started).To(BeTrue())

					testedService.Stop()
					testedService.Stop()
					testedService.Stop()
					Expect(testedService.started).To(BeFalse())
				})
			})
		})
	})

	Context("StartWorker() / StopWorker() functions", func() {
		var testedService *IntervalService

		AfterEach(func() {
			testedService.Stop()
			testedService = nil
		})

		Context("Start a worker, and stop it", func() {
			var (
				jobValue    = 0
				getJobValue = func() int {
					return jobValue
				}
			)

			BeforeEach(func() {
				testedService = NewIntervalService()
				testedService.Add(
					"ww-1", &IntervalConfig{FixedDelay: 100 * time.Millisecond},
					ProcJob(func() error {
						jobValue++
						return nil
					}),
				)
			})

			It("After stopped, the value should not be changed", func() {
				By("Start the job")
				Expect(testedService.StartWorker("ww-1")).To(BeTrue())

				By("Start the job again(un-effective)")
				Expect(testedService.StartWorker("ww-1")).To(BeFalse())

				Eventually(
					getJobValue, 2*time.Second, 150*time.Millisecond,
				).Should(BeNumerically(">=", 3))

				By("Stop the job")
				Expect(testedService.StopWorker("ww-1")).To(BeTrue())

				By("Stop the job again(un-effective)")
				Expect(testedService.StopWorker("ww-1")).To(BeFalse())

				lastJobValue := jobValue

				Consistently(
					getJobValue, time.Second, 150*time.Millisecond,
				).Should(Equal(lastJobValue))
			})

			It("Start non-existing job", func() {
				Expect(testedService.StartWorker("ww-3")).To(BeFalse())
			})

			It("Stop non-existing job", func() {
				Expect(testedService.StartWorker("ww-3")).To(BeFalse())
			})
		})
	})

	Context("Get profile of a worker", func() {
		testedService := NewIntervalService()

		const workerName = "pj-01"
		testedService.Add(
			workerName,
			&IntervalConfig{
				FixedDelay: 200 * time.Millisecond,
			},
			ProcJob(func() error {
				return nil
			}),
		)

		It("The unstarted job should have no profile", func() {
			Expect(testedService.GetWorkerProfile(workerName)).To(BeNil())
		})

		It("The started job should have profile", func() {
			By("Start service")
			testedService.Start()

			Eventually(
				func() *JobProfile { return testedService.GetWorkerProfile(workerName) },
				2*time.Second, 200*time.Millisecond,
			).Should(PointTo(
				MatchAllFields(Fields{
					"Called": BeNumerically(">=", 5),
					"Failed": Equal(0),
				}),
			))
		})

		It("The stopped job should have no profile", func() {
			By("Stop service")
			testedService.Stop()

			Expect(testedService.GetWorkerProfile(workerName)).To(BeNil())
		})
	})
})

var _ = Describe(`Testing on "intervalWorker"`, func() {
	var (
		testedWorker *intervalWorker
	)

	JustBeforeEach(func() {
		if testedWorker == nil {
			return
		}

		testedWorker.start()
	})
	AfterEach(func() {
		if testedWorker == nil {
			return
		}

		testedWorker.stop()
		testedWorker = nil
	})

	Context(`Effective of "InitialDelay"`, func() {
		const initialDelay = 200 * time.Millisecond
		var beforeStart time.Time
		var times []time.Time

		BeforeEach(func() {
			beforeStart = time.Now()
			times = make([]time.Time, 0, 4)
			testedWorker = newIntervalWorker(
				&IntervalConfig{
					InitialDelay: initialDelay,
					FixedDelay:   500 * time.Millisecond,
					job: ProcJob(func() error {
						times = append(times, time.Now())
						return nil
					}),
				},
			)
		})

		It("The elapsed time between the callings on job should be equal or more than the initial delay time", func() {
			Eventually(
				func() int { return len(times) },
				time.Second, 200*time.Millisecond,
			).Should(BeNumerically(">=", 1))

			Expect(times[0].Sub(beforeStart)).To(And(
				BeNumerically(">=", initialDelay),
				BeNumerically("<", initialDelay+(initialDelay>>1)),
			))
		})
	})

	Context("Worker executes and stop normally", func() {
		const fixedDelayValue = 200 * time.Millisecond

		var (
			times = make([]time.Time, 0, 8)
		)

		BeforeEach(func() {
			testedWorker = newIntervalWorker(
				&IntervalConfig{
					FixedDelay: fixedDelayValue,
					job: ProcJob(func() error {
						times = append(times, time.Now())
						return nil
					}),
				},
			)
		})

		It("The elapsed time between the callings on job should be equal or more than the delay time", func() {
			/**
			 * Waiting for at least 5 times of calling on job
			 */
			Eventually(
				func() int { return len(times) },
				2*time.Second, 250*time.Millisecond,
			).Should(BeNumerically(">=", 5))
			// :~)

			testedWorker.stop()

			/**
			 * Asserts the stopped behaviour
			 */
			By("Assert stopping of worker")
			lastTimes := len(times)
			Consistently(
				func() int { return len(times) },
				time.Second, 200*time.Millisecond,
			).Should(Equal(lastTimes), "The times should remain as [%d] after the worker has stopped", lastTimes)
			// :~)

			/**
			 * Asserts the [elapsed time] <= T < [elapsed time * 1.5]
			 * and the initialize delay
			 */
			for i, v := range times {
				if i == 0 {
					continue
				}

				Expect(v.Sub(times[i-1])).To(And(
					BeNumerically(">=", fixedDelayValue),
					BeNumerically("<", fixedDelayValue+(fixedDelayValue>>1)),
				))
			}
			// :~)
		})
	})

	Context("Worker executes with error", func() {
		var times []time.Time

		BeforeEach(func() {
			times = make([]time.Time, 0, 8)
		})

		DescribeTable("The delay time should be as expected",
			func(fixedDelayMillis int, errorDelayMillis int, atLeastMillis int) {
				testedWorker := newIntervalWorker(
					&IntervalConfig{
						FixedDelay: time.Duration(fixedDelayMillis) * time.Millisecond,
						ErrorDelay: time.Duration(errorDelayMillis) * time.Millisecond,
						job: ProcJob(func() error {
							times = append(times, time.Now())
							return errors.New("Sample error")
						}),
					},
				)

				testedWorker.start()

				atLeastValue := time.Duration(atLeastMillis) * time.Millisecond

				/**
				 * Waiting for at least 5 times of calling on job
				 */
				Eventually(
					func() int { return len(times) },
					2*time.Second+5*atLeastValue, 250*time.Millisecond,
				).Should(BeNumerically(">=", 5))
				// :~)

				testedWorker.stop()

				/**
				 * Asserts the [elapsed time] <= T < [elapsed time * 1.5]
				 */
				for i := 1; i < len(times); i++ {
					if i == 0 {
						continue
					}

					Expect(times[i].Sub(times[i-1])).To(And(
						BeNumerically(">=", atLeastValue),
						BeNumerically("<", atLeastValue+(atLeastValue>>1)),
					))
				}
				// :~)
			},
			Entry("Error with 250ms delay", 0, 250, 250),
			Entry("Error with 0 milliseconds delay(use fixed delay(500ms) instead)", 500, 0, 500),
		)
	})

	Context("Execute start()/stop() functions of worker multiple times", func() {
		BeforeEach(func() {
			testedWorker = newIntervalWorker(
				&IntervalConfig{
					InitialDelay: 5 * time.Second,
					job:          ProcJob(func() error { return nil }),
				},
			)
		})

		Context("Start() function", func() {
			It("The flag should be true and nothing failed", func() {
				testedWorker.start()
				testedWorker.start()
				testedWorker.start()

				Expect(testedWorker.started).To(BeTrue())
			})
		})
		Context("Stop() function", func() {
			It("The flag should be false and nothing failed", func() {
				Expect(testedWorker.started).To(BeTrue())

				testedWorker.stop()
				testedWorker.stop()
				testedWorker.stop()
				Expect(testedWorker.started).To(BeFalse())
			})
		})
	})

	Context("Calling of Lifecycle", func() {
		var beforeCalled, beforeStop, afterCalled = false, false, false

		localWorker := newIntervalWorker(
			&IntervalConfig{
				FixedDelay: 5 * time.Second,
				job: &JobInstance{
					BeforeJobStartImpl: func() {
						beforeCalled = true
					},
					BeforeJobStopImpl: func() {
						beforeStop = true
					},
					AfterJobStoppedImpl: func() {
						afterCalled = true
					},
					DoImpl: func() error { return nil },
				},
			},
		)

		It("Functions of life-cycle should get called", func() {
			By("Start the worker - BeforeJobStart() should get called")
			localWorker.start()
			Expect(beforeCalled).To(BeTrue())

			By("Stop the worker")
			localWorker.stop()

			Eventually(
				func() bool { return beforeStop },
				time.Second, 250*time.Millisecond,
			).Should(BeTrue(), "BeforeJobStop() should get called")

			Expect(afterCalled).To(BeTrue())
		})
	})

	Context("The information of profile", func() {
		BeforeEach(func() {
			flip := true

			testedWorker = newIntervalWorker(
				&IntervalConfig{
					FixedDelay: 200 * time.Millisecond,
					ErrorDelay: 200 * time.Millisecond,
					job: ProcJob(func() error {
						if flip {
							flip = false
							return errors.New("Get error")
						}

						flip = true
						return nil
					}),
				},
			)
		})

		It("The times of Called should be more than the times of failed", func() {
			Eventually(
				func() int { return testedWorker.times },
				2*time.Second, 200*time.Millisecond,
			).Should(BeNumerically(">=", 4))

			currentProfile := testedWorker.getProfile()

			GinkgoT().Logf("Profile: %#v", currentProfile)
			Expect(currentProfile).To(PointTo(
				MatchAllFields(Fields{
					"Called": BeNumerically(">=", 4),
					"Failed": And(
						BeNumerically(">=", currentProfile.Called>>1),
						BeNumerically("<", currentProfile.Called),
					),
				}),
			))
		})
	})
})
