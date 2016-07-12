package logruslog

import (
	log "github.com/Sirupsen/logrus"
	"github.com/chyeh/viper"
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
	logLevelStr := viper.GetString("logLevel")
	log.SetLevel(logLevel(logLevelStr))
}
