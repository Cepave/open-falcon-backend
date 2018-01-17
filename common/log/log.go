package log

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

var loggers = make(map[string]*Logger)

type Logger struct {
	name string
	ll   *logrus.Logger
}

func NewLogger(name string) *Logger {
	l, ok := loggers[name]
	if !ok {
		l = &Logger{
			name: name,
			ll:   logrus.New(),
		}
		loggers[name] = l
		l.ll.AddHook(stackHook{})
		l.ll.AddHook(nameHook{name: name})
		//l.ll.WithField("module", name)
		return l
	}
	return l
}

func GetLogger(name string) *Logger {
	return loggers[name]
}

func ListAll() map[string]*Logger {
	return loggers
}

func SetLevel(name string, lv logrus.Level) {
	for k, v := range loggers {
		if strings.HasPrefix(k, name) {
			v.ll.SetLevel(lv)
		}
	}
}

func (l *Logger) NewExtLogger(name string) *Logger {
	ext := NewLogger(l.name + "." + name)
	ext.ll.SetLevel(l.Level())
	return ext
}

func (l *Logger) Level() logrus.Level {
	return l.ll.Level
}

func (l *Logger) SetLevel(lv logrus.Level) {
	SetLevel(l.name, lv)
}

func (l *Logger) Name() string {
	return l.name
}

// Metohds from logrus.Logger

func (l *Logger) Debug(args ...interface{}) {
	l.ll.Debug(args)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.ll.Debugf(format)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.ll.Debugln(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.ll.Error(args)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.ll.Errorf(format)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.ll.Errorln(args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.ll.Fatal(args)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.ll.Fatalf(format)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.ll.Fatalln(args)
}

func (l *Logger) Info(args ...interface{}) {
	l.ll.Info(args)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.ll.Infof(format)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.ll.Infoln(args)
}

func (l *Logger) Panic(args ...interface{}) {
	l.ll.Panic(args)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.ll.Panicf(format)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.ll.Panicln(args)
}

func (l *Logger) Print(args ...interface{}) {
	l.ll.Print(args)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.ll.Printf(format)
}

func (l *Logger) Println(args ...interface{}) {
	l.ll.Println(args)
}

func (l *Logger) SetNoLock() {
	l.ll.SetNoLock()
}

func (l *Logger) Warn(args ...interface{}) {
	l.ll.Warn(args)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.ll.Warnf(format, args)
}

func (l *Logger) Warning(args ...interface{}) {
	l.ll.Warning(args)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.ll.Warningf(format, args)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.ll.Warningln(args)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.ll.Warnln(args)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.ll.WithError(err)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.ll.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.ll.WithFields(fields)
}

func (l *Logger) Writer() *io.PipeWriter {
	return l.ll.Writer()
}

func (l *Logger) WriterLevel(level logrus.Level) *io.PipeWriter {
	return l.ll.WriterLevel(level)
}
