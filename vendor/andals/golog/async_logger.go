/**
* @file async.go
* @brief write msg in independency goroutine
* @author ligang
* @date 2016-07-19
 */

package golog

const (
	ASYNC_MSG_KIND_LOG          = 1
	ASYNC_MSG_KIND_FLUSH        = 2
	ASYNC_MSG_KIND_FREE_LOGGER  = 3
	ASYNC_MSG_KIND_FREE_ROUTINE = 4
)

type asyncMsg struct {
	kind   int
	logger ILogger

	level int
	msg   []byte
}

var alr *asyncLogRoutine

// must be called first
func InitAsyncLogRoutine(msgQueueLen int) {
	alr = &asyncLogRoutine{
		msgCh:  make(chan *asyncMsg, msgQueueLen),
		freeCh: make(chan int),
	}

	go alr.run()
}

func FreeAsyncLogRoutine() {
	alr.msgCh <- &asyncMsg{
		kind: ASYNC_MSG_KIND_FREE_ROUTINE,
	}
	<-alr.freeCh
}

/**
* @name async log routine
* @{ */

type asyncLogRoutine struct {
	msgCh  chan *asyncMsg
	freeCh chan int
}

func (this *asyncLogRoutine) run() {
	for {
		select {
		case am, _ := <-this.msgCh:
			this.processAsyncMsg(am)
		}
	}
}

func (this *asyncLogRoutine) processAsyncMsg(am *asyncMsg) {
	switch am.kind {
	case ASYNC_MSG_KIND_LOG:
		am.logger.Log(am.level, am.msg)
	case ASYNC_MSG_KIND_FLUSH:
		am.logger.Flush()
	case ASYNC_MSG_KIND_FREE_LOGGER:
		am.logger.Free()
	case ASYNC_MSG_KIND_FREE_ROUTINE:
		this.freeCh <- 1
	}
}

/**  @} */

/**
* @name async logger
* @{ */

type asyncLogger struct {
	logger ILogger
}

func NewAsyncLogger(logger ILogger) *asyncLogger {
	this := &asyncLogger{
		logger: logger,
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
	alr.msgCh <- &asyncMsg{
		kind:   ASYNC_MSG_KIND_LOG,
		logger: this.logger,

		msg:   msg,
		level: level,
	}

	return nil
}

func (this *asyncLogger) Flush() error {
	alr.msgCh <- &asyncMsg{
		kind:   ASYNC_MSG_KIND_FLUSH,
		logger: this.logger,
	}

	return nil
}

func (this *asyncLogger) Free() {
	alr.msgCh <- &asyncMsg{
		kind:   ASYNC_MSG_KIND_FREE_LOGGER,
		logger: this.logger,
	}
}

/**  @} */
