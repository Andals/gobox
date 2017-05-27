/**
* @file async.go
* @brief write msg in independency goroutine
* @author ligang
* @date 2016-07-19
 */

package log

import "fmt"

var alr *asyncLogRoutine

// must be called first
func InitAsyncLogRoutine(msgQueueLen int, maxLoggerNum int) {
	alr = &asyncLogRoutine{
		msgCh:   make(chan *asyncMsg, msgQueueLen),
		flushCh: make(chan ILogger, maxLoggerNum),
		freeCh:  make(chan int),

		lockCh: make(chan int, 1),
		addCh:  make(chan *asyncLogger, maxLoggerNum),
		delCh:  make(chan *asyncLogger, maxLoggerNum),

		allist: make(map[string]*asyncLogger),
	}

	alr.lockCh <- 1

	go alr.run()
}

func FreeAsyncLogRoutine() {
	alr.freeCh <- 1
	<-alr.freeCh

	close(alr.msgCh)
	close(alr.flushCh)
	close(alr.freeCh)

	close(alr.lockCh)
	close(alr.addCh)
	close(alr.delCh)

	alr.allist = nil
}

type asyncMsg struct {
	level int
	msg   []byte

	al *asyncLogger
}

/**
* @name async log routine
* @{ */

type asyncLogRoutine struct {
	msgCh   chan *asyncMsg
	flushCh chan ILogger
	freeCh  chan int

	addCh  chan *asyncLogger
	delCh  chan *asyncLogger
	lockCh chan int

	allist map[string]*asyncLogger
}

func (this *asyncLogRoutine) run() {
	for {
		select {
		case logger, _ := <-this.addCh:
			key := asyncLoggerKey(logger)

			<-this.lockCh
			this.allist[key] = logger
			this.lockCh <- 1
		case logger, _ := <-this.delCh:
			key := asyncLoggerKey(logger)

			<-this.lockCh
			delete(this.allist, key)
			this.lockCh <- 1
		case am, _ := <-this.msgCh:
			logAsyncMsg(am)
		case logger, _ := <-this.flushCh:
			logger.Flush()
		case <-this.freeCh:
			for len(this.msgCh) != 0 {
				am, _ := <-this.msgCh
				logAsyncMsg(am)
			}
			for _, al := range this.allist {
				al.logger.Free()
			}
			this.freeCh <- 1
			return
		}
	}
}

func asyncLoggerKey(logger *asyncLogger) string {
	return fmt.Sprintf("%p", logger)
}

func logAsyncMsg(am *asyncMsg) {
	am.al.logger.Log(am.level, am.msg)
	am.al.msgCnt--

	if am.al.waitFree && am.al.msgCnt == 0 {
		am.al.logger.Free()
		alr.delCh <- am.al
	}
}

/**  @} */

/**
* @name async logger
* @{ */

type asyncLogger struct {
	logger ILogger

	msgCnt   int
	waitFree bool
}

func NewAsyncLogger(logger ILogger) *asyncLogger {
	this := &asyncLogger{
		logger: logger,

		msgCnt:   0,
		waitFree: false,
	}

	alr.addCh <- this

	return this
}

func (this *asyncLogger) Debug(msg []byte) {
	this.Log(LEVEL_DEBUG, msg)
}

func (this *asyncLogger) Info(msg []byte) {
	this.Log(LEVEL_INFO, msg)
}

func (this *asyncLogger) Notice(msg []byte) {
	this.Log(LEVEL_NOTICE, msg)
}

func (this *asyncLogger) Warning(msg []byte) {
	this.Log(LEVEL_WARNING, msg)
}

func (this *asyncLogger) Error(msg []byte) {
	this.Log(LEVEL_ERROR, msg)
}

func (this *asyncLogger) Critical(msg []byte) {
	this.Log(LEVEL_CRITICAL, msg)
}

func (this *asyncLogger) Alert(msg []byte) {
	this.Log(LEVEL_ALERT, msg)
}

func (this *asyncLogger) Emergency(msg []byte) {
	this.Log(LEVEL_EMERGENCY, msg)
}

func (this *asyncLogger) Log(level int, msg []byte) error {
	am := &asyncMsg{
		level: level,
		msg:   msg,

		al: this,
	}

	this.msgCnt++
	alr.msgCh <- am

	return nil
}

func (this *asyncLogger) Flush() error {
	alr.flushCh <- this.logger

	return nil
}

func (this *asyncLogger) Free() {
	if this.msgCnt == 0 {
		this.logger.Free()
		alr.delCh<-this
	} else {
		this.waitFree = true
	}
}

/**  @} */
