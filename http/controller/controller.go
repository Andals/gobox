package controller

import (
	//     "fmt"
	"net/http"
	"regexp"
)

const (
	DEF_REMOTE_REAL_IP_HEADER_KEY   = "REMOTE-REAL-IP"
	DEF_REMOTE_REAL_PORT_HEADER_KEY = "REMOTE-REAL-PORT"
)

type ActionFunc func(context *Context, args []string)
type JumpFunc func(context *Context, args ...interface{})

type jumpItem struct {
	jf JumpFunc

	args []interface{}
}

type Controller struct {
	exactMatches map[string]ActionFunc
	regexMatches struct {
		regexSlice  []*regexp.Regexp
		actionSlice []ActionFunc
	}

	beforeActionMatches struct {
		regexSlice  []*regexp.Regexp
		actionSlice []ActionFunc
	}
	afterActionMatches struct {
		regexSlice  []*regexp.Regexp
		actionSlice []ActionFunc
	}

	error404Func JumpFunc

	//eg, access by nginx's proxy_pass
	remoteRealIpHeaderKey   string
	remoteRealPortHeaderKey string
}

func NewController() *Controller {
	this := new(Controller)

	this.exactMatches = make(map[string]ActionFunc)

	this.remoteRealIpHeaderKey = DEF_REMOTE_REAL_IP_HEADER_KEY
	this.remoteRealPortHeaderKey = DEF_REMOTE_REAL_PORT_HEADER_KEY

	return this
}

func (this *Controller) AddExactMatchAction(pattern string, af ActionFunc) {
	this.exactMatches[pattern] = af
}

func (this *Controller) AddRegexMatchAction(pattern string, af ActionFunc) {
	regex := regexp.MustCompile(pattern)

	this.regexMatches.regexSlice = append(this.regexMatches.regexSlice, regex)
	this.regexMatches.actionSlice = append(this.regexMatches.actionSlice, af)
}

func (this *Controller) AddBeforeAction(pattern string, af ActionFunc) {
	regex := regexp.MustCompile(pattern)

	this.beforeActionMatches.regexSlice = append(this.beforeActionMatches.regexSlice, regex)
	this.beforeActionMatches.actionSlice = append(this.beforeActionMatches.actionSlice, af)
}

func (this *Controller) AddAfterAction(pattern string, af ActionFunc) {
	regex := regexp.MustCompile(pattern)

	this.afterActionMatches.regexSlice = append(this.afterActionMatches.regexSlice, regex)
	this.afterActionMatches.actionSlice = append(this.afterActionMatches.actionSlice, af)
}

func (this *Controller) SetError404Func(jf JumpFunc) {
	this.error404Func = jf
}

func (this *Controller) SetRemoteRealIpHeaderKey(key string) {
	this.remoteRealIpHeaderKey = key
}

func (this *Controller) SetRemoteRealPortHeaderKey(key string) {
	this.remoteRealPortHeaderKey = key
}

func (this *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewContext(r, w, ParseRemoteAddr(r, this.remoteRealIpHeaderKey, this.remoteRealPortHeaderKey))

	this.dispatch(context)

	context.RespWriter.Write(context.RespBody)
}

func (this *Controller) dispatch(context *Context) {
	defer func() {
		if e := recover(); e != nil {
			ji, ok := e.(*jumpItem)
			if ok {
				ji.jf(context, ji.args...)
				return
			}
			panic(e)
		}
	}()

	af, args := this.findActionFunc(context.Req)
	if af == nil {
		if this.error404Func == nil {
			this.error404Func = error404
		}

		LongJump(this.error404Func)
	}

	baf, bargs := this.findBeforeActionFunc(context.Req)
	aaf, aargs := this.findAfterActionFunc(context.Req)

	baf(context, bargs)
	af(context, args)
	aaf(context, aargs)
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

func (this *Controller) findBeforeActionFunc(r *http.Request) (ActionFunc, []string) {
	path := r.URL.Path

	for i, regex := range this.beforeActionMatches.regexSlice {
		matches := regex.FindStringSubmatch(path)
		if matches != nil {
			return this.beforeActionMatches.actionSlice[i], matches[1:]
		}
	}

	return NoopBeforeAction, nil
}

func (this *Controller) findAfterActionFunc(r *http.Request) (ActionFunc, []string) {
	path := r.URL.Path

	for i, regex := range this.afterActionMatches.regexSlice {
		matches := regex.FindStringSubmatch(path)
		if matches != nil {
			return this.afterActionMatches.actionSlice[i], matches[1:]
		}
	}

	return NoopAfterAction, nil
}

func NoopBeforeAction(context *Context, args []string) {
}

func NoopAfterAction(context *Context, args []string) {
}

func LongJump(jf JumpFunc, args ...interface{}) {
	ji := &jumpItem{
		jf:   jf,
		args: args,
	}

	panic(ji)
}

func Redirect302(url string) {
	LongJump(redirect302, url)
}

func error404(context *Context, args ...interface{}) {
	http.NotFound(context.RespWriter, context.Req)
}

func redirect302(context *Context, args ...interface{}) {
	http.Redirect(context.RespWriter, context.Req, args[0].(string), 302)
}
