package redis

import (
	"andals/golog"

	"testing"
)

func TestClient(t *testing.T) {
	w, _ := golog.NewFileWriter("/tmp/test_redis.log")
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, golog.NewSimpleFormater())

	config := &Config{
		Host: "127.0.0.1",
		Port: "6379",
		Pass: "123",
	}
	client, _ := NewClient(config, logger)

	reply, _ := client.Do("set", "a", "1")
	t.Log(reply.String())
	reply, _ = client.Do("get", "a")
	t.Log(reply.Int())

	client.Send("set", "a", "a")
	client.Send("set", "b", "b")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, _ := client.ExecPipelining()
	for _, reply := range replies {
		t.Log(reply.String())
	}

	client.BeginTrans()
	client.Send("set", "a", "1")
	client.Send("set", "b", "2")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, _ = client.ExecTrans()
	for _, reply := range replies {
		t.Log(reply.String())
	}

	client.Free()
}
