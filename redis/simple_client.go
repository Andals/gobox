package redis

import (
	//"fmt"
	"time"

	"andals/gobox/log"
	"andals/gobox/misc"

	"github.com/fzzy/radix/redis"
)

type simpleClient struct {
	network string
	addr    string
	pass    string
	timeout time.Duration

	client *redis.Client
	logger log.ILogger
}

func NewSimpleClient(network, addr, pass string, timeout time.Duration) (*simpleClient, error) {
	client, err := newRedisClient(network, addr, pass, timeout)
	if err != nil {
		return nil, err
	}

	this := &simpleClient{
		network: network,
		addr:    addr,
		pass:    pass,
		timeout: timeout,

		client: client,
		logger: new(log.NoopLogger),
	}

	return this, nil
}

func (this *simpleClient) SetLogger(logger log.ILogger) {
	this.logger = logger
}

func (this *simpleClient) Free() {
	this.client.Close()
	this.logger.Free()
}

func (this *simpleClient) runCmd(cmd string, args ...string) *redis.Reply {
	cmdBytes := []byte(cmd)
	for _, s := range args {
		cmdBytes = misc.AppendBytes(cmdBytes, []byte(" "), []byte(s))
	}

	this.logger.Info(cmdBytes)

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

func (this *simpleClient) Expire(key, seconds string) error {
	r := this.runCmd("EXPIRE", key, seconds)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *simpleClient) Del(key string) error {
	r := this.runCmd("DEL", key)
	if r.Err != nil {
		return r.Err
	}

	return nil
}
