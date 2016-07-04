package shardmap

import (
	"andals/gobox/misc"

	//     "fmt"
	"strconv"
	"sync"
	"testing"
)

var smap *ShardMap

func init() {
	smap = New(32)
}

func TestSetGet(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	for i := 0; i < 10000; i++ {
		key := misc.Md5([]byte(strconv.Itoa(i)))
		smap.Set(key, i)

		v, ok := smap.Get(key)
		if !ok || v != i {
			t.Error(v, ok)
		}
	}
}

func TestWalkDel(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	smap.Walk(func(k string, v interface{}) {
		t.Log(k, v)

		smap.Del(k)

		_, ok := smap.Get(k)
		if ok {
			t.Error(v, ok)
		}
	})
}

func BenchmarkRW(b *testing.B) {
	wg := new(sync.WaitGroup)

	for i := 0; i < 10000; i++ {
		key := misc.Md5([]byte(strconv.Itoa(i)))
		wg.Add(1)
		go write(key, i, wg)
	}
	wg.Wait()

	for i := 0; i < 10000; i++ {
		key := misc.Md5([]byte(strconv.Itoa(i)))
		wg.Add(1)
		go read(key, wg)
	}
	wg.Wait()
}

func write(k string, v interface{}, wg *sync.WaitGroup) {
	//     fmt.Println(k, v)
	smap.Set(k, v)
	wg.Done()
}

func read(k string, wg *sync.WaitGroup) {
	//     v, ok := smap.Get(k)
	smap.Get(k)
	//     fmt.Println(k, v, ok)
	wg.Done()
}
