package g

import (
	"runtime"
)

const (
	VERSION = "2.0.2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
