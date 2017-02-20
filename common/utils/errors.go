package utils

import (
	"fmt"
)

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
