package redis

import (
	"github.com/garyburd/redigo/redis"

	"github.com/andals/golog"

	"fmt"
)

type CmdLogFmtFunc func(cmd string, args ...interface{}) []byte

type Client struct {
	config *Config
	logger golog.ILogger
	clff   CmdLogFmtFunc

	conn      redis.Conn
	connected bool
}

func NewClient(config *Config, logger golog.ILogger) *Client {
	if config.LogLevel == 0 {
		config.LogLevel = golog.LEVEL_INFO
	}

	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	this := &Client{
		config: config,
		logger: logger,
	}
	this.clff = this.cmdLogFmt

	return this
}

func (this *Client) SetLogger(logger golog.ILogger) *Client {
	this.logger = logger

	return this
}

func (this *Client) SetCmdLogFmtFunc(clff CmdLogFmtFunc) *Client {
	this.clff = clff

	return this
}

func (this *Client) Connected() bool {
	return this.connected
}

func (this *Client) Free() {
	this.conn.Close()
	this.connected = false
}

func (this *Client) Connect() error {
	options := []redis.DialOption{
		redis.DialConnectTimeout(this.config.ConnectTimeout),
		redis.DialReadTimeout(this.config.ReadTimeout),
		redis.DialWriteTimeout(this.config.WriteTimeout),
	}

	conn, err := redis.Dial("tcp", this.config.Host+":"+this.config.Port, options...)
	if err != nil {
		return err
	}

	_, err = conn.Do("auth", this.config.Pass)
	if err != nil {
		return err
	}

	this.conn = conn
	this.connected = true

	return nil
}

func (this *Client) Do(cmd string, args ...interface{}) (*Reply, error) {
	if !this.connected {
		if err := this.Connect(); err != nil {
			return nil, err
		}
	}

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
	if !this.connected {
		if err := this.Connect(); err != nil {
			return err
		}
	}

	this.log(cmd, args...)

	return this.conn.Send(cmd, args...)
}

func (this *Client) ExecPipelining() ([]*Reply, error) {
	return this.multiDo("")
}

func (this *Client) BeginTrans() error {
	return this.Send("multi")
}

func (this *Client) DiscardTrans() error {
	_, err := this.Do("discard")

	return err
}

func (this *Client) ExecTrans() ([]*Reply, error) {
	return this.multiDo("exec")
}

func (this *Client) log(cmd string, args ...interface{}) {
	msg := this.clff(cmd, args...)
	if msg != nil {
		this.logger.Log(this.config.LogLevel, msg)
	}
}

func (this *Client) cmdLogFmt(cmd string, args ...interface{}) []byte {
	for _, arg := range args {
		cmd += " " + fmt.Sprint(arg)
	}

	return []byte(cmd)
}

func (this *Client) multiDo(cmd string) ([]*Reply, error) {
	if !this.connected {
		if err := this.Connect(); err != nil {
			return nil, err
		}
	}

	this.log(cmd)

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
