package formater

import (
	logLevel "andals/gobox/log/level"
	"andals/gobox/misc"
)

type Simple struct {
}

func (this *Simple) Format(level int, msg []byte) []byte {
	lm, ok := logLevel.LogLevels[level]
	if !ok {
		lm = "-"
	}
	return misc.AppendBytes(
		[]byte("["),
		[]byte(lm),
		[]byte("]\t"),
		msg,
		[]byte("\n"),
	)
}
