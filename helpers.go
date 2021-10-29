package logging

import "go.uber.org/zap"

func Debug(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return Channel().With(fields...)
}

func WithOptions(options ...zap.Option) *zap.Logger {
	return Channel().WithOptions(options...)
}

func Channel(channel ...string) *zap.Logger {
	if l, e := DefaultManager.Channel(channel...); e == nil {
		return l
	}

	return NilLogger
}
