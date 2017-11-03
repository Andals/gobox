package crypto

import (
	"testing"
)

func TestMd5(t *testing.T) {
	md5 := Md5([]byte("abc"))
	if len(md5) != 32 {
		t.Error(string(md5))
	}

	t.Log(string(md5))
}
