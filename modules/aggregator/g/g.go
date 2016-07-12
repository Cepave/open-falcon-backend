package g

import (
	"runtime"
)

const (
	VERSION = "0.0.4"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
