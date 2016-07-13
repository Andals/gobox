package log

import (
	"testing"

	"andals/gobox/misc"
)

func TestSimpleFormater(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	f := new(SimpleFormater)

	b := f.Format(LEVEL_EMERGENCY, []byte("abc"))
	t.Log(string(b))
}

func TestWebFormater(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	f := NewWebFormater("xyz")

	b := f.Format(LEVEL_EMERGENCY, []byte("abc"))
	t.Log(string(b))
}
