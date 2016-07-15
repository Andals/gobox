package controller

import (
	//     "fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"andals/gobox/encoding"
	"andals/gobox/misc"
)

type Context struct {
	Req        *http.Request
	RespWriter http.ResponseWriter

	QueryValues *url.Values
	TransData   map[string]interface{}
	RespBody    []byte

	RemoteRealAddr *RemoteAddr
	Rid            []byte
}

func NewContext(r *http.Request, w http.ResponseWriter, remoteRealAddr *RemoteAddr) *Context {
	this := &Context{
		RespWriter: w,
		Req:        r,

		TransData:      make(map[string]interface{}),
		RemoteRealAddr: remoteRealAddr,
	}
	vs := r.URL.Query()
	this.QueryValues = &vs

	now := time.Now()
	timeInt := now.UnixNano()
	randInt := misc.RandByTime(&now)

	ridStr := this.RemoteRealAddr.String() + "," + strconv.FormatInt(timeInt, 10) + "," + strconv.FormatInt(randInt, 10)
	this.Rid = encoding.Base64Encode([]byte(ridStr))

	return this
}
