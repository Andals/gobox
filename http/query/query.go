package query

import (
	"net/url"
	"strings"

	"andals/gobox/exception"
)

type QuerySet struct {
	formal map[string]Value
}

func NewQuerySet() *QuerySet {
	this := &QuerySet{
		formal: make(map[string]Value),
	}

	return this
}

func (this *QuerySet) Var(name string, v Value) {
	this.formal[name] = v
}

func (this *QuerySet) IntVar(p *int, name string, errno int, msg string, cf CheckInt) {
	this.Var(name, NewIntValue(p, errno, msg, cf))
}

func (this *QuerySet) StringVar(p *string, name string, errno int, msg string, cf CheckString) {
	this.Var(name, NewStringValue(p, errno, msg, cf))
}

func (this *QuerySet) Int64Var(p *int64, name string, errno int, msg string, cf CheckInt64) {
	this.Var(name, NewInt64Value(p, errno, msg, cf))
}

func (this *QuerySet) Parse(actual *url.Values) *exception.Exception {
	for name, v := range this.formal {
		str := strings.TrimSpace(actual.Get(name))
		err := v.Set(str)
		if err != nil {
			return v.Error()
		}
		if v.Check() == false {
			return v.Error()
		}
	}

	return nil
}
