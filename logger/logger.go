package logger

import (
	"github.com/Cepave/alarm/g"
	"github.com/Sirupsen/logrus"
	"os"
)

var logger *logrus.Entry = nil

func InitLogger() {
	config := g.Config()
	if logger == nil {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stderr)
		if config.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.WarnLevel)
		}
		logger = logrus.WithFields(logrus.Fields{})
	}
}

func Logger() *logrus.Entry {
	return logger
}
