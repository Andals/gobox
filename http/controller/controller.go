package controller

import (
	ghttp "andals/gobox/http"
)

type Controller interface {
	BeforeAction(context *ghttp.Context)
	AfterAction(context *ghttp.Context)
	Destruct(context *ghttp.Context)
}
