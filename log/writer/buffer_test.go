package writer

import (
	//     "fmt"
	"testing"
)

func TestBufferFileWriter(t *testing.T) {
	path := "/tmp/test.log"
	bufsize := 4096

	w, _ := NewFileWriter(path)
	writer := NewBufferWriter(w, bufsize)

	writer.Write([]byte("test file writer with buffer\n"))
	writer.Free()

	wd, _ := NewFileWriterWithSplit(path, SPLIT_BY_DAY)
	writer = NewBufferWriter(wd, bufsize)

	writer.Write([]byte("test file writer with buffer and split by day\n"))
	writer.Free()

	wh, _ := NewFileWriterWithSplit(path, SPLIT_BY_HOUR)
	writer = NewBufferWriter(wh, bufsize)

	writer.Write([]byte("test file writer with buffer and split by hour\n"))
	writer.Free()
}
