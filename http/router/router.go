package router

import (
	"reflect"

	ghttp "andals/gobox/http"
	"andals/gobox/http/controller"
)

type Route struct {
	Cl          controller.Controller
	ActionValue *reflect.Value
	ArgsValues  []reflect.Value
}

type Router interface {
	MapRouteItems(cls ...controller.Controller)
	DefineRouteItem(pattern string, cl controller.Controller, actionName string)

	FindRoute(context *ghttp.Context) *Route
}
