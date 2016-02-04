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
	w IWriter

	*bufio.Writer
}

func NewBufferWriter(writer IWriter, bufsize int) *Buffer {
	this := &Buffer{
		w:      writer,
		Writer: bufio.NewWriterSize(writer, bufsize),
	}

	return this
}

func (this *Buffer) Free() {
	this.Flush()
	this.w.Free()
	this.Writer = nil
}
