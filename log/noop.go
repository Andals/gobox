package log

type NoopLogger struct {
}

func (this *NoopLogger) Debug(msg []byte) {
}

func (this *NoopLogger) Info(msg []byte) {
}

func (this *NoopLogger) Notice(msg []byte) {
}

func (this *NoopLogger) Warning(msg []byte) {
}

func (this *NoopLogger) Error(msg []byte) {
}

func (this *NoopLogger) Critical(msg []byte) {
}

func (this *NoopLogger) Alert(msg []byte) {
}

func (this *NoopLogger) Emergency(msg []byte) {
}

func (this *NoopLogger) Log(level int, msg []byte) error {
	return nil
}

func (this *NoopLogger) Flush() error {
	return nil
}

func (this *NoopLogger) Free() {
}
