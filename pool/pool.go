package pool

import (
	"errors"
	"time"
)

type IConn interface {
	Free()
}

type NewConnFunc func() (IConn, error)

type poolItem struct {
	conn IConn

	addTime int64
}

type Pool struct {
	connTimeout int64

	conns chan *poolItem
	ncf   NewConnFunc
}

func NewPool(connTimeout time.Duration, size int, ncf NewConnFunc) *Pool {
	this := &Pool{
		connTimeout: int64(connTimeout),

		conns: make(chan *poolItem, size),
		ncf:   ncf,
	}

	return this
}

func (this *Pool) Get() (IConn, error) {
	select {
	case pItem := <-this.conns:
		if time.Now().UnixNano()-pItem.addTime < this.connTimeout {
			return pItem.conn, nil
		}
		pItem.conn.Free()
	default:
	}

	return this.ncf()
}

func (this *Pool) Put(conn IConn) error {
	pItem := &poolItem{
		conn: conn,

		addTime: time.Now().UnixNano(),
	}

	select {
	case this.conns <- pItem:
		return nil
	default:
		conn.Free()
		pItem = nil
	}

	return errors.New("pool is full")
}
