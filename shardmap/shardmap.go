//ref https://github.com/DeanThompson/syncmap/blob/master/syncmap.go

package shardmap

import (
	"sync"
)

const (
	DEF_SHARD_CNT = 32

	BKDR_SEED = 131 // 31 131 1313 13131 131313 etc...
)

type shardItem struct {
	sync.RWMutex

	data map[string]interface{}
}

type ShardMap struct {
	shardCnt uint8
	shards   []*shardItem
}

/**
* @param uint8, shardCnt must be pow of two
 */
func New(shardCnt uint8) *ShardMap {
	if !isPowOfTwo(shardCnt) {
		shardCnt = DEF_SHARD_CNT
	}

	this := &ShardMap{
		shardCnt: shardCnt,
		shards:   make([]*shardItem, shardCnt),
	}

	for i, _ := range this.shards {
		this.shards[i] = &shardItem{
			data: make(map[string]interface{}),
		}
	}

	return this
}

func (this *ShardMap) Get(key string) (interface{}, bool) {
	si := this.locate(key)

	si.RLock()
	value, ok := si.data[key]
	si.RUnlock()

	return value, ok
}

func (this *ShardMap) Set(key string, value interface{}) {
	si := this.locate(key)

	si.Lock()
	si.data[key] = value
	si.Unlock()
}

func (this *ShardMap) Del(key string) {
	si := this.locate(key)

	si.Lock()
	delete(si.data, key)
	si.Unlock()
}

type kvItem struct {
	key   string
	value interface{}
}

func (this *ShardMap) Walk(wf func(k string, v interface{})) {
	for _, si := range this.shards {
		kvCh := make(chan *kvItem)

		go func() {
			si.RLock()

			for k, v := range si.data {
				si.RUnlock()
				kvCh <- &kvItem{
					key:   k,
					value: v,
				}
				si.RLock()
			}

			si.RUnlock()
			close(kvCh)
		}()

		for {
			kv, ok := <-kvCh
			if !ok {
				break
			}
			wf(kv.key, kv.value)
		}
	}
}

func (this *ShardMap) locate(key string) *shardItem {
	i := bkdrHash(key) & uint32(this.shardCnt-1)

	return this.shards[i]
}

func isPowOfTwo(x uint8) bool {
	return x != 0 && (x&(x-1) == 0)
}

//https://www.byvoid.com/blog/string-hash-compare/
func bkdrHash(str string) uint32 {
	var h uint32

	for _, c := range str {
		h = h*BKDR_SEED + uint32(c)
	}

	return h
}
