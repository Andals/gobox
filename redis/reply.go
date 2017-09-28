package redis

import "github.com/garyburd/redigo/redis"

type Reply struct {
	reply interface{}
}

func (this *Reply) Bool() (bool, error) {
	return redis.Bool(this.reply, nil)
}

func (this *Reply) ByteSlices() ([][]byte, error) {
	return redis.ByteSlices(this.reply, nil)
}

func (this *Reply) Bytes() ([]byte, error) {
	return redis.Bytes(this.reply, nil)
}

func (this *Reply) Float64() (float64, error) {
	return redis.Float64(this.reply, nil)
}

func (this *Reply) Int() (int, error) {
	return redis.Int(this.reply, nil)
}

func (this *Reply) Int64() (int64, error) {
	return redis.Int64(this.reply, nil)
}

func (this *Reply) Int64Map() (map[string]int64, error) {
	return redis.Int64Map(this.reply, nil)
}

func (this *Reply) Ints() ([]int, error) {
	return redis.Ints(this.reply, nil)
}

func (this *Reply) Struct(s interface{}) error {
	return redis.ScanStruct(this.reply.([]interface{}), s)
}

func (this *Reply) String() (string, error) {
	return redis.String(this.reply, nil)
}

func (this *Reply) StringMap() (map[string]string, error) {
	return redis.StringMap(this.reply, nil)
}

func (this *Reply) Strings() ([]string, error) {
	return redis.Strings(this.reply, nil)
}

func (this *Reply) Uint64() (uint64, error) {
	return redis.Uint64(this.reply, nil)
}
