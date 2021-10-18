package logging

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/enorith/logging/writers"
	"go.uber.org/zap"
)

type Config struct {
	//BaseDir base directory of log files
	BaseDir string
}

func WithDefaults(conf Config) {
	zap.RegisterSink("rotate", func(u *url.URL) (zap.Sink, error) {
		format := u.Query().Get("time_format")
		return writers.NewRotateFile(filepath.Join(conf.BaseDir, u.Path), format), nil
	})

	zap.RegisterSink("single", func(u *url.URL) (zap.Sink, error) {
		return os.OpenFile(filepath.Join(conf.BaseDir, u.Path), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	})
}
