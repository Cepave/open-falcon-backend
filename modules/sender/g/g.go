package g

import (
	"runtime"
)

const (
	VERSION = "0.0.1"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
