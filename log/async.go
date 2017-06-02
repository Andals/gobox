/**
* @file async.go
* @brief write msg in independency goroutine
* @author ligang
* @date 2016-07-19
 */

package log

var alr *asyncLogRoutine

// must be called first
func InitAsyncLogRoutine(msgQueueLen int) {
	alr = &asyncLogRoutine{
		msgCh:   make(chan *asyncMsg, msgQueueLen),
		flushCh: make(chan ILogger, msgQueueLen),
		freeCh:  make(chan int),
	}

	go alr.run()
}

func FreeAsyncLogRoutine() {
	alr.freeCh <- 1
	<-alr.freeCh
}

type asyncMsg struct {
	level int
	msg   []byte

	alogger *asyncLogger
}

/**
* @name async log routine
* @{ */

type asyncLogRoutine struct {
	msgCh   chan *asyncMsg
	flushCh chan ILogger
	freeCh  chan int
}

func (this *asyncLogRoutine) run() {
	for {
		select {
		case am, _ := <-this.msgCh:
			logAsyncMsg(am)
		case logger, _ := <-this.flushCh:
			logger.Flush()
		case <-this.freeCh:
			for len(this.msgCh) != 0 {
				am, _ := <-this.msgCh
				logAsyncMsg(am)
			}
			this.freeCh <- 1
			return
		}
	}
}

func logAsyncMsg(am *asyncMsg) {
	am.alogger.logger.Log(am.level, am.msg)
	am.alogger.msgCnt--

	if am.alogger.msgCnt == 0 {
		<-am.alogger.freeLockCh
		if am.alogger.waitFree {
			am.alogger.logger.Free()
		}
		am.alogger.freeLockCh <- 1
	}
}

/**  @} */

/**
* @name async logger
* @{ */

type asyncLogger struct {
	logger ILogger

	msgCnt     int
	waitFree   bool
	freeLockCh chan int
}

func NewAsyncLogger(logger ILogger) *asyncLogger {
	this := &asyncLogger{
		logger: logger,

		msgCnt:     0,
		waitFree:   false,
		freeLockCh: make(chan int, 1),
	}

	this.freeLockCh <- 1

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

		alogger: this,
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
	} else {
		<-this.freeLockCh
		this.waitFree = true
		this.freeLockCh <- 1
	}
}

/**  @} */
