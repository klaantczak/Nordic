package loggers

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type FileLogger struct {
	path    string
	file    *os.File
	writer  io.Writer
	newline []byte
}

func NewFileLogger(path string) (*FileLogger, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(f)
	return &FileLogger{path, f, w, []byte("\n")}, nil
}

func (fl *FileLogger) Print(v ...interface{}) {
	fl.writer.Write([]byte(fmt.Sprint(v...)))
	fl.writer.Write(fl.newline)
}

func (fl *FileLogger) Printf(format string, v ...interface{}) {
	fl.writer.Write([]byte(fmt.Sprintf(format, v...)))
	fl.writer.Write(fl.newline)
}
