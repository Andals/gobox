package pool

import (
	"github.com/andals/gobox/redis"

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

func newRedisClient() (IConn, error) {
	config := &redis.Config{
		Host: "127.0.0.1",
		Port: "6379",
		Pass: "123",
	}

	return redis.NewClient(config, nil)
}

func testPool(pool *Pool, t *testing.T) {
	conn, _ := pool.Get()
	t.Log("conn", conn)

	client := conn.(*redis.Client)
	client.Do("set", "pool", "pool_test")
	reply, _ := client.Do("get", "pool")
	t.Log(reply.String())

	pool.Put(client)
}
