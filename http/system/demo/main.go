package main

import (
	"andals/gobox/http/controller"
	"andals/gobox/http/gracehttp"
	"andals/gobox/http/router"
	"andals/gobox/http/system"
	"net/http"
)

func main() {
	dcl := new(DemoController)
	r := router.NewSimpleRouter()

	r.DefineRouteItem("^/g/([0-9]+)$", dcl, "get")
	r.MapRouteItems(new(IndexController), dcl)

	sys := system.NewSystem(r)

	gracehttp.ListenAndServe(":8001", sys)
}

type DemoActionContext struct {
	Req        *http.Request
	RespWriter http.ResponseWriter
	RespBody   []byte
}

func (this *DemoActionContext) Request() *http.Request {
	return this.Req
}

func (this *DemoActionContext) ResponseWriter() http.ResponseWriter {
	return this.RespWriter
}

func (this *DemoActionContext) ResponseBody() []byte {
	return this.RespBody
}

type IndexController struct {
}

func (this *IndexController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &DemoActionContext{
		Req:        req,
		RespWriter: respWriter,
	}
}

func (this *IndexController) BeforeAction(context controller.ActionContext) {
	acontext := context.(*DemoActionContext)

	acontext.RespBody = append(acontext.RespBody, []byte(" index before ")...)
}

func (this *IndexController) IndexAction(context *DemoActionContext) {
	context.RespBody = append(context.RespBody, []byte(" index action ")...)
}

func (this *IndexController) RedirectAction(context *DemoActionContext) {
	system.Redirect302("https://github.com/Andals/gobox")
}

func (this *IndexController) AfterAction(context controller.ActionContext) {
	acontext := context.(*DemoActionContext)

	acontext.RespBody = append(acontext.RespBody, []byte(" index after ")...)
}

func (this *IndexController) Destruct(context controller.ActionContext) {
	println(" index destruct ")
}

type DemoController struct {
}

func (this *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &DemoActionContext{
		Req:        req,
		RespWriter: respWriter,
	}
}

func (this *DemoController) BeforeAction(context controller.ActionContext) {
	acontext := context.(*DemoActionContext)

	acontext.RespBody = append(acontext.RespBody, []byte(" demo before ")...)
}

func (this *DemoController) DemoAction(context *DemoActionContext) {
	context.RespBody = append(context.RespBody, []byte(" demo action ")...)
}

func (this *DemoController) GetAction(context *DemoActionContext, id string) {
	context.RespBody = append(context.RespBody, []byte(" get action id = "+id)...)
}

func (this *DemoController) AfterAction(context controller.ActionContext) {
	acontext := context.(*DemoActionContext)

	acontext.RespBody = append(acontext.RespBody, []byte(" demo after ")...)
}

func (this *DemoController) Destruct(context controller.ActionContext) {
	println(" demo destruct ")
}
