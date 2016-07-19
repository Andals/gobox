package writer

import (
	//     "fmt"
	"testing"
	"time"
)

func TestBufferFileWriter(t *testing.T) {
	path := "/tmp/test.log"
	bufsize := 4096

	w, _ := NewFileWriter(path)
	wt := NewBufferWriter(w, bufsize, time.Second*3)

	wt.Write([]byte("test file writer with buffer and time interval\n"))
	time.Sleep(time.Second * 5)

	wd, _ := NewFileWriterWithSplit(path, SPLIT_BY_DAY)
	writer := NewBufferWriter(wd, bufsize, 0)

	writer.Write([]byte("test file writer with buffer and split by day\n"))
	writer.Free()

	wh, _ := NewFileWriterWithSplit(path, SPLIT_BY_HOUR)
	writer = NewBufferWriter(wh, bufsize, 0)

	writer.Write([]byte("test file writer with buffer and split by hour\n"))
	writer.Free()
}
