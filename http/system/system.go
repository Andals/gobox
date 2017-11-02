package system

import (
	"andals/gobox/http/controller"
	"andals/gobox/http/router"

	"net/http"
	"reflect"
)

type System struct {
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

		w.Write(context.ResponseBody())
		context.Destruct()
	}()

	context.BeforeAction()
	route.ActionValue.Call(this.makeArgsValues(context, route.Args))
	context.AfterAction()
}

func (this *System) makeArgsValues(context controller.ActionContext, args []string) []reflect.Value {
	argsValues := make([]reflect.Value, len(args)+1)
	argsValues[0] = reflect.ValueOf(context)
	for i, arg := range args {
		argsValues[i+1] = reflect.ValueOf(arg)
	}

	return argsValues
}

type JumpFunc func(context controller.ActionContext, args ...interface{})

type jumpItem struct {
	jf JumpFunc

	args []interface{}
}

func JumpOutAction(jf JumpFunc, args ...interface{}) {
	ji := &jumpItem{
		jf:   jf,
		args: args,
	}

	panic(ji)
}

func Redirect302(url string) {
	JumpOutAction(redirect302, url)
}

func redirect302(context controller.ActionContext, args ...interface{}) {
	http.Redirect(context.ResponseWriter(), context.Request(), args[0].(string), 302)
}
