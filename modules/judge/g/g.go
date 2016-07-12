package g

import (
	"runtime"
)

// change log
// 2.0.1: bugfix HistoryData limit
// 2.0.2: clean stale data
const (
	VERSION = "2.0.2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
