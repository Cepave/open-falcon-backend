package utils

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
