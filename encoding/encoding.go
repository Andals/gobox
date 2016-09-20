package encoding

import (
	"encoding/base64"
)

func Base64Encode(data []byte) []byte {
	coded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(coded, data)

	return coded
}

func Base64Decode(coded []byte) []byte {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(coded)))
	n, _ := base64.StdEncoding.Decode(data, coded)

	return data[:n]
}
