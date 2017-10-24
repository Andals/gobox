package query

import "strconv"

type CheckInt func(v int) bool

func CheckIntIsPositive(v int) bool {
	if v > 0 {
		return true
	}
	return false
}

type intValue struct {
	*baseValue

	p  *int
	cf CheckInt
}

func NewIntValue(p *int, required bool, errno int, msg string, cf CheckInt) *intValue {
	this := &intValue{
		baseValue: newBaseValue(required, errno, msg),

		p:  p,
		cf: cf,
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
	if this.cf == nil {
		return true
	}

	return this.cf(*(this.p))
}
