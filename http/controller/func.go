package controller

import (
	//     "fmt"
	"net/http"
	"regexp"
)

/**
* @name funcMatches
* @{ */

type funcMatches struct {
	exactMatches map[string]interface{}

	regexSlice []*regexp.Regexp
	funcSlice  []interface{}

	defFunc interface{}
}

func newFuncMatches(defFunc interface{}) *funcMatches {
	return &funcMatches{
		exactMatches: make(map[string]interface{}),

		defFunc: defFunc,
	}
}

func (this *funcMatches) addExactFunc(key string, f interface{}) {
	this.exactMatches[key] = f
}

func (this *funcMatches) addRegexFunc(pattern string, f interface{}) {
	regex := regexp.MustCompile(pattern)

	this.regexSlice = append(this.regexSlice, regex)
	this.funcSlice = append(this.funcSlice, f)
}

func (this *funcMatches) findFunc(r *http.Request) (interface{}, []string) {
	path := r.URL.Path

	f, ok := this.exactMatches[path]
	if ok {
		return f, nil
	}

	for i, regex := range this.regexSlice {
		matches := regex.FindStringSubmatch(path)
		if matches != nil {
			return this.funcSlice[i], matches[1:]
		}
	}

	return this.defFunc, nil
}

type ActionFunc func(context *Context, args []string)
type JumpFunc func(context *Context, args ...interface{})
type DestructFunc func(transData map[string]interface{}, args []string)

func NoopAction(context *Context, args []string) {
}

func NoopDestruct(transData map[string]interface{}, args []string) {
}
