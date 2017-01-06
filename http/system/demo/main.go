package main

import (
	ghttp "andals/gobox/http"
	"andals/gobox/http/gracehttp"
	"andals/gobox/http/router"
	"andals/gobox/http/system"
)

func main() {
	dcl := new(DemoController)
	r := router.NewSimpleRouter()

	r.DefineRouteItem("^/g/([0-9]+)$", dcl, "get")
	r.MapRouteItems(new(IndexController), dcl)

	sys := system.NewSystem(r)

	gracehttp.ListenAndServe(":8001", sys)
}

type IndexController struct {
}

func (this *IndexController) BeforeAction(context *ghttp.Context) {
	context.RespBody = append(context.RespBody, []byte(" index before ")...)
}

func (this *IndexController) IndexAction(context *ghttp.Context) {
	context.RespBody = append(context.RespBody, []byte(" index action ")...)
}

func (this *IndexController) RedirectAction(context *ghttp.Context) {
	system.Redirect302("https://github.com/Andals/gobox")
}

func (this *IndexController) AfterAction(context *ghttp.Context) {
	context.RespBody = append(context.RespBody, []byte(" index after ")...)
}

func (this *IndexController) Destruct(context *ghttp.Context) {
	println(" index destruct ")
}

type DemoController struct {
}

func (this *DemoController) BeforeAction(context *ghttp.Context) {
	context.RespBody = append(context.RespBody, []byte(" demo before ")...)
}

func (this *DemoController) DemoAction(context *ghttp.Context) {
	context.RespBody = append(context.RespBody, []byte(" demo action ")...)
}

func (this *DemoController) GetAction(context *ghttp.Context, id string) {
	context.RespBody = append(context.RespBody, []byte(" get action id = "+id)...)
}

func (this *DemoController) AfterAction(context *ghttp.Context) {
	context.RespBody = append(context.RespBody, []byte(" demo after ")...)
}

func (this *DemoController) Destruct(context *ghttp.Context) {
	println(" demo destruct ")
}
