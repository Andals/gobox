/**
* @file noop.go
* @brief writer do nothing
* @author ligang
* @date 2016-02-03
 */

package writer

type Noop struct {
}

func (this *Noop) write(msg []byte) {
}
