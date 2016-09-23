package main

import (
	"andals/gobox/http/controller"
	"andals/gobox/http/gracehttp"
	"andals/gobox/misc"

	"fmt"
)

func main() {
	cl := controller.NewController()

	cl.AddBeforeAction("^/exact", beforeAction).
		AddAfterAction("^/([a-z]+)[0-9]+", afterAction).
		AddExactMatchAction("/exact", exactAction).
		AddExactMatchAction("/redirect", redirectAction).
		AddRegexMatchAction("^/[a-z]+([0-9]+)", regexAction).
		AddDestructFunc("^/([a-z]+)[0-9]+", destruct)

	gracehttp.ListenAndServe(":8001", cl)
}

func beforeAction(context *controller.Context, args []string) {
	context.RespBody = []byte("exact before")
}

func exactAction(context *controller.Context, args []string) {
	context.RespBody = misc.AppendBytes(context.RespBody, []byte(" exact "))
}

func regexAction(context *controller.Context, args []string) {
	context.RespBody = misc.AppendBytes(context.RespBody, []byte(" regex id = "+args[0]+" "))
}

func afterAction(context *controller.Context, args []string) {
	context.RespBody = misc.AppendBytes(context.RespBody, []byte("after "+args[0]+"\n"))

	context.TransData["after"] = "after"
}

func redirectAction(context *controller.Context, args []string) {
	controller.Redirect302("http://www.gobox.com")
}

func destruct(transData map[string]interface{}, args []string) {
	fmt.Println(transData, args)
}
