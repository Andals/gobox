package client

import (
	"andals/gobox/log"
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	logger log.ILogger

	c *http.Client
}

type Response struct {
	T        time.Duration
	Contents []byte

	*http.Response
}

func NewClient(logger log.ILogger) *Client {
	this := new(Client)

	this.logger = logger
	if this.logger == nil {
		this.logger = new(log.NoopLogger)
	}

	this.c = new(http.Client)
	this.c.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return this
}

func (this *Client) SetTimeout(timeout time.Duration) *Client {
	if timeout != 0 {
		this.c.Timeout = timeout
	}

	return this
}

func (this *Client) SetMaxIdleConnsPerHost(value int) *Client {
	if value != 0 {
		this.c.Transport.(*http.Transport).MaxIdleConnsPerHost = value
	}

	return this
}

func (this *Client) Do(req *http.Request, retry int) (*Response, error) {
	start := time.Now()
	resp, err := this.c.Do(req)
	t := time.Since(start)
	if err != nil || resp.StatusCode != 200 {
		for i := 0; i < retry; i++ {
			start = time.Now()
			resp, err = this.c.Do(req)
			t = time.Since(start)
			if err == nil && resp.StatusCode == 200 {
				break
			}
		}
	}

	msg := [][]byte{
		[]byte("Method:" + req.Method),
		[]byte("Host: " + req.Host),
		[]byte("Url:" + req.URL.String()),
		[]byte("Time:" + t.String()),
	}
	if err != nil {
		if resp != nil {
			msg = append(msg, []byte("StatusCode:"+strconv.Itoa(resp.StatusCode)))
		}
		msg = append(msg, []byte("ErrMsg:"+err.Error()+"\n"))
		this.logger.Error(bytes.Join(msg, []byte("\t")))
		return nil, err
	}
	msg = append(msg, []byte("StatusCode:"+strconv.Itoa(resp.StatusCode)+"\n"))
	this.logger.Info(bytes.Join(msg, []byte("\t")))

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		T:        t,
		Contents: contents,
		Response: resp,
	}, nil
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
