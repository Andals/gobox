package query

import (
	"andals/gobox/exception"

	"net/url"
	"strings"
)

type QuerySet struct {
	formal map[string]Value
	exists map[string]bool
}

func NewQuerySet() *QuerySet {
	this := &QuerySet{
		formal: make(map[string]Value),
		exists: make(map[string]bool),
	}

	return this
}

func (this *QuerySet) ExistsInfo() map[string]bool {
	return this.exists
}

func (this *QuerySet) Exist(name string) bool {
	exist, ok := this.exists[name]
	if !ok || exist == false {
		return false
	}

	return true
}

func (this *QuerySet) Var(name string, v Value) *QuerySet {
	this.formal[name] = v

	return this
}

func (this *QuerySet) IntVar(p *int, name string, required bool, errno int, msg string, cf CheckInt) *QuerySet {
	this.Var(name, NewIntValue(p, required, errno, msg, cf))

	return this
}

func (this *QuerySet) StringVar(p *string, name string, required bool, errno int, msg string, cf CheckString) *QuerySet {
	this.Var(name, NewStringValue(p, required, errno, msg, cf))

	return this
}

func (this *QuerySet) Int64Var(p *int64, name string, required bool, errno int, msg string, cf CheckInt64) *QuerySet {
	this.Var(name, NewInt64Value(p, required, errno, msg, cf))

	return this
}

func (this *QuerySet) Parse(actual url.Values) *exception.Exception {
	for name, v := range this.formal {
		if len(actual[name]) == 0 {
			this.exists[name] = false
			if v.Required() {
				return v.Error()
			}
			continue
		}

		this.exists[name] = true
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
