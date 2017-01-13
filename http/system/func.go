package system

import (
	"net/http"

	"andals/gobox/http/controller"
)

type JumpFunc func(context controller.ActionContext, args ...interface{})

type jumpItem struct {
	jf JumpFunc

	args []interface{}
}

func LongJump(jf JumpFunc, args ...interface{}) {
	ji := &jumpItem{
		jf:   jf,
		args: args,
	}

	panic(ji)
}

func Redirect302(url string) {
	LongJump(redirect302, url)
}

func redirect302(context controller.ActionContext, args ...interface{}) {
	http.Redirect(context.ResponseWriter(), context.Request(), args[0].(string), 302)
}
