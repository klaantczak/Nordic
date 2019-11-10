package loggers

import (
	"fmt"
)

type ConsoleLogger struct {
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (cl *ConsoleLogger) Print(v ...interface{}) {
	fmt.Print(v...)
	fmt.Println("")
}

func (cl *ConsoleLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	fmt.Println("")
}
