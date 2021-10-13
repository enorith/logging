package logging_test

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/enorith/logging/writers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Foo struct {
	A int
}

func TestZap(t *testing.T) {
	zap.RegisterSink("rotate", func(u *url.URL) (zap.Sink, error) {
		wd, _ := os.Getwd()

		return writers.NewRotate(filepath.Join(wd, u.Path)), nil
	})

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
		OutputPaths:      []string{"rotate:///tmp/log.log", "stdout"},
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

func TestRotate(t *testing.T) {
	w := writers.NewRotate("tmp/test.log")
	for i := 0; i < 10; i++ {
		_, e := w.Write([]byte(fmt.Sprintf("%s hello\n", time.Now().Format("2006-01-02T15:04:05.999Z07:00"))))
		if e != nil {
			t.Error(e)
		}
	}
}

func Benchmark_Rotate(b *testing.B) {
	w := writers.NewRotate("tmp/test.log")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, e := w.Write([]byte(fmt.Sprintf("%s hello\n", time.Now().Format(time.RFC3339))))
		if e != nil {
			b.Error(e)
		}
	}
}
