package simplecache

import (
	"github.com/andals/gobox/shardmap"

	"time"
)

const (
	TICK_INTERVAL = 30 * time.Second

	NO_EXPIRE = 0
)

type cacheItem struct {
	value  interface{}
	expire int64

	setTime    time.Time
	accessTime time.Time
}

type SimpleCache struct {
	shardMap *shardmap.ShardMap
}

func New(shardCnt uint8) *SimpleCache {
	this := &SimpleCache{
		shardMap: shardmap.New(shardCnt),
	}

	go this.runJanitor()

	return this
}

func (this *SimpleCache) Set(key string, value interface{}, expire time.Duration) {
	this.shardMap.Set(key, newCacheItem(value, expire))
}

func (this *SimpleCache) Get(key string) (interface{}, bool) {
	value, ok := this.shardMap.Get(key)
	if !ok {
		return nil, false
	}

	ci, ok := value.(*cacheItem)

	now := time.Now()
	if expired(now.UnixNano(), ci) {
		return nil, false
	}

	ci.accessTime = now

	return ci.value, true
}

func (this *SimpleCache) runJanitor() {
	ticker := time.NewTicker(TICK_INTERVAL)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			this.shardMap.Walk(func(k string, v interface{}) {
				ci, ok := v.(*cacheItem)
				if ok {
					if expired(now, ci) {
						this.shardMap.Del(k)
					}
				}
			})
		}
	}
}

func expired(now int64, ci *cacheItem) bool {
	if ci.expire == NO_EXPIRE {
		return false
	}

	if (now - ci.expire) >= ci.setTime.UnixNano() {
		return true
	}

	return false
}

func newCacheItem(value interface{}, expire time.Duration) *cacheItem {
	now := time.Now()

	return &cacheItem{
		value:  value,
		expire: int64(expire),

		setTime:    now,
		accessTime: now,
	}
}
