/**
* @file async.go
* @brief write msg in independency goroutine
* @author ligang
* @date 2016-07-19
 */

package log

var alr *asyncLogRoutine

// must be called first
func InitAsyncLogRoutine(msgQueueLen int, maxLoggerNum int) {
	alr = &asyncLogRoutine{
		msgCh:   make(chan *asyncMsg, msgQueueLen),
		flushCh: make(chan ILogger, maxLoggerNum),
		freeCh:  make(chan int),

		addCh: make(chan *asyncLogger, maxLoggerNum),
		delCh: make(chan *asyncLogger, maxLoggerNum),

		asyncLoggers: make(map[uint64]*asyncLogger),
	}

	go alr.run()
}

func FreeAsyncLogRoutine() {
	alr.freeCh <- 1
	<-alr.freeCh

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

type asyncLogRoutine struct {
	msgCh   chan *asyncMsg
	flushCh chan ILogger
	freeCh  chan int

	addCh chan *asyncLogger
	delCh chan *asyncLogger

	curIndex     uint64
	asyncLoggers map[uint64]*asyncLogger
}

func (this *asyncLogRoutine) addAsyncLogger(logger *asyncLogger) {
	this.addCh <- logger
}

func (this *asyncLogRoutine) delAsyncLogger(logger *asyncLogger) {
	this.delCh <- logger
}

func (this *asyncLogRoutine) freeAsyncLogger(logger *asyncLogger) {
	logger.logger.Free()
	delete(this.asyncLoggers, logger.index)
}

func (this *asyncLogRoutine) run() {
	for {
		select {
		case logger, _ := <-this.addCh:
			logger.index = this.curIndex
			this.asyncLoggers[this.curIndex] = logger
			this.curIndex++
		case logger, _ := <-this.delCh:
			this.freeAsyncLogger(logger)
		case am, _ := <-this.msgCh:
			this.logAsyncMsg(am)
		case logger, _ := <-this.flushCh:
			logger.Flush()
		case <-this.freeCh:
			for len(this.msgCh) != 0 {
				am, _ := <-this.msgCh
				this.logAsyncMsg(am)
			}
			for _, al := range this.asyncLoggers {
				this.freeAsyncLogger(al)
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
		this.freeAsyncLogger(am.al)
	}
	am.al.waitFreeLockCh <- 1
}

/**  @} */

/**
* @name async logger
* @{ */

type asyncLogger struct {
	logger ILogger
	index  uint64

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

	this.waitFreeLockCh <- 1
	alr.addAsyncLogger(this)

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
		alr.delAsyncLogger(this)
	} else {
		this.waitFree = true
	}
	this.waitFreeLockCh <- 1
}

/**  @} */
