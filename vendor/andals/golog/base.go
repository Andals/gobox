/**
* @file logger.go
* @author ligang
* @date 2016-02-04
 */

package golog

import "io"

const (
	LEVEL_DEBUG     = 1
	LEVEL_INFO      = 2
	LEVEL_NOTICE    = 3
	LEVEL_WARNING   = 4
	LEVEL_ERROR     = 5
	LEVEL_CRITICAL  = 6
	LEVEL_ALERT     = 7
	LEVEL_EMERGENCY = 8
)

var logLevels map[int]string = map[int]string{
	LEVEL_DEBUG:     "debug",
	LEVEL_INFO:      "info",
	LEVEL_NOTICE:    "notice",
	LEVEL_WARNING:   "warning",
	LEVEL_ERROR:     "error",
	LEVEL_CRITICAL:  "critical",
	LEVEL_ALERT:     "alert",
	LEVEL_EMERGENCY: "emergency",
}

type ILogger interface {
	Debug(msg []byte)
	Info(msg []byte)
	Notice(msg []byte)
	Warning(msg []byte)
	Error(msg []byte)
	Critical(msg []byte)
	Alert(msg []byte)
	Emergency(msg []byte)

	Log(level int, msg []byte) error

	Flush() error
	Free()
}

type IFormater interface {
	Format(level int, msg []byte) []byte
}

type IWriter interface {
	io.Writer

	Flush() error
	Free()
}
