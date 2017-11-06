package redis

import (
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := NewPool(time.Second*5, 300, newRedisClient)

	testPool(pool, t)
	testPool(pool, t)

	time.Sleep(time.Second * 7)
	testPool(pool, t)
}

func newRedisClient() (*Client, error) {
	config := &Config{
		Host: "127.0.0.1",
		Port: "6379",
		Pass: "123",
	}

	client := NewClient(config, nil)
	return client, nil
}

func testPool(pool *Pool, t *testing.T) {
	client, _ := pool.Get()
	client.Do("set", "redis_pool", "pool_test")
	reply, _ := client.Do("get", "redis_pool")
	t.Log(reply.String())

	pool.Put(client)
}
