package formater

import (
	"testing"

	"andals/gobox/log/level"
	"andals/gobox/misc"
)

func TestSimple(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	f := new(Simple)

	b := f.Format(level.LEVEL_EMERGENCY, []byte("abc"))
	t.Log(string(b))
}
