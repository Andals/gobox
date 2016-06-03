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
* @name IntValue
* @{ */

type CheckInt func(v int) bool

type IntValue struct {
	*baseValue

	p  *int
	cf CheckInt
}

func NewIntValue(p *int, errno int, msg string, cf CheckInt) *IntValue {
	this := &IntValue{
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

func (this *IntValue) Set(str string) error {
	v, e := strconv.Atoi(str)
	if e != nil {
		return e
	}

	*(this.p) = v

	return nil
}

func (this *IntValue) Check() bool {
	return this.cf(*(this.p))
}

/**  @} */

/**
* @name StringValue
* @{ */

type CheckString func(v string) bool

type StringValue struct {
	*baseValue

	p  *string
	cf CheckString
}

func NewStringValue(p *string, errno int, msg string, cf CheckString) *StringValue {
	this := &StringValue{
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

func (this *StringValue) Set(str string) error {
	*(this.p) = str

	return nil
}

func (this *StringValue) Check() bool {
	return this.cf(*(this.p))
}

/**  @} */
