package inotify

import "syscall"

const (
	IN_ALL_EVENTS = syscall.IN_ALL_EVENTS
	IN_IGNORED    = syscall.IN_IGNORED

	IN_MODIFY      = syscall.IN_MODIFY
	IN_DELETE_SELF = syscall.IN_DELETE_SELF
	IN_MOVE_SELF   = syscall.IN_MOVE_SELF
)

type Event struct {
	wd     uint32
	mask   uint32
	cookie uint32
	path   string

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

func (this *Event) InMoveSelf() bool {
	return this.mask&IN_MOVE_SELF == IN_MOVE_SELF
}
