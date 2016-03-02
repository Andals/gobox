/**
* @file async.go
* @brief async log use goroutine, one logger one routine
* @author ligang
* @date 2016-02-06
 */

package log

type asyncMsg struct {
	level int
	msg   []byte
}

type asyncLogger struct {
	msgCh   chan *asyncMsg
	flushCh chan int
	freeCh  chan int

	ILogger
}

type asyncLoggerList struct {
	// prevent asyncLoggerList race
	lockCh chan int

	loggers []*asyncLogger
}

var allist asyncLoggerList

func init() {
	allist.lockCh = make(chan int, 1)
	allist.lockCh <- 1
}

func NewAsyncLogger(logger ILogger, queueLen int) (*asyncLogger, error) {
	defer func() {
		allist.lockCh <- 1
	}()

	<-allist.lockCh

	this := &asyncLogger{
		msgCh:   make(chan *asyncMsg, queueLen),
		flushCh: make(chan int),
		freeCh:  make(chan int),

		ILogger: logger,
	}

	go this.logRoutine()

	allist.loggers = append(allist.loggers, this)

	return this, nil
}

func FreeAllAsyncLogger() {
	for _, logger := range allist.loggers {
		logger.Free()
	}
}

func (this *asyncLogger) Log(level int, msg []byte) error {
	am := &asyncMsg{
		level: level,
		msg:   msg,
	}

	this.msgCh <- am

	return nil
}

func (this *asyncLogger) Flush() error {
	this.flushCh <- 1

	return nil
}

func (this *asyncLogger) Free() {
	defer func() {
		<-this.freeCh
	}()

	this.freeCh <- 1
}

func (this *asyncLogger) logRoutine() {
	defer func() {
		this.freeCh <- 1
	}()

	for {
		select {
		case am, _ := <-this.msgCh:
			this.ILogger.Log(am.level, am.msg)
		case <-this.flushCh:
			this.ILogger.Flush()
		case <-this.freeCh:
			for 0 != len(this.msgCh) {
				am, _ := <-this.msgCh
				this.ILogger.Log(am.level, am.msg)
			}
			this.ILogger.Free()
			return
		}
	}
}
