/**
* @file exception.go
* @brief exception with errno and msg
* @author ligang
* @date 2016-02-01
 */

package exception

import (
	"strconv"
)

type Exception struct {
	errno int
	msg   string
}

func New(errno int, msg string) *Exception {
	return &Exception{
		errno: errno,
		msg:   msg,
	}
}

func (this *Exception) Error() string {
	result := "errno: " + strconv.Itoa(this.errno) + ", "
	result += "msg: " + this.msg

	return result
}

func (this *Exception) Errno() int {
	return this.errno
}

func (this *Exception) Msg() string {
	return this.msg
}
