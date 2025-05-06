package logger

import (
	"context"
)

func GetNewLogger(ctx context.Context, loggerType string, level ...string) Logger {
	var logger Logger

	var logLevel = ""
	if len(level) > 0 {
		logLevel = level[0]
	}

	switch loggerType {
	case Logrus:
		logger = NewLogrusLogger(ctx, logLevel)
	case Zap:
		logger = NewZapLogger(ctx, logLevel)
	default:
		logger = NewLogrusLogger(ctx, logLevel)
	}

	return logger
}
