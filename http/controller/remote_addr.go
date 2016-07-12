package controller

import (
	"net/http"
	"strings"
)

type RemoteAddr struct {
	Ip   string
	Port string
}

func ParseRemoteAddr(r *http.Request, remoteRealIpHeaderKey, remoteRealPortHeaderKey string) *RemoteAddr {
	rs := strings.Split(r.RemoteAddr, ":")
	this := &RemoteAddr{
		Ip:   rs[0],
		Port: rs[1],
	}

	if this.Ip == "127.0.0.1" {
		rip := strings.TrimSpace(r.Header.Get(remoteRealIpHeaderKey))
		rport := strings.TrimSpace(r.Header.Get(remoteRealPortHeaderKey))
		if rip != "" && rport != "" {
			this.Ip = rip
			this.Port = rport
		}
	}

	return this
}

func (this *RemoteAddr) String() string {
	return this.Ip + ":" + this.Port
}
