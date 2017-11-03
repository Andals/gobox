package router

import (
	"github.com/andals/gobox/http/controller"

	"reflect"
	"regexp"
	"strings"
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

func (this *SimpleRouter) FindRoute(path string) *Route {
	path = strings.ToLower(path)

	rg := this.findRouteGuideByDefined(path)
	if rg == nil {
		rg = this.findRouteGuideByGeneral(path)
	}

	ri, ok := this.routeTable[rg.controllerName]
	if !ok {
		return nil
	}

	ai, ok := ri.actionMap[rg.actionName]
	if !ok {
		return nil
	}

	return &Route{
		Cl:          ri.cl,
		ActionValue: ai.value,
		Args:        this.makeActionArgs(rg.actionArgs, ai.argsNum),
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

func (this *SimpleRouter) makeActionArgs(args []string, validArgsNum int) []string {
	rgArgsNum := len(args)
	missArgsNum := validArgsNum - rgArgsNum
	switch {
	case missArgsNum == 0:
	case missArgsNum > 0:
		for i := 0; i < missArgsNum; i++ {
			args = append(args, "")
		}
	case missArgsNum < 0:
		args = args[:validArgsNum]
	}

	return args
}
