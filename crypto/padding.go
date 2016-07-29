package crypto

import (
	"bytes"
)

type PaddingInterface interface {
	Padding(data []byte) []byte
	UnPadding(data []byte) []byte
}

type PKCS5Padding struct {
	BlockSize int
}

func (this *PKCS5Padding) Padding(data []byte) []byte {
	padding := this.BlockSize - len(data)%this.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padtext...)
}

func (this *PKCS5Padding) UnPadding(data []byte) []byte {
	l := len(data)
	unpadding := int(data[l-1])

	return data[:(l - unpadding)]
}
