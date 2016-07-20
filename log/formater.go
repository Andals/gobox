/**
* @file base.go
* @brief format msg before send to writer
* @author ligang
* @date 2016-07-12
 */

package log

import (
	"time"

	"andals/gobox/misc"
)

type IFormater interface {
	Format(level int, msg []byte) []byte
}

/**
* @name NoopFormater
* @{ */

type NoopFormater struct {
}

func (this *NoopFormater) Format(level int, msg []byte) []byte {
	return msg
}

/**  @} */

/**
* @name SimpleFormater which add [logLevelMsg] ahead of msg
* @{ */

type SimpleFormater struct {
}

func (this *SimpleFormater) Format(level int, msg []byte) []byte {
	lm, ok := logLevels[level]
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

/**  @} */

/**
* @name WebFormater which add [logLevelMmsg], [time] and [logId] ahead of msg
* @{ */

type WebFormater struct {
	logId []byte
}

func NewWebFormater(logId []byte) *WebFormater {
	this := new(WebFormater)
	copy(this.logId, logId)

	return this
}

func (this *WebFormater) Format(level int, msg []byte) []byte {
	lm, ok := logLevels[level]
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

/**  @} */
