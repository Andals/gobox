package crypto

import (
	"andals/gobox/misc"

	"testing"
)

func TestMd5(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	md5 := Md5([]byte("abc"))
	if len(md5) != 32 {
		t.Error(string(md5))
	}

	t.Log(string(md5))
}

func TestPKCS5Padding(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	data := []byte("abcd")
	bs := 16
	pd := PKCS5Padding(data, bs)
	t.Log(data, pd)

	if len(pd)%bs != 0 {
		t.Error(pd)
	}
}

func TestPKCS5UnPadding(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	data := []byte("abcd")
	bs := 16
	pd := PKCS5Padding(data, bs)
	t.Log(data, pd)

	ud := PKCS5UnPadding(pd)
	t.Log(ud)

	if string(data) != string(ud) {
		t.Error(ud)
	}
}

func TestAesEncryptDecrypt(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	key := []byte("gobox")
	data := []byte("abc")

	crypted := AesEncrypt(key, data)
	t.Log(crypted)

	d := AesDecrypt(key, crypted)
	t.Log(d)

	if string(d) != string(data) {
		t.Error(d, data)
	}
}
