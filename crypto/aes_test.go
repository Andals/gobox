package crypto

import (
	"andals/gobox/misc"

	"testing"
)

func TestAesCBCCrypter(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	key := Md5([]byte("gobox"))
	iv := Md5([]byte("andals"))[:AES_BLOCK_SIZE]
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
