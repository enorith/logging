package logging

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/enorith/logging/writers"
	"go.uber.org/zap"
)

var FallbackLogger = NewNilLogger()

type Config struct {
	//BaseDir base directory of log files
	BaseDir string
	//Fallback logger when logger not found
	Fallback *zap.Logger
}

func WithDefaults(conf Config) {
	zap.RegisterSink("rotate", func(u *url.URL) (zap.Sink, error) {
		format := u.Query().Get("time_format")
		makeDir(conf.BaseDir, u.Path)

		return writers.NewRotateFile(filepath.Join(conf.BaseDir, u.Path), format), nil
	})

	zap.RegisterSink("single", func(u *url.URL) (zap.Sink, error) {
		makeDir(conf.BaseDir, u.Path)
		return os.OpenFile(filepath.Join(conf.BaseDir, u.Path), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	})

	if conf.Fallback != nil {
		FallbackLogger = conf.Fallback
	}
}

func makeDir(baseDir, path string) {
	dir := filepath.Join(baseDir, filepath.Dir(path))
	os.MkdirAll(dir, 0775)
}

func NewNilLogger() *zap.Logger {
	conf := zap.NewProductionConfig()

	conf.OutputPaths = nil
	conf.ErrorOutputPaths = nil
	l, _ := conf.Build()
	return l
}

func NewStdLogger() *zap.Logger {
	conf := zap.NewProductionConfig()

	conf.OutputPaths = []string{"stdout"}
	conf.ErrorOutputPaths = []string{"stderr"}
	l, _ := conf.Build()
	return l
}
