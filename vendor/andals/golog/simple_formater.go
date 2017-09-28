/**
* @file base.go
* @brief format msg before send to writer
* @author ligang
* @date 2016-07-12
 */

package golog

import (
	"time"
)

type simpleFormater struct {
}

func NewSimpleFormater() *simpleFormater {
	return new(simpleFormater)
}

func (this *simpleFormater) Format(level int, msg []byte) []byte {
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
		msg,
		[]byte("\n"),
	)
}
