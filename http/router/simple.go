package router

import (
	"reflect"
	"regexp"
	"strings"

	ghttp "andals/gobox/http"
	"andals/gobox/http/controller"
)

type actionItem struct {
	argsNum int
	value   *reflect.Value
}

type routeItem struct {
	cl  controller.Controller
	clv *reflect.Value
	clt reflect.Type

	controllerName string
	actionMap      map[string]*actionItem
}

type routeDefined struct {
	regex *regexp.Regexp

	controllerName string
	actionName     string
}

type routeGuide struct {
	controllerName string
	actionName     string
	actionArgs     []string
}

type SimpleRouter struct {
	defaultControllerName string
	defaultActionName     string

	cregex *regexp.Regexp
	aregex *regexp.Regexp

	routeDefinedList []*routeDefined
	routeTable       map[string]*routeItem
}

func NewSimpleRouter() *SimpleRouter {
	return &SimpleRouter{
		defaultActionName:     "index",
		defaultControllerName: "index",

		cregex: regexp.MustCompile("([A-Z][A-Za-z0-9_]*)Controller$"),
		aregex: regexp.MustCompile("^([A-Z][A-Za-z0-9_]*)Action$"),

		routeTable: make(map[string]*routeItem),
	}
}

func (this *SimpleRouter) SetDefaultControllerName(name string) *SimpleRouter {
	this.defaultControllerName = name

	return this
}

func (this *SimpleRouter) SetDefaultActionName(name string) *SimpleRouter {
	this.defaultActionName = name

	return this
}

func (this *SimpleRouter) MapRouteItems(cls ...controller.Controller) {
	for _, cl := range cls {
		this.mapRouteItem(cl)
	}
}

func (this *SimpleRouter) mapRouteItem(cl controller.Controller) {
	ri := this.getRouteItem(cl)
	if ri == nil {
		return
	}

	for i := 0; i < ri.clv.NumMethod(); i++ {
		am := ri.clt.Method(i)
		actionName := this.getActionName(am.Name)
		if actionName == "" {
			continue
		}
		_, ok := ri.actionMap[actionName]
		if ok {
			continue
		}
		actionArgsNum := this.getActionArgsNum(am, ri.clt)
		if actionArgsNum == -1 {
			continue
		}

		av := ri.clv.Method(i)
		ri.actionMap[actionName] = &actionItem{
			argsNum: actionArgsNum,
			value:   &av,
		}
	}
}

func (this *SimpleRouter) DefineRouteItem(pattern string, cl controller.Controller, actionName string) {
	methodName := strings.Title(actionName) + "Action"
	actionName = strings.ToLower(methodName)
	if actionName == "" {
		return
	}

	ri := this.getRouteItem(cl)
	if ri == nil {
		return
	}

	am, ok := ri.clt.MethodByName(methodName)
	if !ok {
		return
	}
	actionArgsNum := this.getActionArgsNum(am, ri.clt)
	if actionArgsNum == -1 {
		return
	}

	av := ri.clv.MethodByName(methodName)
	ri.actionMap[actionName] = &actionItem{
		argsNum: actionArgsNum,
		value:   &av,
	}

	this.routeDefinedList = append(this.routeDefinedList, &routeDefined{
		regex: regexp.MustCompile(pattern),

		controllerName: strings.ToLower(ri.controllerName),
		actionName:     strings.ToLower(actionName),
	})
}

func (this *SimpleRouter) getRouteItem(cl controller.Controller) *routeItem {
	v := reflect.ValueOf(cl)
	t := v.Type()

	controllerName := this.getControllerName(t.String())
	if controllerName == "" {
		return nil
	}

	ri, ok := this.routeTable[controllerName]
	if !ok {
		ri = &routeItem{
			cl:  cl,
			clv: &v,
			clt: t,

			controllerName: controllerName,
			actionMap:      make(map[string]*actionItem),
		}
		this.routeTable[controllerName] = ri
	}

	return ri
}

func (this *SimpleRouter) getControllerName(typeString string) string {
	matches := this.cregex.FindStringSubmatch(typeString)
	if matches == nil {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (this *SimpleRouter) getActionName(methodName string) string {
	matches := this.aregex.FindStringSubmatch(methodName)
	if matches == nil {
		return ""
	}

	actionName := strings.ToLower(matches[1])
	if actionName != "before" && actionName != "after" {
		return actionName
	}

	return ""
}

func (this *SimpleRouter) getActionArgsNum(actionMethod reflect.Method, controllerType reflect.Type) int {
	n := actionMethod.Type.NumIn()
	if n < 2 {
		return -1
	}

	if actionMethod.Type.In(0).String() != controllerType.String() {
		return -1
	}
	if actionMethod.Type.In(1).String() != "*http.Context" {
		return -1
	}
	if n > 2 {
		valid := true
		for i := 2; i < n; i++ {
			if actionMethod.Type.In(i).String() != "string" {
				valid = false
				break
			}
		}
		if !valid {
			return -1
		}
	}

	return n - 2 //delete this and context
}

func (this *SimpleRouter) FindRoute(context *ghttp.Context) *Route {
	path := strings.ToLower(context.Req.URL.Path)

	rg := this.findRouteGuideByDefined(path)
	if rg == nil {
		rg = this.findRouteGuideByGeneral(path)
	}

	ri, ok := this.routeTable[rg.controllerName]
	if !ok {
		return nil
	}

	actionItem, ok := ri.actionMap[rg.actionName]
	if !ok {
		return nil
	}

	return &Route{
		Cl:          ri.cl,
		ActionValue: actionItem.value,
		ArgsValues:  this.makeArgsValues(rg.actionArgs, actionItem.argsNum, context),
	}
}

func (this *SimpleRouter) findRouteGuideByDefined(path string) *routeGuide {
	for _, rd := range this.routeDefinedList {
		matches := rd.regex.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		return &routeGuide{
			controllerName: rd.controllerName,
			actionName:     rd.actionName,
			actionArgs:     matches[1:],
		}
	}

	return nil
}

func (this *SimpleRouter) findRouteGuideByGeneral(path string) *routeGuide {
	rg := new(routeGuide)

	path = strings.Trim(path, "/")
	sl := strings.Split(path, "/")

	sl[0] = strings.TrimSpace(sl[0])
	if sl[0] == "" {
		rg.controllerName = this.defaultControllerName
		rg.actionName = this.defaultActionName
	} else {
		rg.controllerName = sl[0]
		if len(sl) > 1 {
			sl[1] = strings.TrimSpace(sl[1])
			if sl[1] != "" {
				rg.actionName = sl[1]
			} else {
				rg.actionName = this.defaultActionName
			}
		} else {
			rg.actionName = this.defaultActionName
		}
	}

	return rg
}

func (this *SimpleRouter) makeArgsValues(routeGuideActionArgs []string, validArgsNum int, context *ghttp.Context) []reflect.Value {
	rgArgsNum := len(routeGuideActionArgs)
	missArgsNum := validArgsNum - rgArgsNum
	switch {
	case missArgsNum == 0:
	case missArgsNum > 0:
		for i := 0; i < missArgsNum; i++ {
			routeGuideActionArgs = append(routeGuideActionArgs, "")
		}
	case missArgsNum < 0:
		routeGuideActionArgs = routeGuideActionArgs[:validArgsNum]
	}

	argsValues := make([]reflect.Value, validArgsNum+1)
	argsValues[0] = reflect.ValueOf(context)
	for i, arg := range routeGuideActionArgs {
		argsValues[i+1] = reflect.ValueOf(arg)
	}

	return argsValues
}
