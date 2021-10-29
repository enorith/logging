package logging

import (
	"fmt"

	"go.uber.org/zap"
)

func Debug(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

func Debugf(msg string, args ...interface{}) {
	Debug(fmt.Sprintf(msg, args...))
}

func Info(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

func Infof(msg string, args ...interface{}) {
	Info(fmt.Sprintf(msg, args...))
}

func Warn(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

func Warnf(msg string, args ...interface{}) {
	Warn(fmt.Sprintf(msg, args...))
}

func Fatal(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

func Fatalf(msg string, args ...interface{}) {
	Fatal(fmt.Sprintf(msg, args...))
}

func Panic(msg string, fields ...zap.Field) {
	WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
}

func Panicf(msg string, args ...interface{}) {
	Panic(fmt.Sprintf(msg, args...))
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
