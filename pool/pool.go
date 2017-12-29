package pool

import (
	"errors"
	"time"
)

type IConn interface {
	Free()
}

type Config struct {
	Size              int
	MaxIdleTime       time.Duration
	KeepAliveInterval time.Duration

	NewConnFunc   func() (IConn, error)
	KeepAliveFunc func(conn IConn) error
}

type poolItem struct {
	conn IConn

	accessTime time.Time
}

var ErrPoolIsFull = errors.New("pool is full")

type Pool struct {
	config *Config

	conns chan *poolItem
}

func NewPool(config *Config) *Pool {
	this := &Pool{
		config: config,

		conns: make(chan *poolItem, config.Size),
	}

	if config.KeepAliveInterval > 0 && config.KeepAliveFunc != nil {
		go this.keepAliveRoutine()
	}

	return this
}

func (this *Pool) Get() (IConn, error) {
	pi := this.get()
	if pi != nil {
		if time.Now().Sub(pi.accessTime) < this.config.MaxIdleTime {
			return pi.conn, nil
		}

		pi.conn.Free()
	}

	return this.config.NewConnFunc()
}

func (this *Pool) Put(conn IConn) error {
	pi := &poolItem{
		conn: conn,

		accessTime: time.Now(),
	}

	notFull := this.put(pi)
	if notFull {
		return nil
	}

	conn.Free()

	return ErrPoolIsFull
}

func (this *Pool) get() *poolItem {
	select {
	case pi := <-this.conns:
		return pi
	default:
	}

	return nil
}

func (this *Pool) put(pi *poolItem) bool {
	select {
	case this.conns <- pi:
		return true
	default:
	}

	return false
}

func (this *Pool) keepAliveRoutine() {
	ticker := time.NewTicker(this.config.KeepAliveInterval)

	for {
		select {
		case <-ticker.C:
			this.keepAlive()
		}
	}
}

func (this *Pool) keepAlive() {
	maxIdleNum := len(this.conns)

	for i := 0; i < maxIdleNum; i++ {
		pi := this.get()
		if pi == nil {
			return
		}

		if time.Now().Sub(pi.accessTime) < this.config.MaxIdleTime {
			err := this.config.KeepAliveFunc(pi.conn)
			if err == nil {
				if this.put(pi) {
					continue
				}
			}
		}

		pi.conn.Free()
	}
}
