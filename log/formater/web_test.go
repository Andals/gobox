package formater

import (
	"testing"

	"andals/gobox/log/level"
	"andals/gobox/misc"
)

func TestWeb(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	f := NewWeb("xyz")

	b := f.Format(level.LEVEL_EMERGENCY, []byte("abc"))
	t.Log(string(b))
}
