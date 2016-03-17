package g

import (
	"runtime"
)

// change log
// 1.3.2 add last api for querying last item
// 1.3.3 rm debug log in http.graph
// 1.3.4 add http-api /graph/last/raw
// 1.3.5 fill response with endpoint & counter when rpc Graph.Query getting errors
// 1.4.0 restruct query: use simple rpc conn pool
// 1.4.1 add last item counter, add proc for connpool
// 1.4.2 rm nil items in http.responses
// 1.4.3 spell check, make config consistent with previous

const (
	VERSION = "1.4.3"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
