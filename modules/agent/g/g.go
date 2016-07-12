package g

import (
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
