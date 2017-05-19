package redis

import (
	"github.com/fzzy/radix/redis"

	"time"
)

type IClient interface {
	LastCmd() []byte
	Free()

	Hset(key, field, value string) error
	Hmset(key string, fieldValuePairs ...string) error
	Hget(key, field string) (string, error)
	Hgetall(key string) (map[string]string, error)
}

func newRedisClient(network, addr, pass string, timeout time.Duration) (*redis.Client, error) {
	client, e := redis.DialTimeout(network, addr, timeout)
	if e != nil {
		return nil, e
	}

	r := client.Cmd("AUTH", pass)
	if r.Err != nil {
		client.Close()

		return nil, r.Err
	}

	return client, nil
}
