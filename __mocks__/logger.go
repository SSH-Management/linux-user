package mocks

import (
	"fmt"
)

type BufferLogger struct {
	data []string
}

func NewLogger() *BufferLogger {
	return &BufferLogger{
		data: make([]string, 0, 100),
	}
}

func (b *BufferLogger) Print(msg string, data ...interface{}) {
	b.data = append(b.data, fmt.Sprintf(msg, data...))
}

func (b BufferLogger) Len() int {
	return len(b.data)
}

func (b BufferLogger) Data() []string {
	return b.data
}
