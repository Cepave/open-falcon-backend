package log

import (
	"github.com/sirupsen/logrus"
)

type nameHook struct {
	name string
}

func (hook nameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook nameHook) Fire(entry *logrus.Entry) error {
	entry.Data["module"] = hook.name
	return nil
}
