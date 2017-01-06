package system

import (
	"net/http"

	ghttp "andals/gobox/http"
)

type JumpFunc func(context *ghttp.Context, args ...interface{})

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

func redirect302(context *ghttp.Context, args ...interface{}) {
	http.Redirect(context.RespWriter, context.Req, args[0].(string), 302)
}

func error404(context *ghttp.Context) {
	http.NotFound(context.RespWriter, context.Req)
}
