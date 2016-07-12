package formater

type Noop struct {
}

func (this *Noop) Format(level int, msg []byte) []byte {
	return msg
}
