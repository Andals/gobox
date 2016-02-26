/**
* @file buffer.go
* @brief writer with buffer
* @author ligang
* @date 2016-02-04
 */

package writer

import (
	"bufio"
)

type Buffer struct {
	IWriter

	buf *bufio.Writer
}

func NewBufferWriter(writer IWriter, bufsize int) *Buffer {
	this := &Buffer{
		IWriter: writer,
		buf:     bufio.NewWriterSize(writer, bufsize),
	}

	return this
}

func (this *Buffer) Flush() error {
	return this.buf.Flush()
}

func (this *Buffer) Free() {
	this.Flush()
	this.IWriter.Free()
	this.buf = nil
}
