package redis

func (this *SimpleClient) Hset(key, field, value string) error {
	r := this.runCmd("HSET", key, field, value)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *SimpleClient) Hmset(key string, fieldValuePairs ...string) error {
	r := this.runCmd("HMSET", append([]string{key}, fieldValuePairs...)...)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *SimpleClient) Hget(key, field string) (string, error) {
	r := this.runCmd("HGET", key, field)
	if r.Err != nil {
		return "", r.Err
	}

	str, err := r.Str()
	if err != nil {
		return "", err
	}
	return str, nil
}

func (this *SimpleClient) Hgetall(key string) (map[string]string, error) {
	r := this.runCmd("HGETALL", key)
	if r.Err != nil {
		return nil, r.Err
	}

	hash, err := r.Hash()
	if err != nil {
		return nil, err
	}
	return hash, nil
}
