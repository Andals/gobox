package mysql

import (
	"github.com/andals/gobox/pool"

	"time"
)

type Pool struct {
	p *pool.Pool

	ncf NewClientFunc
}

type NewClientFunc func() (*Client, error)

func NewPool(clientTimeout time.Duration, size int, ncf NewClientFunc) *Pool {
	this := &Pool{
		ncf: ncf,
	}

	this.p = pool.NewPool(clientTimeout, size, this.newConn)

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
	return this.ncf()
}
