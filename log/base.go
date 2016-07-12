/**
* @file logger.go
* @author ligang
* @date 2016-02-04
 */

package log

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
