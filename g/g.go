package g

import (
	"runtime"
)

const (
	VERSION = "0.0.3"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
