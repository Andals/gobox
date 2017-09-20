package redis

import (
	"github.com/fzzy/radix/redis"

	"andals/gobox/log"

	"time"
)

type StringResult struct {
	Value string
	Err   error
}

type HashResult struct {
	Value map[string]string
	Err   error
}

type IClient interface {
	SetLogger(logger log.ILogger)
	Free()

	Expire(key, seconds string) error
	Del(key string) error

	//args：[EX seconds] [PX milliseconds] [NX|XX]
	Set(key, value string, args ...string) error
	Setex(key, seconds, value string) error
	Get(key string) *StringResult

	Hset(key, field, value string) error
	Hmset(key string, fieldValuePairs ...string) error
	Hget(key, field string) *StringResult
	Hgetall(key string) *HashResult
	Hdel(key string, fields ...string) error
	RunCmd(key string, fields ...string) *redis.Reply
}

func newStringResult(r *redis.Reply) *StringResult {
	if r.Type == redis.NilReply {
		return nil
	}

	result := new(StringResult)
	if r.Err != nil {
		result.Err = r.Err
		return result
	}

	var err error
	result.Value, err = r.Str()
	if err != nil {
		result.Err = err
	}

	return result
}

func newHashResult(r *redis.Reply) *HashResult {
	if r.Type == redis.NilReply {
		return nil
	}

	result := new(HashResult)
	if r.Err != nil {
		result.Err = r.Err
		return result
	}

	var err error
	result.Value, err = r.Hash()
	if err != nil {
		result.Err = err
	}

	return result
}

func newRedisClient(network, addr, pass string, timeout time.Duration) (*redis.Client, error) {
	client, e := redis.DialTimeout(network, addr, timeout)
	if e != nil {
		return nil, e
	}

	r := client.Cmd("AUTH", pass)
	if r.Err != nil {
		client.Close()

		return nil, r.Err
	}

	return client, nil
}
