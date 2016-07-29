package crypto

import (
	"andals/gobox/misc"

	"testing"
)

func TestPKCS5Padding(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	padding := &PKCS5Padding{
		BlockSize: 16,
	}

	data := []byte("abcd")
	pd := padding.Padding(data)
	t.Log(data, pd)

	upd := padding.UnPadding(pd)
	t.Log(upd)

	if string(data) != string(upd) {
		t.Error(upd)
	}
}
