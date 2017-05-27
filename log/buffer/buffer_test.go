package buffer

import (
	//     "fmt"
	"testing"
	"time"

	"andals/gobox/log/writer"
)

func TestBufferFileWriter(t *testing.T) {
	Init(1024, time.Second*7)

	path := "/tmp/test_buffer.log"
	bufsize := 4096

	fw, _ := writer.NewFileWriter(path)
	bw := NewBuffer(fw, bufsize)

	bw.Write([]byte("test file writer with buffer and time interval\n"))
	time.Sleep(time.Second * 5)
	bw.Free()
}
