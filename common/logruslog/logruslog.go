package logruslog

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	log "github.com/Sirupsen/logrus"
)

type funcInfo struct {
	file string
	line int
	name string
}

func (f funcInfo) String() string {
	return fmt.Sprintf("%s (%s:%d)", path.Base(f.name), path.Base(f.file), f.line)
}

type stackHook struct{}

func (hook stackHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook stackHook) Fire(entry *log.Entry) error {
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
		if !strings.Contains(name, "github.com/Sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			f := funcInfo{
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

func logLevel(l string) log.Level {
	switch l {
	case "debug", "Debug", "DEBUG":
		return log.DebugLevel
	case "info", "Info", "INFO":
		return log.InfoLevel
	case "warn", "Warn", "WARN":
		return log.WarnLevel
	case "error", "Error", "ERROR":
		return log.ErrorLevel
	case "fatal", "Fatal", "FATAL":
		return log.FatalLevel
	case "panic", "Panic", "PANIC":
		return log.PanicLevel
	default:
		return log.InfoLevel

	}
}

func Init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	logLevelStr := vipercfg.Config().GetString("logLevel")
	log.SetLevel(logLevel(logLevelStr))
	log.AddHook(stackHook{})
}
