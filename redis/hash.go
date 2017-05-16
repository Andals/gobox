package Redis

func (this *Client) Hset(key, field, value string) error {
	r := this.runCmd("HSET", key, field, value)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *Client) Hmset(key string, fieldValues ...string) error {
	r := this.runCmd("HMSET", append([]string{key}, fieldValues...)...)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *Client) Hget(key, field string) (string, error) {
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

func (this *Client) Hgetall(key string) (map[string]string, error) {
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
