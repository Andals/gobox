/**
* @file noop.go
* @brief writer do nothing
* @author ligang
* @date 2016-02-03
 */

package writer

type Noop struct {
}

func (this *Noop) Write(msg []byte) (int, error) {
	return 0, nil
}

func (this *Noop) Flush() error {
	return nil
}

func (this *Noop) Free() {
}
