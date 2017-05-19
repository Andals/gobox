package redis

import (
	"andals/gobox/log"

	"time"
)

type poolItem struct {
	client *SimpleClient

	addTime int64
}

type Pool struct {
	clientTimeout int64

	ch chan *poolItem
}

func NewPool(clientTimeout time.Duration, size int) *Pool {
	this := &Pool{
		clientTimeout: int64(clientTimeout),

		ch: make(chan *poolItem, size),
	}

	return this
}

func (this *Pool) get() *SimpleClient {
	select {
	case pItem := <-this.ch:
		if time.Now().Unix()-pItem.addTime < this.clientTimeout {
			return pItem.client
		}
	default:
	}

	return nil
}

func (this *Pool) put(client *SimpleClient) bool {
	pItem := &poolItem{
		client: client,

		addTime: time.Now().Unix(),
	}

	select {
	case this.ch <- pItem:
		return true
	default:
	}

	return false
}

type PoolClient struct {
	client *SimpleClient
	pool   *Pool

	logger log.ILogger
}

func NewPoolClient(network, addr, pass string, timeout time.Duration, logger log.ILogger, pool *Pool) (*PoolClient, error) {
	var client *SimpleClient
	var err error

	client = pool.get()
	if client == nil {
		client, err = NewSimpleClient(network, addr, pass, timeout, new(log.NoopLogger))
		if err != nil {
			return nil, err
		}
	}

	return &PoolClient{
		client: client,
		pool:   pool,

		logger: logger,
	}, nil
}

func (this *PoolClient) Free() {
	this.pool.put(this.client)

	this.logger.Free()
}

func(this *PoolClient) LastCmd() []byte{
	return this.client.LastCmd()
}

func (this *PoolClient)Hset(key, field, value string) error {
	err:=this.client.Hset(key, field,value)

	this.logger.Info(this.client.LastCmd())

	return err
}

func (this *PoolClient)Hmset(key string, fieldValuePairs ...string) error {
	err:=this.client.Hmset(key string, fieldValuePairs ...)

	this.logger.Info(this.client.LastCmd())

	return err
}