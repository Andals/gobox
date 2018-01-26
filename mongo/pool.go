package mongo

import (
	"github.com/andals/gobox/pool"
)

type PConfig struct {
	pool.Config

	NewClientFunc func() (*Client, error)
}

type Pool struct {
	p *pool.Pool

	config *PConfig
}

func NewPool(config *PConfig) *Pool {
	this := &Pool{
		config: config,
	}

	config.NewConnFunc = this.newConn

	this.p = pool.NewPool(&this.config.Config)

	return this
}

func (this *Pool) Get() (*Client, error) {
	conn, err := this.p.Get()
	if err != nil {
		return nil, err
	}

	return conn.(*Client), nil
}

func (this *Pool) Put(client *Client) error {
	return this.p.Put(client)
}

func (this *Pool) newConn() (pool.IConn, error) {
	return this.config.NewClientFunc()
}
