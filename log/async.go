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

		addCh: make(chan *asyncLoggerAddChItem, maxLoggerNum),
		delCh: make(chan string, maxLoggerNum),

		asyncLoggers: make(map[string]*asyncLogger),
	}

	go alr.run()
}

func FreeAsyncLogRoutine() {
	alr.freeCh <- 1
	<-alr.freeCh

	close(alr.msgCh)
	close(alr.flushCh)
	close(alr.freeCh)

	close(alr.addCh)
	close(alr.delCh)

	alr.asyncLoggers = nil
}

type asyncMsg struct {
	level int
	msg   []byte

	al *asyncLogger
}

/**
* @name async log routine
* @{ */

type asyncLoggerAddChItem struct {
	key    string
	logger *asyncLogger
}

type asyncLogRoutine struct {
	msgCh   chan *asyncMsg
	flushCh chan ILogger
	freeCh  chan int

	addCh chan *asyncLoggerAddChItem
	delCh chan string

	asyncLoggers map[string]*asyncLogger
}

func (this *asyncLogRoutine) addAsyncLogger(key string, logger *asyncLogger) {
	this.addCh <- &asyncLoggerAddChItem{key, logger}
}

func (this *asyncLogRoutine) delAsyncLogger(key string) {
	this.delCh <- key
}

func (this *asyncLogRoutine) run() {
	for {
		select {
		case item, _ := <-this.addCh:
			this.asyncLoggers[item.key] = item.logger
		case key, _ := <-this.delCh:
			delete(this.asyncLoggers, key)
		case am, _ := <-this.msgCh:
			this.logAsyncMsg(am)
		case logger, _ := <-this.flushCh:
			logger.Flush()
		case <-this.freeCh:
			for len(this.msgCh) != 0 {
				am, _ := <-this.msgCh
				this.logAsyncMsg(am)
			}
			for key, al := range this.asyncLoggers {
				al.logger.Free()
				delete(this.asyncLoggers, key)
			}
			this.freeCh <- 1
			return
		}
	}
}

func (this *asyncLogRoutine) logAsyncMsg(am *asyncMsg) {
	am.al.logger.Log(am.level, am.msg)

	<-am.al.waitFreeLockCh
	am.al.msgCnt--
	if am.al.waitFree && am.al.msgCnt == 0 {
		am.al.logger.Free()
		delete(this.asyncLoggers, am.al.key)
	}
	am.al.waitFreeLockCh <- 1
}

/**  @} */

/**
* @name async logger
* @{ */

type asyncLogger struct {
	logger ILogger
	key    string

	msgCnt         int
	waitFree       bool
	waitFreeLockCh chan int
}

func NewAsyncLogger(logger ILogger) *asyncLogger {
	this := &asyncLogger{
		logger: logger,

		msgCnt:         0,
		waitFree:       false,
		waitFreeLockCh: make(chan int, 1),
	}

	this.key = fmt.Sprintf("%p", this)
	this.waitFreeLockCh <- 1
	alr.addAsyncLogger(this.key, this)

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
	<-this.waitFreeLockCh
	if this.msgCnt == 0 {
		this.logger.Free()
		alr.delAsyncLogger(this.key)
	} else {
		this.waitFree = true
	}
	this.waitFreeLockCh <- 1
}

/**  @} */
