package redis

import (
	"andals/gobox/log"
	"andals/gobox/log/writer"

	"testing"
	"time"
)

func TestPool(t *testing.T) {
	w, _ := writer.NewFileWriter("/tmp/test_redis_pool.log")
	logger, _ := log.NewSimpleLogger(w, log.LEVEL_INFO, new(log.SimpleFormater))
	pool := NewPool(time.Second*3600, 300,
		func() (IClient, error) {
			return NewSimpleClient("tcp", "127.0.0.1:6379", "123", time.Second*3, logger)
		})

	client, _ := pool.Get()
	key := "ligang_pool:1"
	client.Set(key, "11")
	sr := client.Get(key)
	if sr.Value != "11" {
		t.Error("pool set get error")
	}
}
