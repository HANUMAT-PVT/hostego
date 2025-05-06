package logger

import (
	"time"

	"github.com/sirupsen/logrus"
)

type MyLogger struct {
	*logrus.Logger
}

var logger = logrus.New()
var myLogger = &MyLogger{logger}

func setFormat(format string) {
	switch format {
	case "json":
		setupJsonLogging()
	default:
		setupTextLogging()
	}
}

func setLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		myLogger.Fatal(err)
	}
	myLogger.SetLevel(lvl)
}
func (ml *MyLogger) InfoWithStruct(message string, data interface{}) {
	fields := logrus.Fields{
		"data": data,
	}
	ml.WithFields(fields).Info(message)
}

func (ml *MyLogger) WarnWithStruct(message string, data interface{}) {
	fields := logrus.Fields{
		"data": data,
	}
	ml.WithFields(fields).Warn(message)
}

func (ml *MyLogger) ErrorWithStruct(message string, data interface{}) {
	fields := logrus.Fields{
		"data": data,
	}
	ml.WithFields(fields).Error(message)
}

// Provides logger
func GetLogger() *MyLogger {
	return myLogger
}

func GetLoggerWithFormatAndLogging(loggingFormat string, loggingLevel string) *MyLogger {
	setFormat(loggingFormat)
	setLevel(loggingLevel)
	return GetLogger()
}

func setupJsonLogging() {
	myLogger.Logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "levelname",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}) // suitable for prod env
	myLogger.Logger.SetReportCaller(true)
}

func setupTextLogging() {
	myLogger.Logger.SetFormatter(&logrus.TextFormatter{})                   // suitable for local env
	myLogger.Logger.SetFormatter(&logrus.TextFormatter{PadLevelText: true}) // by adding padding to the level text.
	myLogger.Logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	myLogger.Logger.SetReportCaller(true) // to include caller name
}
