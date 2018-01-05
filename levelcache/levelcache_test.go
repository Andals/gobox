package levelcache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache, _ := NewCache("/tmp/levelcache_test", 5*time.Second)

	key := []byte("k1")
	value := []byte("v1")

	cache.Set(key, value, 3)

	value, _ = cache.Get(key)
	sv := string(value)
	t.Log(sv)
	if sv != "v1" {
		t.Fatal("set get error")
	}

	time.Sleep(6 * time.Second)

	v, err := cache.Get(key)
	sv = string(v)
	t.Log(sv, err)
}
