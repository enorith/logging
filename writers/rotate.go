package writers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var DefaultRotateTimeFormat = "2006-01-02"

type RotateFileWriter struct {
	lock                           sync.Mutex
	filename, timeFormat, fileTime string
	name, ext, dir                 string

	out *os.File

	limit int
}

func NewRotateFile(filename string, timeFormat ...string) *RotateFileWriter {
	var tf string
	if len(timeFormat) > 0 && timeFormat[0] != "" {
		tf = timeFormat[0]
	} else {
		tf = DefaultRotateTimeFormat
	}

	return &RotateFileWriter{filename: filename, timeFormat: tf}
}

// Write satisfies the io.Writer interface.
func (w *RotateFileWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	out, e := w.RotateWriterNoLock()
	if e != nil {
		return 0, e
	}

	return out.Write(output)
}

func (w *RotateFileWriter) Close() error {
	if w.out != nil {
		return w.out.Close()
	}

	return nil
}

func (w *RotateFileWriter) Sync() error {
	if w.out != nil {
		return w.out.Sync()
	}

	return nil
}

func (w *RotateFileWriter) RotateWriterNoLock() (out *os.File, err error) {
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
	err = w.prepareFile()
	if err != nil {
		return
	}

	realname := fmt.Sprintf("%s.%s%s", w.name, w.fileTime, w.ext)

	out, err = os.OpenFile(realname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	w.out = out

	return
}

func (w *RotateFileWriter) prepareFile() error {
	if w.ext != "" && w.name != "" && w.dir != "" {
		return nil
	}

	var realpath string
	realpath, err := filepath.Abs(w.filename)
	if err != nil {
		return err
	}
	w.dir = filepath.Dir(realpath)
	w.ext = filepath.Ext(realpath)
	if w.ext == "" {
		w.ext = "log"
	}
	dot := strings.LastIndexByte(realpath, '.')
	if dot == -1 {
		dot = len(realpath)
	}
	w.name = realpath[0:dot]

	return nil
}

func (w *RotateFileWriter) SetLimit(limit int) *RotateFileWriter {
	w.limit = limit
	return w
}

func (w *RotateFileWriter) Cleanup() error {
	if w.limit < 1 {
		return nil
	}
	err := w.prepareFile()
	if err != nil {
		return err
	}
	var files []string

	filepath.WalkDir(w.dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if strings.Index(path, w.name) == 0 {
				files = append(files, path)
			}
		}

		return nil
	})

	for len(files) > w.limit {
		var file string
		file, files = files[0], files[1:]
		os.Remove(file)
	}

	return nil
}
