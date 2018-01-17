package log

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type stackHook struct{}

func (hook stackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook stackHook) Fire(entry *logrus.Entry) error {
	var skipFrames int
	if len(entry.Data) == 0 {
		// When WithField(s) is not used, we have 8 logrus frames to skip.
		skipFrames = 8
	} else {
		// When WithField(s) is used, we have 6 logrus frames to skip.
		skipFrames = 6
	}
	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(skipFrames, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "common/log") && !strings.Contains(name, "sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			f := struct {
				file string
				line int
				name string
			}{
				file: file,
				line: line,
				name: name,
			}
			entry.Data["func"] = f
			break
		}
	}
	return nil
}
