package redis

import (
	"github.com/garyburd/redigo/redis"

	"github.com/andals/golog"

	"fmt"
	"io"
)

type CmdLogFmtFunc func(cmd string, args ...interface{}) []byte

type cmdArgs struct {
	cmd  string
	args []interface{}
}

type Client struct {
	config *Config
	logger golog.ILogger
	clff   CmdLogFmtFunc

	conn      redis.Conn
	connected bool

	pipeCmds []*cmdArgs
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

		pipeCmds: []*cmdArgs{},
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
	if this.conn != nil {
		this.conn.Close()
	}

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

func (this *Client) Do(cmd string, args ...interface{}) *Reply {
	if !this.connected {
		if err := this.Connect(); err != nil {
			return NewReply(nil, err)
		}
	}

	this.log(cmd, args...)
	defer func() {
		this.pipeCmds = []*cmdArgs{}
	}()

	for _, ca := range this.pipeCmds {
		if err := this.conn.Send(ca.cmd, ca.args...); err != nil {
			return NewReply(nil, err)
		}
	}

	reply, err := this.conn.Do(cmd, args...)
	if err != nil {
		if err != io.EOF {
			return NewReply(nil, err)
		}
		if !this.config.TimeoutAutoReconnect {
			return NewReply(nil, err)
		}
		err = this.reconnect()
		if err != nil {
			return NewReply(nil, err)
		}

		for _, ca := range this.pipeCmds {
			if err = this.conn.Send(ca.cmd, ca.args...); err != nil {
				return NewReply(nil, err)
			}
		}
		reply, err = this.conn.Do(cmd, args...)
		if err != nil {
			return NewReply(nil, err)
		}
	}

	return NewReply(reply, err)
}

func (this *Client) Send(cmd string, args ...interface{}) {
	this.log(cmd, args...)
	this.pipeCmds = append(this.pipeCmds, &cmdArgs{cmd, args})
}

func (this *Client) ExecPipelining() ([]*Reply, []int) {
	if !this.connected {
		if err := this.Connect(); err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}

	defer func() {
		this.pipeCmds = []*cmdArgs{}
	}()

	for _, ca := range this.pipeCmds {
		if err := this.conn.Send(ca.cmd, ca.args...); err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}
	if err := this.conn.Flush(); err != nil {
		return []*Reply{NewReply(nil, err)}, []int{0}
	}

	reply, err := this.conn.Receive()
	if err != nil {
		if err != io.EOF {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
		if !this.config.TimeoutAutoReconnect {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
		err = this.reconnect()
		if err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}

		for _, ca := range this.pipeCmds {
			if err = this.conn.Send(ca.cmd, ca.args...); err != nil {
				return []*Reply{NewReply(nil, err)}, []int{0}
			}
		}

		if err = this.conn.Flush(); err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
		reply, err = this.conn.Receive()
		if err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}

	replies := make([]*Reply, len(this.pipeCmds))
	var errIndexes []int

	replies[0] = NewReply(reply, nil)
	for i := 1; i < len(this.pipeCmds); i++ {
		reply, err := this.conn.Receive()
		replies[i] = NewReply(reply, err)
		if err != nil {
			errIndexes = append(errIndexes, i)
		}
	}

	return replies, errIndexes
}

func (this *Client) BeginTrans() {
	this.Send("multi")
}

func (this *Client) DiscardTrans() error {
	return this.Do("discard").Err
}

func (this *Client) ExecTrans() ([]*Reply, error) {
	reply := this.Do("exec")
	values, err := redis.Values(reply.reply, reply.Err)
	if err != nil {
		return nil, err
	}

	replies := make([]*Reply, len(values))
	for i, value := range values {
		replies[i] = NewReply(value, nil)
	}

	return replies, nil
}

func (this *Client) log(cmd string, args ...interface{}) {
	if len(cmd) == 0 {
		return
	}

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

func (this *Client) reconnect() error {
	this.Free()

	return this.Connect()
}
