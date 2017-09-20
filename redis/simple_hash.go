package redis

func (this *simpleClient) Hset(key, field, value string) error {
	r := this.RunCmd("HSET", key, field, value)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *simpleClient) Hmset(key string, fieldValuePairs ...string) error {
	r := this.RunCmd("HMSET", append([]string{key}, fieldValuePairs...)...)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *simpleClient) Hget(key, field string) *StringResult {
	r := this.RunCmd("HGET", key, field)

	return newStringResult(r)
}

func (this *simpleClient) Hgetall(key string) *HashResult {
	r := this.RunCmd("HGETALL", key)

	return newHashResult(r)
}

func (this *simpleClient) Hdel(key string, fields ...string) error {
	r := this.RunCmd("HDEL", append([]string{key}, fields...)...)
	if r.Err != nil {
		return r.Err
	}

	return nil
}
