package simplecache

import (
	"github.com/andals/gobox/crypto"
	"github.com/andals/gobox/misc"

	//     "fmt"
	"strconv"
	"testing"
	"time"
)

var sc *SimpleCache

func init() {
	sc = New(32)
}

func TestSetGet(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	for i := 0; i < 10000; i++ {
		key := string(crypto.Md5([]byte(strconv.Itoa(i))))
		sc.Set(key, i, 10*time.Second)

		v, ok := sc.Get(key)
		if !ok || v != i {
			t.Error(v, ok)
		}
	}

	time.Sleep(16 * time.Second)

	for i := 0; i < 10000; i++ {
		key := string(crypto.Md5([]byte(strconv.Itoa(i))))

		v, ok := sc.Get(key)
		if ok || v == i {
			t.Error(v, ok)
		}
	}
}
