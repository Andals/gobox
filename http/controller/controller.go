package controller

import (
	//     "fmt"
	"net/http"
	"net/url"
	"regexp"
)

const (
	DEF_REMOTE_REAL_IP_HEADER_KEY   = "REMOTE-REAL-IP"
	DEF_REMOTE_REAL_PORT_HEADER_KEY = "REMOTE-REAL-PORT"
)

type Context struct {
	RespWriter http.ResponseWriter
	Req        *http.Request

	QueryValues *url.Values
	TransData   map[string]interface{}
	RespBody    []byte

	RemoteRealAddr *RemoteAddr
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

	remoteRealIpHeaderKey   string
	remoteRealPortHeaderKey string
}

func NewController() *Controller {
	this := new(Controller)

	this.exactMatches = make(map[string]ActionFunc)
	this.beforeAction = DefaultBeforeAction
	this.afterAction = DefaultAfterAction

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

func (this *Controller) SetBeforeAction(af ActionFunc) {
	this.beforeAction = af
}

func (this *Controller) SetAfterAction(af ActionFunc) {
	this.afterAction = af
}

func (this *Controller) SetRemoteRealIpHeaderKey(key string) {
	this.remoteRealIpHeaderKey = key
}

func (this *Controller) SetRemoteRealPortHeaderKey(key string) {
	this.remoteRealPortHeaderKey = key
}

func (this *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	af, args := this.findActionFunc(r)
	if af == nil {
		http.NotFound(w, r)
		return
	}

	context := &Context{
		RespWriter: w,
		Req:        r,

		TransData:      make(map[string]interface{}),
		RemoteRealAddr: ParseRemoteAddr(r, this.remoteRealIpHeaderKey, this.remoteRealPortHeaderKey),
	}
	vs := r.URL.Query()
	context.QueryValues = &vs

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
