package controller

import (
	"net/http"
)

type ActionContext interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	ResponseBody() []byte
}

type Controller interface {
	NewActionContext(req *http.Request, respWriter http.ResponseWriter) ActionContext

	BeforeAction(context ActionContext)
	AfterAction(context ActionContext)
	Destruct(context ActionContext)
}
