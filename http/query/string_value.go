package query

type CheckString func(v string) bool

func CheckStringNotEmpty(v string) bool {
	if v == "" {
		return false
	}
	return true
}

type stringValue struct {
	*baseValue

	p  *string
	cf CheckString
}

func NewStringValue(p *string, required bool, errno int, msg string, cf CheckString) *stringValue {
	this := &stringValue{
		baseValue: newBaseValue(required, errno, msg),

		p:  p,
		cf: cf,
	}

	return this
}

func (this *stringValue) Set(str string) error {
	*(this.p) = str

	return nil
}

func (this *stringValue) Check() bool {
	if this.cf == nil {
		return true
	}

	return this.cf(*(this.p))
}
