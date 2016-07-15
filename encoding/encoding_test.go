package encoding

import (
	"andals/gobox/misc"

	"testing"
)

func TestBase64EncodeDecode(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	od := []byte("abc")

	cd := Base64Encode(od)
	t.Log(string(cd))

	dd := Base64Decode(cd)
	t.Log(string(dd))

	if string(dd) != string(od) {
		t.Error("coding error")
	}
}
