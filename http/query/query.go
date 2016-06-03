package query

import (
	"net/url"

	"andals/gobox/exception"
)

type QuerySet struct {
	actual *url.Values
	formal map[string]Value
}

func NewQuerySet(values *url.Values) *QuerySet {
	this := &QuerySet{
		actual: values,
		formal: make(map[string]Value),
	}

	return this
}

func (this *QuerySet) Var(name string, v Value) {
	this.formal[name] = v
}

func (this *QuerySet) Parse() *exception.Exception {
	for name, v := range this.formal {
		err := v.Set(this.actual.Get(name))
		if err != nil {
			return v.Error()
		}
		if v.Check() == false {
			return v.Error()
		}
	}

	return nil
}
