package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger       *zap.Logger
	ctx          context.Context
	fields       map[string]interface{}
	customFields map[string]interface{}
	mu           *sync.RWMutex
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case LogLevelDebug:
		return zapcore.DebugLevel
	case LogLevelInfo:
		return zapcore.InfoLevel
	case LogLevelWarn:
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel // defaults to InfoLevel if LOG_LEVEL is not set or invalid
	}
}

func NewZapLogger(ctx context.Context, level string) *ZapLogger {
	envLevel := strings.ToLower(os.Getenv(LogLevel))
	if level != "" {
		envLevel = level
	}

	zapLevel := getZapLevel(envLevel)
	atomicLevel := zap.NewAtomicLevelAt(zapLevel)
	config := zap.NewProductionConfig()
	config.Level = atomicLevel
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableCaller = true // we're using fileAndFuncInfo() for actual log location

	logger, _ := config.Build()

	m := make(map[string]interface{})
	c := make(map[string]interface{})

	return &ZapLogger{logger: logger, ctx: ctx, fields: m, mu: &sync.RWMutex{}, customFields: c}
}

func (l *ZapLogger) AddFields(fields map[string]interface{}) {
	l.addCustomFields(fields)
}

func (l *ZapLogger) Debug(msg string) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Debug(msg, zapFields...)
}

func (l *ZapLogger) Info(msg string) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Info(msg, zapFields...)
}

func (l *ZapLogger) Warn(msg string) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Warn(msg, zapFields...)
}

func (l *ZapLogger) Error(msg string) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Error(msg, zapFields...)
}

func (l *ZapLogger) Fatal(msg string) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Fatal(msg, zapFields...)
}

func (l *ZapLogger) Debugf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Debug(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) Infof(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Info(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) Warnf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Warn(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) Errorf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Error(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) Fatalf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.fields)
	var zapFields = l.getZapFields()
	l.logger.Fatal(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) DebugCf(msg string) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Debug(msg, zapFields...)
}

func (l *ZapLogger) InfoCf(msg string) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Info(msg, zapFields...)
}

func (l *ZapLogger) WarnCf(msg string) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Warn(msg, zapFields...)
}

func (l *ZapLogger) ErrorCf(msg string) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Error(msg, zapFields...)
}

func (l *ZapLogger) FatalCf(msg string) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Fatal(msg, zapFields...)
}

func (l *ZapLogger) DebugfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Debug(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) InfofCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Info(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) WarnfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Warn(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) ErrorfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Error(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) FatalfCf(msg string, args ...interface{}) {
	l.addContextCommonFields(l.customFields)
	var zapFields = l.getCustomZapFields()
	l.logger.Fatal(fmt.Sprintf(msg, args...), zapFields...)
}

func (l *ZapLogger) addContextCommonFields(fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.ctx != nil && l.ctx.Value(CommonFieldsKey) != nil {
		for k, v := range l.ctx.Value(CommonFieldsKey).(map[string]interface{}) {
			if _, ok := fields[k]; !ok {
				fields[k] = v
			}
		}
	}
	fields["func"], fields["file"] = fileAndFuncInfo(FileAndInfoSkipLevel)
}

func (l *ZapLogger) addCustomFields(fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.customFields = fields
}

func (l *ZapLogger) getZapFields() []zap.Field {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var zapFields []zap.Field
	for k, v := range l.fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func (l *ZapLogger) getCustomZapFields() []zap.Field {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var zapFields []zap.Field
	for k, v := range l.customFields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
