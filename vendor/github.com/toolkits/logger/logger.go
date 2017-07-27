package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	level       int
	traceLogger = log.New(os.Stdout, "[T] ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "[D] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "[I] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stdout, "[W] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stdout, "[E] ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(os.Stdout, "[F] ", log.Ldate|log.Ltime|log.Lshortfile)
)

func SetLevelWithDefault(lv, defaultLv string) {
	err := SetLevel(lv)
	if err != nil {
		SetLevel(defaultLv)
	}
}

func SetLevel(lv string) error {
	if lv == "" {
		return fmt.Errorf("log level is blank")
	}

	l := strings.ToUpper(lv)

	switch l[0] {
	case 'T':
		level = 0
	case 'D':
		level = 1
	case 'I':
		level = 2
	case 'W':
		level = 3
	case 'E':
		level = 4
	case 'F':
		level = 5
	default:
		level = 6
	}

	if level == 6 {
		return fmt.Errorf("log level setting error")
	}

	return nil
}

func Trace(format string, v ...interface{}) {
	if 0 >= level {
		traceLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debug(format string, v ...interface{}) {
	if 1 >= level {
		debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	if 2 >= level {
		infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Warn(format string, v ...interface{}) {
	if 3 >= level {
		warnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	if 4 >= level {
		errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Fatal(format string, v ...interface{}) {
	if 5 >= level {
		fatalLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Traceln(v ...interface{}) {
	if 0 >= level {
		traceLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Debugln(v ...interface{}) {
	if 1 >= level {
		debugLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Infoln(v ...interface{}) {
	if 2 >= level {
		infoLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Warnln(v ...interface{}) {
	if 3 >= level {
		warnLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Errorln(v ...interface{}) {
	if 4 >= level {
		errorLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Fatalln(v ...interface{}) {
	if 5 >= level {
		fatalLogger.Output(2, fmt.Sprintln(v...))
	}
}
