package controller

import (
	"net/http"
	"net/url"
	"regexp"
)

type Context struct {
	RespWriter http.ResponseWriter
	Req        *http.Request

	queryValues url.Values
	TransData   map[string]interface{}
	RespBody    []byte
}

type ActionFunc func(context *Context, args []string)

type Controller struct {
	exactMatches map[string]ActionFunc
	regexMatches struct {
		regexSlice  []*regexp.Regexp
		actionSlice []ActionFunc
	}

	beforeAction ActionFunc
	afterAction  ActionFunc
}

func NewController() *Controller {
	this := new(Controller)

	this.exactMatches = make(map[string]ActionFunc)
	this.beforeAction = DefaultBeforeAction
	this.afterAction = DefaultAfterAction

	return this
}

func (this *Controller) ExactMatchAction(pattern string, af ActionFunc) {
	this.exactMatches[pattern] = af
}

func (this *Controller) RegexMatchAction(pattern string, af ActionFunc) {
	regex := regexp.MustCompile(pattern)

	this.regexMatches.regexSlice = append(this.regexMatches.regexSlice, regex)
	this.regexMatches.actionSlice = append(this.regexMatches.actionSlice, af)
}

func (this *Controller) BeforeAction(af ActionFunc) {
	this.beforeAction = af
}

func (this *Controller) AfterAction(af ActionFunc) {
	this.afterAction = af
}

func (this *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{
		RespWriter: w,
		Req:        r,

		queryValues: r.URL.Query(),
		TransData:   make(map[string]interface{}),
	}

	af, args := this.findActionFunc(r)
	if af == nil {
		http.NotFound(w, r)
		return
	}

	this.beforeAction(context, args)
	af(context, args)
	this.afterAction(context, args)

	context.RespWriter.Write(context.RespBody)
}

func (this *Controller) findActionFunc(r *http.Request) (ActionFunc, []string) {
	path := r.URL.Path

	af, ok := this.exactMatches[path]
	if ok {
		return af, nil
	}

	for i, regex := range this.regexMatches.regexSlice {
		matches := regex.FindStringSubmatch(path)
		if matches != nil {
			return this.regexMatches.actionSlice[i], matches[1:]
		}
	}

	return nil, nil
}

func DefaultBeforeAction(context *Context, args []string) {
}

func DefaultAfterAction(context *Context, args []string) {
}
