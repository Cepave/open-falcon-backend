package g

import (
	"runtime"
)

const (
	VERSION = "0.0.2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
