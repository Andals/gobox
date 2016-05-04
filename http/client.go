package http

import (
	"andals/gobox/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	logger log.ILogger

	c    *http.Client
	resp *http.Response
	t    time.Duration
}

func NewClient(logger log.ILogger) *Client {
	this := new(Client)

	this.logger = logger
	if this.logger == nil {
		this.logger = new(log.NoopLogger)
	}

	this.c = new(http.Client)

	return this
}

func (this *Client) SetTimeout(timeout time.Duration) *Client {
	if timeout != 0 {
		this.c.Timeout = timeout
	}

	return this
}

func (this *Client) Do(req *http.Request, retry int) ([]byte, error) {
	var err error

	start := time.Now()
	this.resp, err = this.c.Do(req)
	this.t = time.Since(start)
	if err != nil {
		for i := 0; i < retry; i++ {
			start = time.Now()
			this.resp, err = this.c.Do(req)
			this.t = time.Since(start)
		}
	}

	msg := "Method:" + req.Method + "\t" + "Url:" + req.URL.String() + "\t" + "Time:" + this.t.String() + "\n"
	this.logger.Info([]byte(msg))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(this.resp.Body)
}

func NewRequestForGet(url string, ip string, extHeaders map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return setRequestCommon(req, ip, extHeaders), nil
}

func setRequestCommon(req *http.Request, ip string, extHeaders map[string]string) *http.Request {
	req.Host = req.URL.Host

	if ip != "" {
		s := strings.Split(req.URL.Host, ":")
		s[0] = ip
		req.URL.Host = strings.Join(s, ":")
	}

	if extHeaders != nil {
		for k, v := range extHeaders {
			req.Header.Set(k, v)
		}
	}

	return req
}
