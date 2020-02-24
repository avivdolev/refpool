package example

import (
	"bytes"
)

// Buffer is an example for Element for refpool
type Buffer struct {
	bytes.Buffer
	count int64
}

// Counter implement Element interface:
func (b *Buffer) Counter() *int64 {
	return &b.count
}
