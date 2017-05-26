package redis

func (this *simpleClient) Hset(key, field, value string) error {
	r := this.runCmd("HSET", key, field, value)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *simpleClient) Hmset(key string, fieldValuePairs ...string) error {
	r := this.runCmd("HMSET", append([]string{key}, fieldValuePairs...)...)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *simpleClient) Hget(key, field string) *StringResult {
	r := this.runCmd("HGET", key, field)

	return newStringResult(r)
}

func (this *simpleClient) Hgetall(key string) *HashResult {
	r := this.runCmd("HGETALL", key)

	return newHashResult(r)
}
