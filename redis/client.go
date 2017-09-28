package redis

import (
	"github.com/garyburd/redigo/redis"

	"andals/golog"

	"fmt"
)

type Client struct {
	config *Config
	logger golog.ILogger

	conn redis.Conn
}

func NewClient(config *Config, logger golog.ILogger) (*Client, error) {
	if config.LogLevel == 0 {
		config.LogLevel = golog.LEVEL_INFO
	}

	options := []redis.DialOption{
		redis.DialConnectTimeout(config.ConnectTimeout),
		redis.DialReadTimeout(config.ReadTimeout),
		redis.DialWriteTimeout(config.WriteTimeout),
	}

	conn, err := redis.Dial("tcp", config.Host+":"+config.Port, options...)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("auth", config.Pass)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	return &Client{
		config: config,
		logger: logger,

		conn: conn,
	}, nil
}

func (this *Client) SetLogger(logger golog.ILogger) {
	this.logger = logger
}

func (this *Client) Free() {
	this.logger.Free()
	this.conn.Close()
}

func (this *Client) Do(cmd string, args ...interface{}) (*Reply, error) {
	this.log(cmd, args...)

	reply, err := this.conn.Do(cmd, args...)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, nil
	}

	return &Reply{reply}, nil
}

func (this *Client) Send(cmd string, args ...interface{}) error {
	this.log(cmd, args...)

	return this.conn.Send(cmd, args...)
}

func (this *Client) ExecPipelining() ([]*Reply, error) {
	return this.multiDo("")
}

func (this *Client) BeginTrans() error {
	cmd := "multi"
	this.log(cmd)

	return this.conn.Send(cmd)
}

func (this *Client) DiscardTrans() error {
	cmd := "discard"
	this.log(cmd)

	_, err := this.conn.Do(cmd)
	return err
}

func (this *Client) ExecTrans() ([]*Reply, error) {
	cmd := "exec"
	this.log(cmd)

	return this.multiDo(cmd)
}

func (this *Client) log(cmd string, args ...interface{}) {
	for _, arg := range args {
		cmd += " " + fmt.Sprint(arg)
	}

	this.logger.Log(this.config.LogLevel, []byte(cmd))
}

func (this *Client) multiDo(cmd string) ([]*Reply, error) {
	r, err := this.conn.Do(cmd)
	if err != nil {
		return nil, err
	}

	rs := r.([]interface{})
	replies := make([]*Reply, len(rs))
	for i, v := range rs {
		replies[i] = &Reply{v}
	}

	return replies, nil
}
