package system

import (
	"net/http"
	"reflect"

	"andals/gobox/http/controller"
	"andals/gobox/http/router"
)

type System struct {
	//eg, access by nginx's proxy_pass
	remoteRealIpHeaderKey   string
	remoteRealPortHeaderKey string

	router router.Router
}

func NewSystem(r router.Router) *System {
	return &System{
		router: r,
	}
}

func (this *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := this.router.FindRoute(r.URL.Path)
	if route == nil {
		http.NotFound(w, r)
		return
	}

	context := route.Cl.NewActionContext(r, w)

	defer func() {
		if e := recover(); e != nil {
			ji, ok := e.(*jumpItem)
			if !ok {
				panic(e)
			}
			ji.jf(context, ji.args...)
		}

		route.Cl.Destruct(context)
	}()

	route.Cl.BeforeAction(context)
	route.ActionValue.Call(this.makeArgsValues(context, route.Args))
	route.Cl.AfterAction(context)

	w.Write(context.ResponseBody())
}

func (this *System) makeArgsValues(context controller.ActionContext, args []string) []reflect.Value {
	argsValues := make([]reflect.Value, len(args)+1)
	argsValues[0] = reflect.ValueOf(context)
	for i, arg := range args {
		argsValues[i+1] = reflect.ValueOf(arg)
	}

	return argsValues
}
