package redis

import (
	"testing"
	"time"
)

func testHash(t *testing.T, client IClient) {
	key := "ligang_hash:1"
	client.Hset(key, "2017-05-01", "10")
	client.Hmset(key, "2017-05-02", "3", "2017-05-03", "4", "2017-05-04", "5")

	sr := client.Hget(key, "2017-05-01")
	if sr.Value != "10" {
		t.Error("str must be equal 10")
	}

	client.Hdel(key, "2017-05-01")
	sr = client.Hget(key, "2017-05-01")
	if sr != nil {
		t.Error("HDEL error")
	}

	om := map[string]string{
		"2017-05-01": "10",
		"2017-05-02": "3",
		"2017-05-03": "4",
		"2017-05-04": "5",
	}
	mr := client.Hgetall(key)
	for k, v := range mr.Value {
		ov, ok := om[k]
		if !ok {
			t.Error("key: " + k + " not exists")
		}
		if ov != v {
			t.Error("key: " + k + " value: " + v + " not equal ov: " + ov)
		}
	}
}

func testString(t *testing.T, client IClient) {
	key := "ligang_string:1"
	client.Set(key, "10")
	sr := client.Get(key)
	if sr.Value != "10" {
		t.Error("string set get error")
	}

	client.Setex(key, "2", "10")
	time.Sleep(time.Second * 3)
	sr = client.Get(key)
	if sr != nil {
		t.Error("string set ex error", sr)
	}
}

func testExpire(t *testing.T, client IClient) {
	key := "ligang_string:1"
	client.Set(key, "10")
	client.Expire(key, "2")
	time.Sleep(time.Second * 3)
	sr := client.Get(key)
	if sr != nil {
		t.Error("expire error", sr)
	}
}

func testDel(t *testing.T, client IClient) {
	key := "ligang_string:1"
	client.Set(key, "10")
	client.Del(key)
	sr := client.Get(key)
	if sr != nil {
		t.Error("hdel error", sr)
	}
}
