package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getDefaultLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment(zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))
	return logger
}

func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	getDefaultLogger().Info(msg, fields...)
}

func Infof(ctx context.Context, format string, args ...any) {
	getDefaultLogger().Info(fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	getDefaultLogger().Error(msg, fields...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	getDefaultLogger().Error(fmt.Sprintf(format, args...))
}

func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	getDefaultLogger().Warn(msg, fields...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	getDefaultLogger().Warn(fmt.Sprintf(format, args...))
}
