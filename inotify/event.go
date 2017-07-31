package inotify

import "syscall"

const (
	IN_ALL_EVENTS = syscall.IN_ALL_EVENTS
	IN_IGNORED    = syscall.IN_IGNORED

	IN_MODIFY      = syscall.IN_MODIFY
	IN_DELETE_SELF = syscall.IN_DELETE_SELF
)

type Event struct {
	mask   uint32
	cookie uint32

	Name string
}

func (this *Event) InIgnored() bool {
	return this.mask&IN_IGNORED == IN_IGNORED
}

func (this *Event) InModify() bool {
	return this.mask&IN_MODIFY == IN_MODIFY
}

func (this *Event) InDeleteSelf() bool {
	return this.mask&IN_DELETE_SELF == IN_DELETE_SELF
}
