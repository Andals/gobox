package golog

import (
	"testing"

	"andals/gobox/misc"
)

func TestSimpleFormater(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	f := NewSimpleFormater()

	b := f.Format(LEVEL_EMERGENCY, []byte("abc"))
	t.Log(string(b))
}

func TestWebFormater(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	f := NewWebFormater([]byte("xyz"), []byte("10.0.0.1"))

	b := f.Format(LEVEL_EMERGENCY, []byte("abc"))
	t.Log(string(b))
}
