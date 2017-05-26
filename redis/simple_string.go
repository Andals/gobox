package redis

func (this *simpleClient) Set(key, value string, args ...string) error {
	r := this.runCmd("SET", append([]string{key, value}, args...)...)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *simpleClient) Setex(key, seconds, value string) error {
	r := this.runCmd("SETEX", key, seconds, value)
	if r.Err != nil {
		return r.Err
	}

	return nil
}

func (this *simpleClient) Get(key string) *StringResult {
	r := this.runCmd("GET", key)

	return newStringResult(r)
}
