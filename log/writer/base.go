/**
* @file base.go
* @brief writer interface
* @author ligang
* @date 2016-02-04
 */

package writer

import (
	"io"
)

type IWriter interface {
	io.Writer

	Flush()
	Free()
}
