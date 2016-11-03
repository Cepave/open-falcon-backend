package logruslog

import (
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	log "github.com/Sirupsen/logrus"
	"strings"
)

const RootLoggerName = "ROOT"
const DefaultLoggerLevel = "INFO"

// Mapping between logger name and level
type LoggerLevelMapping map[string]string

type LoggerFactory struct {
	LoggerLevels LoggerLevelMapping
}

// Gets logger with pre-configured level
func (factory *LoggerFactory) GetLogger(name string) *log.Logger {
	level, ok := factory.LoggerLevels[name]
	if ok {
		return NewDefaultLogger(level)
	}

	rootLevel, rootOk := factory.LoggerLevels[RootLoggerName]
	if rootOk {
		return NewDefaultLogger(rootLevel)
	}

	return NewDefaultLogger(DefaultLoggerLevel)
}

// Initialize a logger with string value of level
func NewDefaultLogger(logLevelValue string) *log.Logger {
	newLogger := log.New()
	newLogger.Formatter = newDefulatFormatter()
	newLogger.Level = logLevel(logLevelValue)

	return newLogger
}

// Sets log level by various string
//
// "debug", "Debug", "DEBUG"
// "info", "Info", "INFO"
// "warn", "Warn", "WARN"
// "error", "Error", "ERROR"
// "fatal", "Fatal", "FATAL"
// "panic", "Panic", "PANIC"
//
// Otherwise, use **INFO** level
func SetLogLevelByString(logLevelValue string) {
	log.SetFormatter(newDefulatFormatter())
	log.SetLevel(logLevel(logLevelValue))
}

func newDefulatFormatter() *log.TextFormatter {
	return &log.TextFormatter{FullTimestamp: true}
}

func logLevel(l string) log.Level {
	switch strings.ToLower(l) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func Init() {
	SetLogLevelByString(vipercfg.Config().GetString("logLevel"))
}
