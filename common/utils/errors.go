//
// Error Handler
//
// PanicToError() could be used to transfer a panic code to customized error object
//
// PanicToSimpleError() could be used to transfer a panic code to simple error object
//
// 	func your_func() (err error) {
// 		defer PanicToSimpleError(&err)()
//
// 		// Code may cause panic...
// 	}
//
// You can wrap a function to a function with error returned:
//
// 	func your_func(int) { /**/ }
//
// 	errFunc := PanicToSimpleErrorWrapper(
// 		func() {
// 			your_func(var_1)
// 		},
// 	)
//
// 	err := errFunc()
//
// Capture panic to special handler
//
// If a go routine gets panic, the process would be terminated by default,
// you could use "BuildPanicCapture" to prevent the behavior.
//
// 	panicFreeFunc := utils.BuildPanicCapture(
//		func() {
//			/* Your code... */
//		},
//		func(panicObject interface{}) {
//			/* Gets executed if there is non-nil panic while executing your code */
//		},
// 	)
//
// 	go panicFreeFunc()
//
// Alternatively, there is a "BuildPanicToError" to set error object into "error pointer" if something gets panic.
//
// Stack error
//
// If you are writing IoC framework, the code position in framework may be too complex while output panic message,
// you could use "StackError" to record caller's position:
//
//	func yourLibrary() error {
//		err := someFunc()
//		if err != nil {
//			return BuildErrorWithCaller(err)
//		}
//
//		return nil
//	}
//
// For lambda-style programming, you could use "common/runtime.GetCallerInfo()" to keep the position of caller:
//
// 	func yourIoc(clientCode func()) {
// 		callerInfo := runtime.GetCallerInfo()
//
// 		defer func() {
// 			p := recover()
// 			if p != nil {
// 				panic(BuildErrorWithCallerInfo(SimpleErrorConverter(p), callerInfo))
// 			}
// 		}
//
// 		clinetCode()
// 	}
package utils

import (
	"fmt"
	gr "github.com/Cepave/open-falcon-backend/common/runtime"
)

// Defines the converter from converting any value to error object
type ErrorConverter func(p interface{}) error

// Defines the handler of panic
type PanicHandler func(interface{})

// Converts the panic content to error object
//
// 	err - The holder of error object
// 	errConverter - The builder function for converting non-error object to error object
func PanicToError(err *error, errConverter ErrorConverter) func() {
	return func() {
		p := recover()
		if p == nil {
			return
		}

		*err = errConverter(p)
	}
}

// Converts the panic content to error object by SimpleErrorConverter()
func PanicToSimpleError(err *error) func() {
	return PanicToError(err, SimpleErrorConverter)
}

// Simple converter for converting non-error object to error object by:
//
// 	fmt.Errorf("%v", object)
func SimpleErrorConverter(p interface{}) error {
	if errObject, ok := p.(error); ok {
		return errObject
	}

	return fmt.Errorf("%v", p)
}

// Convert a lambda function to function with error returned
func PanicToErrorWrapper(mainFunc func(), errConverter ErrorConverter) func() error {
	return func() (err error) {
		defer PanicToError(&err, errConverter)
		mainFunc()
		return
	}
}

// Convert a lambda function to function with error returned(by SimpleErrorConverter)
func PanicToSimpleErrorWrapper(mainFunc func()) func() error {
	return PanicToErrorWrapper(mainFunc, SimpleErrorConverter)
}

// Builds simple function, which executes target function with panic handler(panic-free)
func BuildPanicCapture(targetFunc func(), panicHandler func(interface{})) func() {
	return func() {
		defer func() {
			p := recover()
			if p != nil {
				panicHandler(p)
			}
		}()

		targetFunc()
	}
}

// Builds sample function, which captures panic object and convert it to error object
func BuildPanicToError(targetFunc func(), errHolder *error) func() {
	return func() {
		defer func() {
			p := recover()
			if p != nil {
				*errHolder = fmt.Errorf("%v", p)
			}
		}()

		targetFunc()
	}
}

type StackError struct {
	cause      error
	callerInfo *gr.CallerInfo
}

func (e *StackError) Error() string {
	return fmt.Sprintf("%s:%d:%v", e.callerInfo.GetFile(), e.callerInfo.Line, e.cause)
}

func DeferCatchPanicWithCaller() func() {
	callerInfo := gr.GetCallerInfoWithDepth(1)

	return func() {
		p := recover()
		if p == nil {
			return
		}

		panic(BuildErrorWithCallerInfo(
			SimpleErrorConverter(p), callerInfo,
		))
	}
}

func BuildErrorWithCaller(err error) *StackError {
	if err == nil {
		return nil
	}

	return BuildErrorWithCallerDepth(err, 2)
}
func BuildErrorWithCallerDepth(err error, depth int) *StackError {
	if err == nil {
		return nil
	}

	return BuildErrorWithCallerInfo(
		err,
		gr.GetCallerInfoWithDepth(depth),
	)
}

// Builds the error object with caller info, if the error object is type of *StackError,
// replaces the caller info with the new one
func BuildErrorWithCallerInfo(err error, callerInfo *gr.CallerInfo) *StackError {
	if err == nil {
		return nil
	}

	if stackError, ok := err.(*StackError); ok {
		stackError.callerInfo = callerInfo
		return stackError
	}

	return &StackError{
		cause:      err,
		callerInfo: callerInfo,
	}
}
