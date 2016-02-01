/**
* @file error.go
* @brief error with errno and msg
* @author ligang
* @date 2016-02-01
 */

package error

import (
	"strconv"
)

type Error struct {
	errno int
	msg   string
}

func NewError(errno int, msg string) *Error {
	return &Error{
		errno: errno,
		msg:   msg,
	}
}

func (this *Error) Error() string {
	result := "errno: " + strconv.Itoa(this.errno) + ", "
	result += "msg: " + this.msg

	return result
}

func (this *Error) GetErrno() int {
	return this.errno
}

func (this *Error) GetMsg() string {
	return this.msg
}
