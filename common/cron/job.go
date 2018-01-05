package cron

import (
	"fmt"
	"reflect"
)

// Defines the events of life cycle on a job
type JobLifecycle interface {
	// Gets called before the job starts
	BeforeJobStart()
	// Gets called before the job is going to be stopped
	//
	// This function get called asynchronously.
	BeforeJobStop()
	// Gets called after the job has been stopped
	AfterJobStopped()
}

// Defines all of the actions of a job
type Job interface {
	JobLifecycle
	// Main work to be executed of a job
	Do() error
}

// Represents the profile of job.
type JobProfile struct {
	// The times of the job get executed
	Called int
	// The times of the job's executions are failed
	Failed int
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Converts any object to Job object,
// this function would check the signature of "Job" interface.
//
// This function "Do() error" is **mandatory**.
//
// Other methods of life-cycle are optional.
func AutoJob(v interface{}) Job {
	valueOfImpl := reflect.ValueOf(v)

	/**
	 * Load function of main works of job
	 */
	doFunc := getValidFunction(
		valueOfImpl.MethodByName("Do"),
		nil, []reflect.Type{errorType},
	)
	if !doFunc.IsValid() {
		panic(fmt.Sprintf(`Object %v doesn't have "Do() error" function`, v))
	}
	// :~)

	/**
	 * Loading methods of lifecycle
	 */
	beforeStartFunc := getValidFunction(
		valueOfImpl.MethodByName("BeforeJobStart"),
		nil, nil,
	)
	beforeStopFunc := getValidFunction(
		valueOfImpl.MethodByName("BeforeJobStop"),
		nil, nil,
	)
	afterStoppedFunc := getValidFunction(
		valueOfImpl.MethodByName("AfterJobStopped"),
		nil, nil,
	)
	// :~)

	doImpl := doFunc.Interface().(func() error)
	beforeStartImpl := func() {}
	if beforeStartFunc.IsValid() {
		beforeStartImpl = beforeStartFunc.Interface().(func())
	}
	beforeStopImpl := func() {}
	if beforeStopFunc.IsValid() {
		beforeStopImpl = beforeStopFunc.Interface().(func())
	}
	afterStoppedImpl := func() {}
	if afterStoppedFunc.IsValid() {
		afterStoppedImpl = afterStoppedFunc.Interface().(func())
	}

	return &JobInstance{
		BeforeJobStartImpl:  beforeStartImpl,
		BeforeJobStopImpl:   beforeStopImpl,
		AfterJobStoppedImpl: afterStoppedImpl,
		DoImpl:              doImpl,
	}
}

// Constructs a job, which does nothing.
func EmptyJob() Job {
	return Job(new(emptyJob))
}

// As struct-typed of a "Job"
type JobInstance struct {
	BeforeJobStartImpl  func()
	BeforeJobStopImpl   func()
	AfterJobStoppedImpl func()
	DoImpl              func() error
}

// Delegates to "obj.BeforeJobStartImpl". This is optional.
func (j *JobInstance) BeforeJobStart() {
	if j.BeforeJobStartImpl != nil {
		j.BeforeJobStartImpl()
	}
}

// Delegates to "obj.BeforeJobStopImpl". This is optional.
func (j *JobInstance) BeforeJobStop() {
	if j.BeforeJobStopImpl != nil {
		j.BeforeJobStopImpl()
	}
}

// Delegates to "obj.AfterJobStoppedImpl". This is optional.
func (j *JobInstance) AfterJobStopped() {
	if j.AfterJobStoppedImpl != nil {
		j.AfterJobStoppedImpl()
	}
}

// Delegates to "obj.DoImpl". This is **mandatory**.
func (j *JobInstance) Do() error {
	return j.DoImpl()
}

// As function-typed of a "job"
type ProcJob func() error

// Implementation of "Job" interface. Do nothing.
func (p ProcJob) BeforeJobStart() {}

// Implementation of "Job" interface. Do nothing.
func (p ProcJob) BeforeJobStop() {}

// Implementation of "Job" interface. Do nothing.
func (p ProcJob) AfterJobStopped() {}

// Delegates to the body of this function type.
func (p ProcJob) Do() error {
	return p()
}

type emptyJob bool

func (b *emptyJob) BeforeJobStart()  {}
func (b *emptyJob) BeforeJobStop()   {}
func (b *emptyJob) AfterJobStopped() {}
func (b *emptyJob) Do() error        { return nil }

func getValidFunction(funcValue reflect.Value, inTypes []reflect.Type, outTypes []reflect.Type) reflect.Value {
	zeroValue := reflect.ValueOf(nil)

	if !funcValue.IsValid() {
		return zeroValue
	}

	funcType := funcValue.Type()

	/**
	 * Checks the input parameters of a function
	 */
	if funcType.NumIn() != len(inTypes) {
		return zeroValue
	}
	for i, inType := range inTypes {
		if funcType.In(i) != inType {
			return zeroValue
		}
	}
	// :~)

	/**
	 * Checks the returned types of a function
	 */
	if funcType.NumOut() != len(outTypes) {
		return zeroValue
	}
	for i, outType := range outTypes {
		if funcType.Out(i) != outType {
			return zeroValue
		}
	}
	// :~)

	return funcValue
}
