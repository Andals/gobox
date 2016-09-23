package controller

import (
	//     "fmt"
	"net/http"
)

const (
	DEF_REMOTE_REAL_IP_HEADER_KEY   = "REMOTE-REAL-IP"
	DEF_REMOTE_REAL_PORT_HEADER_KEY = "REMOTE-REAL-PORT"
)

type jumpItem struct {
	jf JumpFunc

	args []interface{}
}

type Controller struct {
	actionMatches       *funcMatches
	beforeActionMatches *funcMatches
	afterActionMatches  *funcMatches
	destructFuncMatches *funcMatches

	missActionFunc JumpFunc

	//eg, access by nginx's proxy_pass
	remoteRealIpHeaderKey   string
	remoteRealPortHeaderKey string
}

func NewController() *Controller {
	return &Controller{
		actionMatches:       newFuncMatches(nil),
		beforeActionMatches: newFuncMatches(ActionFunc(NoopAction)),
		afterActionMatches:  newFuncMatches(ActionFunc(NoopAction)),
		destructFuncMatches: newFuncMatches(DestructFunc(NoopDestruct)),

		remoteRealIpHeaderKey:   DEF_REMOTE_REAL_IP_HEADER_KEY,
		remoteRealPortHeaderKey: DEF_REMOTE_REAL_PORT_HEADER_KEY,
	}
}

func (this *Controller) AddExactMatchAction(key string, f ActionFunc) *Controller {
	this.actionMatches.addExactFunc(key, f)

	return this
}

func (this *Controller) AddRegexMatchAction(pattern string, f ActionFunc) *Controller {
	this.actionMatches.addRegexFunc(pattern, f)

	return this
}

func (this *Controller) AddBeforeAction(pattern string, f ActionFunc) *Controller {
	this.beforeActionMatches.addRegexFunc(pattern, f)

	return this
}

func (this *Controller) AddAfterAction(pattern string, f ActionFunc) *Controller {
	this.afterActionMatches.addRegexFunc(pattern, f)

	return this
}

func (this *Controller) AddDestructFunc(pattern string, f DestructFunc) *Controller {
	this.destructFuncMatches.addRegexFunc(pattern, f)

	return this
}

func (this *Controller) SetMissActionFunc(f JumpFunc) *Controller {
	this.missActionFunc = f

	return this
}

func (this *Controller) SetRemoteRealIpHeaderKey(key string) *Controller {
	this.remoteRealIpHeaderKey = key

	return this
}

func (this *Controller) SetRemoteRealPortHeaderKey(key string) *Controller {
	this.remoteRealPortHeaderKey = key

	return this
}

func (this *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewContext(r, w, ParseRemoteAddr(r, this.remoteRealIpHeaderKey, this.remoteRealPortHeaderKey))

	this.dispatch(context)

	df, dargs := this.findDestructFunc(context.Req)
	df(context.TransData, dargs)

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
		if this.missActionFunc == nil {
			this.missActionFunc = error404
		}

		LongJump(this.missActionFunc)
	}

	baf, bargs := this.findBeforeActionFunc(context.Req)
	aaf, aargs := this.findAfterActionFunc(context.Req)

	baf(context, bargs)
	af(context, args)
	aaf(context, aargs)
}

func (this *Controller) findActionFunc(r *http.Request) (ActionFunc, []string) {
	f, args := this.actionMatches.findFunc(r)

	return f.(ActionFunc), args
}

func (this *Controller) findBeforeActionFunc(r *http.Request) (ActionFunc, []string) {
	f, args := this.beforeActionMatches.findFunc(r)

	return f.(ActionFunc), args
}

func (this *Controller) findAfterActionFunc(r *http.Request) (ActionFunc, []string) {
	f, args := this.afterActionMatches.findFunc(r)

	return f.(ActionFunc), args
}

func (this *Controller) findDestructFunc(r *http.Request) (DestructFunc, []string) {
	f, args := this.destructFuncMatches.findFunc(r)

	return f.(DestructFunc), args
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

func redirect302(context *Context, args ...interface{}) {
	http.Redirect(context.RespWriter, context.Req, args[0].(string), 302)
}

func error404(context *Context, args ...interface{}) {
	http.NotFound(context.RespWriter, context.Req)
}
