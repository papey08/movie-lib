package logger

import (
	"io"
	"log"
)

// defaultLogger are 3 loggers of different levels
type defaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func (dl *defaultLogger) InfoLog(s string) {
	dl.infoLogger.Println(s)
}

func (dl *defaultLogger) ErrorLog(s string) {
	dl.errorLogger.Println(s)
}

func (dl *defaultLogger) FatalLog(s string) {
	dl.fatalLogger.Fatal(s)
}

// DefaultLogger initializes Logger writing logs to writer
func DefaultLogger(writer io.Writer) Logger {
	return &defaultLogger{
		infoLogger:  log.New(writer, "INFO: ", log.Ldate|log.Ltime),
		errorLogger: log.New(writer, "ERROR: ", log.Ldate|log.Ltime),
		fatalLogger: log.New(writer, "FATAL: ", log.Ldate|log.Ltime),
	}
}
