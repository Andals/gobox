package query

import (
	"andals/gobox/exception"
)

type Value interface {
	Required() bool
	Set(str string) error
	Check() bool
	Error() *exception.Exception
}

type baseValue struct {
	required bool
	errno    int
	msg      string
}

func newBaseValue(required bool, errno int, msg string) *baseValue {
	return &baseValue{
		required: required,
		errno:    errno,
		msg:      msg,
	}
}

func (this *baseValue) Required() bool {
	return this.required
}

func (this *baseValue) Error() *exception.Exception {
	return exception.New(this.errno, this.msg)
}
