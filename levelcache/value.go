package levelcache

import (
	"bytes"
	"encoding/binary"
)

type CacheBin struct {
	AddTime int64
	Expire  int64
	Size    int64
}

type CacheValue struct {
	*CacheBin

	Value []byte
}

func parseByBinary(bs []byte) (*CacheValue, error) {
	buf := bytes.NewBuffer(bs)

	cv := &CacheValue{}
	cv.CacheBin = &CacheBin{}

	err := binary.Read(buf, binary.LittleEndian, cv.CacheBin)
	if err != nil {
		return nil, err
	}

	cv.Value = make([]byte, cv.Size)
	err = binary.Read(buf, binary.LittleEndian, cv.Value)
	if err != nil {
		return nil, err
	}

	return cv, nil
}

func (this *CacheValue) toBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, this.CacheBin)
	if err != nil {
		return nil, err
	}

	buf.Write(this.Value)

	return buf.Bytes(), nil
}
