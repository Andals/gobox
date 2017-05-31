package redis

import (
	"errors"
	"time"
)

type poolItem struct {
	client IClient

	addTime time.Time
}

type Pool struct {
	clientTimeout time.Duration

	ch chan *poolItem

	newClientFunc func() (IClient, error)
}

func NewPool(clientTimeout time.Duration, size int, newClientFunc func() (IClient, error)) *Pool {
	this := &Pool{
		clientTimeout: clientTimeout,

		ch: make(chan *poolItem, size),

		newClientFunc: newClientFunc,
	}

	return this
}

func (this *Pool) Get() (IClient, error) {
	select {
	case pItem := <-this.ch:
		if time.Now().Unix()-pItem.addTime.Unix() < int64(this.clientTimeout) {
			return pItem.client, nil
		}
	default:
	}

	return this.newClientFunc()
}

func (this *Pool) Put(client IClient) error {
	pItem := &poolItem{
		client: client,

		addTime: time.Now(),
	}

	select {
	case this.ch <- pItem:
		return nil
	default:
		client.Free()
	}

	return errors.New("pool is full")
}
