package logruslog

import (
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	log "github.com/Sirupsen/logrus"
)

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
}
