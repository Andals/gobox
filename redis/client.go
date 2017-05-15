package Redis

import (
	//"fmt"
	"time"

	"github.com/fzzy/radix/redis"
)

type Client struct {
	client *redis.Client
}

func NewClient(network, addr, pass string, timeout time.Duration) (*Client, error) {
	c, e := redis.DialTimeout(network, addr, timeout)
	if e != nil {
		return nil, e
	}

	this := &Client{
		client: c,
	}
	r := this.client.Cmd("AUTH", pass)
	if r.Err != nil {
		this.client.Close()

		return nil, r.Err
	}

	return this, nil
}

func (this *Client) Close() {
	this.client.Close()
}
