package system

import (
	"net/http"

	ghttp "andals/gobox/http"
	"andals/gobox/http/router"
)

const (
	DEF_REMOTE_REAL_IP_HEADER_KEY   = "REMOTE-REAL-IP"
	DEF_REMOTE_REAL_PORT_HEADER_KEY = "REMOTE-REAL-PORT"
)

type System struct {
	//eg, access by nginx's proxy_pass
	remoteRealIpHeaderKey   string
	remoteRealPortHeaderKey string

	router router.Router
}

func NewSystem(r router.Router) *System {
	return &System{
		remoteRealIpHeaderKey:   DEF_REMOTE_REAL_IP_HEADER_KEY,
		remoteRealPortHeaderKey: DEF_REMOTE_REAL_PORT_HEADER_KEY,

		router: r,
	}
}

func (this *System) SetRemoteRealIpHeaderKey(key string) *System {
	this.remoteRealIpHeaderKey = key

	return this
}

func (this *System) SetRemoteRealPortHeaderKey(key string) *System {
	this.remoteRealPortHeaderKey = key

	return this
}

func (this *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := ghttp.NewContext(r, w, ghttp.ParseRemoteAddr(r, this.remoteRealIpHeaderKey, this.remoteRealPortHeaderKey))

	this.dispatch(context)

	context.RespWriter.Write(context.RespBody)
}

func (this *System) dispatch(context *ghttp.Context) {
	r := this.router.FindRoute(context)
	if r == nil {
		error404(context)
		return
	}

	defer func() {
		r.Cl.Destruct(context)
	}()
	defer func() {
		if e := recover(); e != nil {
			ji, ok := e.(*jumpItem)
			if !ok {
				panic(e)
			}
			ji.jf(context, ji.args...)
		}
	}()

	r.Cl.BeforeAction(context)
	r.ActionValue.Call(r.ArgsValues)
	r.Cl.AfterAction(context)
}
