package loggers

import (
	"fmt"
	"hps"
)

type ContextLogger struct {
	logger  hps.ILogger
	context string
}

func NewContextLogger(logger hps.ILogger, context string) *ContextLogger {
	cl := &ContextLogger{logger, context}
	return cl
}

func (cl *ContextLogger) Print(v ...interface{}) {
	cl.logger.Printf("[%s] %s", cl.context, fmt.Sprint(v...))
}

func (cl *ContextLogger) Printf(format string, v ...interface{}) {
	cl.logger.Printf("[%s] %s", cl.context, fmt.Sprintf(format, v...))
}
