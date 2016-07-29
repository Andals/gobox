package crypto

import (
	"crypto/md5"
	"fmt"
)

func Md5(data []byte) []byte {
	str := fmt.Sprintf("%x", md5.Sum(data))

	return []byte(str)
}
