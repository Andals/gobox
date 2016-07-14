package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
)

func Md5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func AesEncrypt(key, data []byte) []byte {
	key = []byte(Md5(key))
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)

	return crypted
}

func AesDecrypt(key, crypted []byte) []byte {
	key = []byte(Md5(key))
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()

	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	data := make([]byte, len(crypted))
	blockMode.CryptBlocks(data, crypted)

	return PKCS5UnPadding(data)
}

func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padtext...)
}

func PKCS5UnPadding(data []byte) []byte {
	l := len(data)
	unpadding := int(data[l-1])

	return data[:(l - unpadding)]
}
