package formater

import (
	"time"

	logLevel "andals/gobox/log/level"
	"andals/gobox/misc"
)

type Web struct {
	logId []byte
}

func New(logId string) *Web {
	return &Web{
		logId: []byte(logId),
	}
}

func (this *Web) Format(level int, msg []byte) []byte {
	lm, ok := logLevel.LogLevels[level]
	if !ok {
		lm = "-"
	}

	return misc.AppendBytes(
		[]byte("["),
		[]byte(lm),
		[]byte("]\t"),
		[]byte("["),
		[]byte(time.Now().Format(misc.TimeGeneralLayout())),
		[]byte("]\t"),
		this.logId,
		[]byte("\t"),
		msg,
		[]byte("\n"),
	)
}
