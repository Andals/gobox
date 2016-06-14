package query

import (
	"strconv"

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

/**
* @name intValue
* @{ */

type CheckInt func(v int) bool

type intValue struct {
	*baseValue

	p  *int
	cf CheckInt
}

func NewIntValue(p *int, errno int, msg string, cf CheckInt) *intValue {
	this := &intValue{
		baseValue: newBaseValue(errno, msg),

		p: p,
	}

	if cf == nil {
		this.cf = func(v int) bool {
			if v == 0 {
				return false
			}
			return true
		}
	} else {
		this.cf = cf
	}

	return this
}

func (this *intValue) Set(str string) error {
	var v int = 0
	var e error = nil

	if str != "" {
		v, e = strconv.Atoi(str)
	}

	if e != nil {
		return e
	}

	*(this.p) = v

	return nil
}

func (this *intValue) Check() bool {
	return this.cf(*(this.p))
}

/**  @} */

/**
* @name stringValue
* @{ */

type CheckString func(v string) bool

type stringValue struct {
	*baseValue

	p  *string
	cf CheckString
}

func NewStringValue(p *string, errno int, msg string, cf CheckString) *stringValue {
	this := &stringValue{
		baseValue: newBaseValue(errno, msg),

		p: p,
	}

	if cf == nil {
		this.cf = func(v string) bool {
			if v == "" {
				return false
			}
			return true
		}
	} else {
		this.cf = cf
	}

	return this
}

func (this *stringValue) Set(str string) error {
	*(this.p) = str

	return nil
}

func (this *stringValue) Check() bool {
	return this.cf(*(this.p))
}

/**  @} */
