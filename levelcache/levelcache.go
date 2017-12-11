package levelcache

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"

	"time"
)

type Cache struct {
	db *leveldb.DB

	jstopCh chan bool
}

func NewCache(cacheDir string, janitorInterval time.Duration) (*Cache, error) {
	db, err := leveldb.OpenFile(cacheDir, nil)
	if err != nil {
		return nil, err
	}

	this := &Cache{
		db: db,

		jstopCh: make(chan bool),
	}

	go this.runJanitor(janitorInterval)

	return this, nil
}

func (this *Cache) Free() {
	this.db.Close()
	close(this.jstopCh)

	this.jstopCh <- true
}

func (this *Cache) runJanitor(jinterval time.Duration) {
	ticker := time.NewTicker(jinterval)

	for {
		select {
		case <-this.jstopCh:
			return
		case <-ticker.C:
			now := time.Now().Unix()
			iter := this.db.NewIterator(nil, nil)
			for iter.Next() {
				key := iter.Key()
				cv, err := parseByBinary(iter.Value())
				if err != nil {
					this.Delete(key)
					continue
				}
				if cv.Expire == 0 {
					continue
				}
				if now-cv.AddTime > cv.Expire {
					this.Delete(key)
				}
			}
			iter.Release()
		}
	}
}

func (this *Cache) Set(key, value []byte, expireSeconds int64) error {
	cb := &CacheBin{
		AddTime: time.Now().Unix(),
		Expire:  expireSeconds,
		Size:    int64(len(value)),
	}
	cv := &CacheValue{
		CacheBin: cb,

		Value: value,
	}

	bv, err := cv.toBinary()
	if err != nil {
		return err
	}

	return this.db.Put(key, bv, nil)
}

func (this *Cache) Get(key []byte) ([]byte, error) {
	bv, err := this.db.Get(key, nil)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	cv, err := parseByBinary(bv)
	if err != nil {
		return nil, err
	}

	return cv.Value, nil
}

func (this *Cache) Delete(key []byte) error {
	return this.db.Delete(key, nil)
}
