// This package provides scheduled(by interval) service for jobs.
//
// IntervalService
//
// The main service used to manage/start/stop job you want to execute continuously.
//
// Time Interval
//
// As "IntervalConfig" defined, there are "duration"s for the period of job execution:
//
//	InitialDelay - The delay duration before the first execution
//	FixedDelay - The delay duration before the next execution(after the complete of previous execution)
//	ErrorDelay - The delay duration before the next execution(after the complete of previous execution which is failed)
//
// By default(all of the durations are 0), your job would be executed one by another.
//
// Job
//
// A "job" is a object implements "Job" interface. You could use "ProcJob" as functional paradigm.
//
// The "AutoJob()" is an out-of-box function to convert your partial-implementation object to a "Job" object.
// 	newJob := AutoJob(new(YourStruct))
//	service.Add("worker-1", &IntervalConfig{ }, newJob)
package cron

import (
	"fmt"
	"sync"
	"time"
)

// Configuration of interval jobs
type IntervalConfig struct {
	// The delay time before first execution
	InitialDelay time.Duration
	// The delay time between executions(after the job is finished)
	FixedDelay time.Duration
	// The delay time for next job if the last execution gave error
	//
	// If the value <= 0(no duration), the value of "FixedDelay" would be used instead.
	ErrorDelay time.Duration

	job Job
}

type IntervalService struct {
	namedJobs      map[string]*IntervalConfig
	startedWorkers map[string]*intervalWorker

	locker  *sync.Mutex
	started bool
}

func NewIntervalService() *IntervalService {
	return &IntervalService{
		namedJobs:      make(map[string]*IntervalConfig),
		startedWorkers: make(map[string]*intervalWorker),
		locker:         &sync.Mutex{},
		started:        false,
	}
}

// Adds a job with worker name.
//
// The same name of worker would be overrode.
func (s *IntervalService) Add(name string, config *IntervalConfig, targetJob Job) {
	newConfig := *config
	newConfig.job = targetJob

	s.namedJobs[name] = &newConfig
}

// Gets the profiling information after a job has started.
//
// Once the job is stopped, whatever by "Stop()" or "StopWorker()",
// the information of profiling would be abandoned.
//
// This method gives "nil" the job is not running.
func (s *IntervalService) GetWorkerProfile(name string) *JobProfile {
	worker, ok := s.startedWorkers[name]

	if !ok {
		return nil
	}

	return worker.getProfile()
}

// Starts all of the added jobs
func (s *IntervalService) Start() {
	s.locker.Lock()
	defer s.locker.Unlock()

	if s.started {
		return
	}

	for name, config := range s.namedJobs {
		/**
		 * Skip the workers have started
		 */
		if _, ok := s.startedWorkers[name]; ok {
			continue
		}
		// :~)

		s.startedWorkers[name] = newIntervalWorker(config)
		s.startedWorkers[name].start()
	}

	s.started = true
}

// Stops all of the added jobs
func (s *IntervalService) Stop() {
	s.locker.Lock()
	defer s.locker.Unlock()

	if !s.started {
		return
	}

	for _, worker := range s.startedWorkers {
		worker.stop()
	}

	s.startedWorkers = make(map[string]*intervalWorker)

	s.started = false
}

// Starts a worker by name
//
// This method returns "true" if the job is started effectively.
func (s *IntervalService) StartWorker(name string) bool {
	s.locker.Lock()
	defer s.locker.Unlock()

	_, ok := s.startedWorkers[name]
	if ok {
		return false
	}

	config, ok := s.namedJobs[name]
	if !ok {
		return false
	}

	s.startedWorkers[name] = newIntervalWorker(config)
	s.startedWorkers[name].start()

	return true
}

// Stops a worker by name
//
// This method returns "true" if the job is started effectively.
func (s *IntervalService) StopWorker(name string) bool {
	s.locker.Lock()
	defer s.locker.Unlock()

	worker, ok := s.startedWorkers[name]
	if !ok {
		return false
	}

	worker.stop()
	delete(s.startedWorkers, name)

	return true
}

type intervalWorker struct {
	*IntervalConfig

	locker  *sync.Mutex
	started bool

	stopChannel chan bool

	times       int
	failedTimes int
}

func newIntervalWorker(config *IntervalConfig) *intervalWorker {
	return &intervalWorker{
		IntervalConfig: config,
		locker:         &sync.Mutex{},
		stopChannel:    make(chan bool),
		started:        false,
	}
}

func (w *intervalWorker) getProfile() *JobProfile {
	return &JobProfile{
		Called: w.times,
		Failed: w.failedTimes,
	}
}
func (w *intervalWorker) start() {
	w.locker.Lock()
	defer w.locker.Unlock()

	if w.started {
		return
	}

	jobHandler := func() (err error) {
		defer func() {
			p := recover()
			if p != nil {
				err = fmt.Errorf("Job has panic!! %v", p)
				logger.Warnf("%v", err)
			}
		}()

		return w.job.Do()
	}

	w.job.BeforeJobStart()
	go func() {
		timer := time.NewTimer(w.InitialDelay)
		defer timer.Stop()

		for {
			select {
			case <-w.stopChannel:
				return
			case <-timer.C:
				w.times++
				err := jobHandler()

				if err != nil {
					w.failedTimes++

					if w.ErrorDelay > 0 {
						timer.Reset(w.ErrorDelay)
						continue
					}
				}

				timer.Reset(w.FixedDelay)
			}
		}
	}()

	w.started = true
}
func (w *intervalWorker) stop() {
	w.locker.Lock()
	defer w.locker.Unlock()

	if !w.started {
		return
	}

	go w.job.BeforeJobStop()
	w.stopChannel <- true
	w.started = false

	w.job.AfterJobStopped()
}
