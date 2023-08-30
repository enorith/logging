package logging_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/enorith/logging"
	"github.com/enorith/logging/writers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Foo struct {
	A int
}

func TestZap(t *testing.T) {
	wd, _ := os.Getwd()
	logging.WithDefaults(logging.Config{
		BaseDir: wd,
	})

	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"rotate:///tmp/test.log?time_format=2006-01-02T15", "single:///tmp/single.log"}
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	conf.ErrorOutputPaths = []string{"stderr"}
	logger, _ := conf.Build()
	defer logger.Sync()

	logger.Error("failed to fetch URL",
		zap.String("url", "www.google.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

func TestRotate(t *testing.T) {
	w := writers.NewRotateFile("tmp/test.log", "2006-01-02T15 04 05")
	for i := 0; i < 10; i++ {
		_, e := w.Write([]byte(fmt.Sprintf("%s hello\n", time.Now().Format("2006-01-02T15:04:05.999Z07:00"))))
		if e != nil {
			t.Error(e)
		}
	}
	w.SetLimit(5)
	w.Cleanup()
}

func Benchmark_Rotate(b *testing.B) {
	w := writers.NewRotateFile("tmp/test.log")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, e := w.Write([]byte(fmt.Sprintf("%s hello\n", time.Now().Format(time.RFC3339))))
		if e != nil {
			b.Error(e)
		}
	}
}

func TestManager(t *testing.T) {
	wd, _ := os.Getwd()
	logging.WithDefaults(logging.Config{
		BaseDir: wd,
	})
	logging.DefaultManager.Resolve("default", func(conf zap.Config) (*zap.Logger, error) {
		conf.OutputPaths = []string{"rotate:///tmp/rotates/enorith.log?time_format=2006-01-02T15 04 05&limit=3", "stdout"}
		conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.999")
		conf.EncoderConfig.StacktraceKey = "trace"
		conf.Encoding = "console"
		return conf.Build()
	})

	logger, e := logging.DefaultManager.Channel()
	if e != nil {
		t.Error(e)
		t.Fail()
	}

	logger.Error("test", zap.Int("answer", 42))

	logging.Info("info snappy")
	logging.Infof("info snappy arg %s", "data")
	logging.With(zap.String("foo", "bar")).Info("info snappy with")

	logging.DefaultManager.Sync()
	logging.DefaultManager.Cleanup()
}

func TestNil(t *testing.T) {
	logging.FallbackLogger.Info("check")
	logging.FallbackLogger.Error("check err")
}
