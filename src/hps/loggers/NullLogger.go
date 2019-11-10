package loggers

type NullLogger struct {
}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

func (nl *NullLogger) Print(v ...interface{}) {
}

func (nl *NullLogger) Printf(format string, v ...interface{}) {
}
