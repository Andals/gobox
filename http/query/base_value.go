package query

import (
	"andals/gobox/exception"
)

type Value interface {
	Set(str string) error
	Check() bool
	Error() *exception.Exception
}

type baseValue struct {
	errno int
	msg   string
}

func newBaseValue(errno int, msg string) *baseValue {
	return &baseValue{
		errno: errno,
		msg:   msg,
	}
}

func (this *baseValue) Error() *exception.Exception {
	return exception.New(this.errno, this.msg)
}
