package misc

type BoolSet map[string]bool

func (this BoolSet) IsTrue(key string) bool {
	value, ok := this[key]
	if ok && value {
		return true
	}

	return false
}
