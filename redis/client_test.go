package Redis

import (
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	client, _ := NewClient("tcp", "127.0.0.1:6379", "123", time.Duration(10)*time.Second)

	defer func() {
		client.Close()
	}()

	key := "ligang_hash:1"
	client.Hset(key, "2017-05-01", "10")
	client.Hmset(key, "2017-05-02", "3", "2017-05-03", "4", "2017-05-04", "5")

	str, _ := client.Hget(key, "2017-05-01")
	if str != "10" {
		t.Error("str must be equal 10")
	}

	om := map[string]string{
		"2017-05-01": "10",
		"2017-05-02": "3",
		"2017-05-03": "4",
		"2017-05-04": "5",
	}
	m, _ := client.Hgetall(key)
	for k, v := range m {
		ov, ok := om[k]
		if !ok {
			t.Error("key: " + k + " not exists")
		}
		if ov != v {
			t.Error("key: " + k + " value: " + v + " not equal ov: " + ov)
		}
	}
}
