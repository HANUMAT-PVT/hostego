package logger

import (
	"context"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger       *logrus.Logger
	ctx          context.Context
	fields       map[string]interface{}
	customFields map[string]interface{}
}

func getLogrusLevel(level string) logrus.Level {
	switch level {
	case LogLevelDebug:
		return logrus.DebugLevel
	case LogLevelInfo:
		return logrus.InfoLevel
	case LogLevelWarn:
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel // defaults to InfoLevel if LOG_LEVEL is not set or invalid
	}
}

func NewLogrusLogger(ctx context.Context, level string) *LogrusLogger {
	envLevel := strings.ToLower(os.Getenv(LogLevel))
	if level != "" {
		envLevel = level
	}

	logrusLevel := getLogrusLevel(envLevel)
	logger := logrus.New()
	logger.SetLevel(logrusLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFunc:        "func",
			logrus.FieldKeyLevel:       "level",
			logrus.FieldKeyMsg:         "msg",
			logrus.FieldKeyLogrusError: "error",
			logrus.FieldKeyTime:        "timestamp",
		},
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // ISO8601 format
	})
	logger.SetReportCaller(false) // we're using fileAndFuncInfo() for actual log location

	m := make(map[string]interface{})
	c := make(map[string]interface{})

	return &LogrusLogger{logger: logger, ctx: ctx, fields: m, customFields: c}
}

func (l *LogrusLogger) AddFields(fields map[string]interface{}) {
	l.addCustomFields(fields)
}

func (l *LogrusLogger) Debug(msg string) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Debug(msg)
}

func (l *LogrusLogger) Info(msg string) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Info(msg)
}

func (l *LogrusLogger) Warn(msg string) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Warn(msg)
}

func (l *LogrusLogger) Error(msg string) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Error(msg)
}

func (l *LogrusLogger) Fatal(msg string) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Fatal(msg)
}

func (l *LogrusLogger) Infof(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Infof(msg, args...)
}

func (l *LogrusLogger) Warnf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Warnf(msg, args...)

}

func (l *LogrusLogger) Errorf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Errorf(msg, args...)

}

func (l *LogrusLogger) Fatalf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Fatalf(msg, args...)
}

func (l *LogrusLogger) Debugf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	l.logger.WithFields(l.fields).Debugf(msg, args...)

}

func (l *LogrusLogger) DebugCf(msg string) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Debug(msg)
}

func (l *LogrusLogger) InfoCf(msg string) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Info(msg)
}

func (l *LogrusLogger) WarnCf(msg string) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Warn(msg)
}

func (l *LogrusLogger) ErrorCf(msg string) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Error(msg)
}

func (l *LogrusLogger) FatalCf(msg string) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Fatal(msg)
}

func (l *LogrusLogger) InfofCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Infof(msg, args...)

}

func (l *LogrusLogger) WarnfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Warnf(msg, args...)

}

func (l *LogrusLogger) ErrorfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Errorf(msg, args...)

}

func (l *LogrusLogger) FatalfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Fatalf(msg, args...)

}

func (l *LogrusLogger) DebugfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	l.logger.WithFields(l.customFields).Debugf(msg, args...)

}

func (l *LogrusLogger) addContextCommonFields(fields map[string]interface{}) {
	if l.ctx != nil && l.ctx.Value(CommonFieldsKey) != nil {
		for k, v := range l.ctx.Value(CommonFieldsKey).(map[string]interface{}) {
			if _, ok := fields[k]; !ok {
				fields[k] = v
			}
		}
	}

	fields["func"], fields["file"] = fileAndFuncInfo(FileAndInfoSkipLevel)

}

func (l *LogrusLogger) addCustomFields(fields map[string]interface{}) {
	l.customFields = fields
}
