package writers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var DefaultRotateTimeFormat = "2006-01-02"

type RotateWriter struct {
	lock                           sync.Mutex
	filename, timeFormat, fileTime string

	out *os.File
}

func NewRotate(filename string, timeFormat ...string) *RotateWriter {
	var tf string
	if len(timeFormat) > 0 {
		tf = timeFormat[0]
	} else {
		tf = DefaultRotateTimeFormat
	}

	return &RotateWriter{filename: filename, timeFormat: tf}
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	out, e := w.RotateWriterNoLock()
	if e != nil {
		return 0, e
	}

	return out.Write(output)
}

func (w *RotateWriter) Close() error {
	if w.out != nil {
		return w.out.Close()
	}

	return nil
}

func (w *RotateWriter) Sync() error {
	if w.out != nil {
		return w.out.Sync()
	}

	return nil
}

func (w *RotateWriter) RotateWriterNoLock() (out *os.File, err error) {
	currentTime := time.Now().Format(w.timeFormat)
	if w.out != nil && currentTime == w.fileTime {
		return w.out, nil
	}
	// rotate to another file
	if w.out != nil {
		err = w.out.Close()
		if err != nil {
			return
		}
	}
	w.fileTime = currentTime
	var realpath string
	realpath, err = filepath.Abs(w.filename)
	if err != nil {
		return
	}
	ext := filepath.Ext(realpath)
	if ext == "" {
		ext = "log"
	}
	dot := strings.LastIndexByte(realpath, '.')
	if dot == -1 {
		dot = len(realpath)
	}
	name := realpath[0:dot]
	realname := fmt.Sprintf("%s.%s%s", name, w.fileTime, ext)

	out, err = os.OpenFile(realname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	w.out = out

	return
}
