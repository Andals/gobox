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

	logger ILogger
}

/**
* @name async log routine
* @{ */

type asyncLogRoutineCh struct {
	msgCh   chan *asyncMsg
	flushCh chan ILogger
	freeCh  chan int
}

var asyncLogRoutineList []asyncLogRoutineCh

func NewAsyncLogRoutine(queueLen int) asyncLogRoutineCh {
	this := asyncLogRoutineCh{
		msgCh:   make(chan *asyncMsg, queueLen),
		flushCh: make(chan ILogger, queueLen),
		freeCh:  make(chan int),
	}

	go logRoutine(this)
	asyncLogRoutineList = append(asyncLogRoutineList, this)

	return this
}

func (this *asyncLogRoutineCh) Free() {
	this.freeCh <- 1
	<-this.freeCh

	close(this.msgCh)
	close(this.flushCh)
	close(this.freeCh)
}

func logRoutine(ach asyncLogRoutineCh) {
	for {
		select {
		case am, _ := <-ach.msgCh:
			am.logger.Log(am.level, am.msg)
		case logger, _ := <-ach.flushCh:
			logger.Flush()
		case <-ach.freeCh:
			for len(ach.msgCh) != 0 {
				am, _ := <-ach.msgCh
				am.logger.Log(am.level, am.msg)
			}
			ach.freeCh <- 1
			return
		}
	}
}

/**  @} */

/**
* @name async logger
* @{ */

type asyncLogger struct {
	logger ILogger

	ach asyncLogRoutineCh
}

var asyncLoggerList []*asyncLogger

func NewAsyncLogger(logger ILogger, ach asyncLogRoutineCh) *asyncLogger {
	this := &asyncLogger{
		logger: logger,

		ach: ach,
	}

	asyncLoggerList = append(asyncLoggerList, this)

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
	this.ach.msgCh <- &asyncMsg{
		level: level,
		msg:   msg,

		logger: this.logger,
	}

	return nil
}

func (this *asyncLogger) Flush() error {
	this.ach.flushCh <- this.logger

	return nil
}

func (this *asyncLogger) Free() {
}

/**  @} */

func FreeAsyncLogList() {
	for _, ach := range asyncLogRoutineList {
		ach.Free()
	}

	for _, al := range asyncLoggerList {
		al.logger.Free()
	}
}
