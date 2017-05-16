package Redis

import (
	//"fmt"
	"time"

	"andals/gobox/log"
	"andals/gobox/misc"
	"github.com/fzzy/radix/redis"
)

type Client struct {
	client *redis.Client

	logger log.ILogger
}

func NewClient(network, addr, pass string, timeout time.Duration, logger log.ILogger) (*Client, error) {
	c, e := redis.DialTimeout(network, addr, timeout)
	if e != nil {
		return nil, e
	}

	this := &Client{
		client: c,
		logger: logger,
	}
	r := this.client.Cmd("AUTH", pass)
	if r.Err != nil {
		this.client.Close()

		return nil, r.Err
	}

	return this, nil
}

func (this *Client) Close() {
	this.client.Close()
}

func (this *Client) runCmd(cmd string, args ...string) *redis.Reply {
	this.logCmd(cmd, args...)

	cnt := len(args)
	iargs := make([]interface{}, cnt)
	for i := 0; i < cnt; i++ {
		iargs[i] = args[i]
	}

	return this.client.Cmd(cmd, iargs...)
}
func (this *Client) logCmd(cmd string, args ...string) {
	msg := []byte(cmd)
	for _, s := range args {
		msg = misc.AppendBytes(msg, []byte(" "), []byte(s))
	}

	this.logger.Info(msg)
}
