package router

import (
	"reflect"

	"andals/gobox/http/controller"
)

type Route struct {
	Cl          controller.Controller
	ActionValue *reflect.Value
	Args        []string
}

type Router interface {
	MapRouteItems(cls ...controller.Controller)
	DefineRouteItem(pattern string, cl controller.Controller, actionName string)

	FindRoute(path string) *Route
}
