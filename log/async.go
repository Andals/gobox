/**
* @file async.go
* @brief write msg in independency goroutine
* @author ligang
* @date 2016-07-19
 */

package log

type asyncMsg struct {
	level int
	msg   []byte

	al *asyncLogger
}

/**
* @name async log routine
* @{ */

type AsyncLogRoutineCh struct {
	msgCh   chan *asyncMsg
	flushCh chan ILogger
	freeCh  chan int
}

var asyncLogRoutineList []*AsyncLogRoutineCh

func NewAsyncLogRoutine(queueLen int) *AsyncLogRoutineCh {
	this := &AsyncLogRoutineCh{
		msgCh:   make(chan *asyncMsg, queueLen),
		flushCh: make(chan ILogger, queueLen),
		freeCh:  make(chan int),
	}

	go logRoutine(this)
	asyncLogRoutineList = append(asyncLogRoutineList, this)

	return this
}

func FreeAsyncLogRoutines() {
	for _, ach := range asyncLogRoutineList {
		ach.Free()
	}
}

func (this *AsyncLogRoutineCh) Free() {
	this.freeCh <- 1
	<-this.freeCh

	close(this.msgCh)
	close(this.flushCh)
	close(this.freeCh)
}

func logRoutine(ach *AsyncLogRoutineCh) {
	for {
		select {
		case am, _ := <-ach.msgCh:
			logAsyncMsg(am)
		case logger, _ := <-ach.flushCh:
			logger.Flush()
		case <-ach.freeCh:
			for len(ach.msgCh) != 0 {
				am, _ := <-ach.msgCh
				logAsyncMsg(am)
			}
			ach.freeCh <- 1
			return
		}
	}
}

func logAsyncMsg(am *asyncMsg) {
	am.al.logger.Log(am.level, am.msg)
	am.al.msgCnt--

	if am.al.msgCnt == 0 && am.al.waitFree {
		am.al.logger.Free()
	}
}

/**  @} */

/**
* @name async logger
* @{ */

type asyncLogger struct {
	msgCnt   int
	waitFree bool

	logger ILogger
	ach    *AsyncLogRoutineCh
}

func NewAsyncLogger(logger ILogger, ach *AsyncLogRoutineCh) *asyncLogger {
	this := &asyncLogger{
		msgCnt:   0,
		waitFree: false,

		logger: logger,
		ach:    ach,
	}

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
	this.ach.msgCh <- am

	return nil
}

func (this *asyncLogger) Flush() error {
	this.ach.flushCh <- this.logger

	return nil
}

func (this *asyncLogger) Free() {
	this.waitFree = true
}

/**  @} */
