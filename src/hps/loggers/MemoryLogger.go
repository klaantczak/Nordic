package loggers

import (
	"fmt"
	"strings"
)

type MemoryLogger struct {
	messages []string
}

func NewMemoryLogger() *MemoryLogger {
	return &MemoryLogger{[]string{}}
}

func (ml *MemoryLogger) Print(v ...interface{}) {
	ml.messages = append(ml.messages, fmt.Sprint(v...))
}

func (ml *MemoryLogger) Printf(format string, v ...interface{}) {
	ml.messages = append(ml.messages, fmt.Sprintf(format, v...))
}

func (ml *MemoryLogger) ToString() string {
	return strings.Join(ml.messages, "\n")
}
