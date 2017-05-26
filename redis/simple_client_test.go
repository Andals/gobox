package redis

import (
	"andals/gobox/log"
	"andals/gobox/log/writer"

	"testing"
	"time"
)

func TestSimpleClientHash(t *testing.T) {
	client := getSimpleClient()

	testHash(t, client)
}

func TestSimpleClientString(t *testing.T) {
	client := getSimpleClient()

	testString(t, client)
}

func TestSimpleClientExpire(t *testing.T) {
	client := getSimpleClient()

	testExpire(t, client)
}

func getSimpleClient() *simpleClient {
	w, _ := writer.NewFileWriter("/tmp/test_redis_simple_client.log")
	logger, _ := log.NewSimpleLogger(w, log.LEVEL_INFO, new(log.SimpleFormater))
	client, _ := NewSimpleClient("tcp", "127.0.0.1:6379", "123", time.Second*3, logger)

	return client
}
