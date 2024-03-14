package logger

import "fmt"

type Logger interface {
	InfoLog(s string)
	ErrorLog(s string)
	FatalLog(s string)
}

var l Logger

// InitLogger initializes logger
func InitLogger(logger Logger) {
	l = logger
}

// Info writes log with INFO prefix
func Info(format string, v ...any) {
	if len(v) == 0 {
		l.InfoLog(format)
	} else {
		l.InfoLog(fmt.Sprintf(format, v...))
	}
}

// Error writes log with ERROR prefix
func Error(format string, v ...any) {
	if len(v) == 0 {
		l.ErrorLog(format)
	} else {
		l.ErrorLog(fmt.Sprintf(format, v...))
	}
}

// Fatal writes log with FATAL prefix and makes os.Exit(1)
func Fatal(format string, v ...any) {
	if len(v) == 0 {
		l.FatalLog(format)
	} else {
		l.FatalLog(fmt.Sprintf(format, v...))
	}
}
