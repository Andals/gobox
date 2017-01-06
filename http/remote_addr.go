package http

import (
	"net/http"
	"strings"
)

type RemoteAddr struct {
	Ip   string
	Port string
}

func ParseRemoteAddr(r *http.Request, remoteRealIpHeaderKey, remoteRealPortHeaderKey string) *RemoteAddr {
	this := new(RemoteAddr)

	rip := strings.TrimSpace(r.Header.Get(remoteRealIpHeaderKey))
	rport := strings.TrimSpace(r.Header.Get(remoteRealPortHeaderKey))

	if rip != "" && rport != "" {
		this.Ip = rip
		this.Port = rport
	} else {
		rs := strings.Split(r.RemoteAddr, ":")
		this.Ip = rs[0]
		this.Port = rs[1]
	}

	return this
}

func (this *RemoteAddr) String() string {
	return this.Ip + ":" + this.Port
}
