package main

import (
	"andals/gobox/http/controller"
	"andals/gobox/http/gracehttp"
	"andals/gobox/misc"
)

func main() {
	cl := controller.NewController()

	cl.BeforeAction(beforeAction)
	cl.AfterAction(afterAction)

	cl.ExactMatchAction("/exact", exactAction)
	cl.RegexMatchAction("^/[a-z]+([0-9]+)", regexAction)

	gracehttp.ListenAndServe(":8001", cl)
}

func beforeAction(context *controller.Context, args []string) {
	context.RespBody = []byte("before")
}

func exactAction(context *controller.Context, args []string) {
	context.RespBody = misc.AppendBytes(context.RespBody, []byte(" exact "))
}

func regexAction(context *controller.Context, args []string) {
	context.RespBody = misc.AppendBytes(context.RespBody, []byte(" regex id = "+args[0]+" "))
}

func afterAction(context *controller.Context, args []string) {
	context.RespBody = misc.AppendBytes(context.RespBody, []byte("after\n"))
}
