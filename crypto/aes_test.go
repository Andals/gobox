package crypto

import (
	"andals/gobox/misc"

	"crypto/aes"
	"testing"
)

func TestAesCBCCrypter(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	key := Md5([]byte("gobox"))
	iv := Md5([]byte("andals"))[:aes.BlockSize]
	data := []byte("abc")

	acc, err := NewAesCBCCrypter(key, iv)
	t.Log(err)
	t.Log(acc.BlockSize())

	crypted := acc.Encrypt(data)
	t.Log(crypted)

	d := acc.Decrypt(crypted)
	t.Log(d)

	if string(d) != string(data) {
		t.Error(d, data)
	}
}
