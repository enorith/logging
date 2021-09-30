package logging_test

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Foo struct {
	A int
}

func TestZap(t *testing.T) {
	conf := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			TimeKey:       "time",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			StacktraceKey: "trace",
		},
		OutputPaths:      []string{"./tmp/log.log", "stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := conf.Build()
	defer logger.Sync()

	logger.Error("failed to fetch URL",
		zap.String("url", "www.google.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
