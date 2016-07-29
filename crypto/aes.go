package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

type AesCBCCrypter struct {
	blockSize int

	encryptBlockMode cipher.BlockMode
	decryptBlockMode cipher.BlockMode

	padding PaddingInterface
}

func NewAesCBCCrypter(key []byte, iv []byte) (*AesCBCCrypter, error) {
	l := len(key)
	if l != 32 && l != 24 && l != 16 {
		return nil, errors.New("The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.")
	}

	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return nil, errors.New("The length of iv must be the same as the Block's block size " + strconv.Itoa(blockSize))
	}

	this := &AesCBCCrypter{
		blockSize: blockSize,

		encryptBlockMode: cipher.NewCBCEncrypter(block, iv),
		decryptBlockMode: cipher.NewCBCDecrypter(block, iv),

		padding: &PKCS5Padding{
			BlockSize: blockSize,
		},
	}

	return this, nil
}

func (this *AesCBCCrypter) BlockSize() int {
	return this.blockSize
}

func (this *AesCBCCrypter) SetPadding(padding PaddingInterface) {
	this.padding = padding
}

func (this *AesCBCCrypter) Encrypt(data []byte) []byte {
	data = this.padding.Padding(data)

	crypted := make([]byte, len(data))
	this.encryptBlockMode.CryptBlocks(crypted, data)

	return crypted
}

func (this *AesCBCCrypter) Decrypt(crypted []byte) []byte {
	data := make([]byte, len(crypted))
	this.decryptBlockMode.CryptBlocks(data, crypted)

	return this.padding.UnPadding(data)
}
