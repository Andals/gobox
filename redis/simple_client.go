package redis

import (
	//"fmt"
	"time"

	"andals/gobox/log"
	"andals/gobox/misc"

	"github.com/fzzy/radix/redis"
)

type SimpleClient struct {
	network string
	addr    string
	pass    string
	timeout time.Duration

	lastCmd []byte

	client *redis.Client
	logger log.ILogger
}

func NewSimpleClient(network, addr, pass string, timeout time.Duration, logger log.ILogger) (*SimpleClient, error) {
	client, err := newRedisClient(network, addr, pass, timeout)
	if err != nil {
		return nil, err
	}

	this := &SimpleClient{
		network: network,
		addr:    addr,
		pass:    pass,
		timeout: timeout,

		client: client,
		logger: logger,
	}

	return this, nil
}

func (this *SimpleClient) LastCmd() []byte {
	return this.lastCmd
}

func (this *SimpleClient) Free() {
	this.client.Close()
	this.logger.Free()
}

func (this *SimpleClient) runCmd(cmd string, args ...string) *redis.Reply {
	this.lastCmd = []byte(cmd)
	for _, s := range args {
		this.lastCmd = misc.AppendBytes(this.lastCmd, []byte(" "), []byte(s))
	}

	this.logger.Info(this.lastCmd)

	cnt := len(args)
	iargs := make([]interface{}, cnt)
	for i := 0; i < cnt; i++ {
		iargs[i] = args[i]
	}

	r := this.client.Cmd(cmd, iargs...)
	if r.Err == nil {
		return r
	}

	this.client.Close()
	c, e := newRedisClient(this.network, this.addr, this.pass, this.timeout)
	if e != nil {
		return r
	}

	this.client = c
	return this.client.Cmd(cmd, iargs...)
}
