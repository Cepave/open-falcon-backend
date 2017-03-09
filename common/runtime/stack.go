package runtime

import (
	"runtime"
	"strings"
)

var prefixLength int
var goPath string

func init() {
	_, file, _, _ := runtime.Caller(0)

	/**
	 * Figures out the $GOPATH of current file
	 */
	size := len(file)
	suffix := len("github.com/Cepave/open-falcon-backend/common/runtime/stack.go")

	goPath = file[0:size-suffix]
	prefixLength = len(goPath)
	// :~)
}

type CallerInfo struct {
	Line int

	file string
	rawFile string
}
// Gets the file by trimming of $GOPATH
//
// The processing of file is delayed to improve performance
func (c *CallerInfo) GetFile() string {
	if c.file != "" {
		return c.file
	}

	c.file = trimGoPath(c.rawFile)
	return c.file
}

func GetCallerInfo() *CallerInfo {
	return GetCallerInfoWithDepth(1)
}
func GetCallerInfoWithDepth(depth int) *CallerInfo {
	_, file, line, _ := runtime.Caller(depth + 2)

	return &CallerInfo {
		rawFile: file,
		Line: line,
	}
}

func trimGoPath(filename string) string {
	if strings.HasPrefix(filename, goPath) {
		return filename[prefixLength:]
	}

	return filename
}
