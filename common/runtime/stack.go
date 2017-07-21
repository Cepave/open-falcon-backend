//
// Because go language does not have industrial level of exception handing mechanism,
// using the information of calling state is the only way to expose secret in code.
//
// Obtain CallerInfo
//
// CallerInfo is the main struct which holds essential information of detail on code.
//
// You can obtain CallerInfo by various functions:
//
// 	GetCallerInfo() - Obtains the caller(the previous calling point to current function)
// 	GetCallerInfoWithDepth() - Obtains the caller of caller by numeric depth
//
package runtime

import (
	"fmt"
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

	goPath = file[0 : size-suffix]
	prefixLength = len(goPath)
	// :~)
}

// Information of the position in file and where line number is targeted
type CallerInfo struct {
	Line int

	file    string
	rawFile string
}

type CallerStack []*CallerInfo

func (s CallerStack) AsStringStack() []string {
	callerStackString := make([]string, 0)
	for _, caller := range s {
		callerStackString = append(callerStackString, caller.String())
	}

	return callerStackString
}
func (s CallerStack) ConcatStringStack(sep string) string {
	return strings.Join(s.AsStringStack(), sep)
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
func (c *CallerInfo) String() string {
	return fmt.Sprintf("%s:%d", c.GetFile(), c.Line)
}

// Gets stack of caller info
func GetCallerInfoStack(startDepth int, endDepth int) CallerStack {
	callers := make([]*CallerInfo, 0)
	for i := startDepth + 1; i < endDepth+1; i++ {
		callerInfo := GetCallerInfoWithDepth(i)
		if callerInfo == nil {
			break
		}

		callers = append(callers, callerInfo)
	}

	return callers
}

// Gets caller info from current function
func GetCallerInfo() *CallerInfo {
	return GetCallerInfoWithDepth(1)
}

// Gets caller info with depth.
//
// N means the Nth caller of caller.
func GetCallerInfoWithDepth(depth int) *CallerInfo {
	_, file, line, ok := runtime.Caller(depth + 2)
	if !ok {
		return nil
	}

	return &CallerInfo{
		rawFile: file,
		Line:    line,
	}
}

func trimGoPath(filename string) string {
	if strings.HasPrefix(filename, goPath) {
		return filename[prefixLength:]
	}

	return filename
}
