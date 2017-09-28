package golog

import "time"

type webFormater struct {
	logId []byte
	ip    []byte
}

func NewWebFormater(logId, ip []byte) *webFormater {
	return &webFormater{
		logId: logId[:],
		ip:    ip[:],
	}
}

func (this *webFormater) Format(level int, msg []byte) []byte {
	lm, ok := logLevels[level]
	if !ok {
		lm = "-"
	}

	return AppendBytes(
		[]byte("["),
		[]byte(lm),
		[]byte("]\t"),
		[]byte("["),
		[]byte(time.Now().Format(TimeGeneralLayout())),
		[]byte("]\t"),
		this.ip,
		[]byte("\t"),
		this.logId,
		[]byte("\t"),
		msg,
		[]byte("\n"),
	)
}
