package inotify

import "syscall"

const (
	IN_ALL_EVENTS = syscall.IN_ALL_EVENTS
	IN_IGNORED    = syscall.IN_IGNORED

	IN_MODIFY = syscall.IN_MODIFY
	IN_ATTRIB = syscall.IN_ATTRIB

	IN_MOVE_SELF  = syscall.IN_MOVE_SELF
	IN_MOVED_FROM = syscall.IN_MOVED_FROM
	IN_MOVED_TO   = syscall.IN_MOVED_TO

	IN_DELETE_SELF = syscall.IN_DELETE_SELF
	IN_DELETE      = syscall.IN_DELETE

	IN_CREATE = syscall.IN_CREATE
)

type Event struct {
	wd     uint32
	mask   uint32
	cookie uint32

	Path string
	Name string
}

func (this *Event) InIgnored() bool {
	return this.mask&IN_IGNORED == IN_IGNORED
}

func (this *Event) InModify() bool {
	return this.mask&IN_MODIFY == IN_MODIFY
}

func (this *Event) InAttrib() bool {
	return this.mask&IN_ATTRIB == IN_ATTRIB
}

func (this *Event) InMoveSelf() bool {
	return this.mask&IN_MOVE_SELF == IN_MOVE_SELF
}

func (this *Event) InMovedFrom() bool {
	return this.mask&IN_MOVED_FROM == IN_MOVED_FROM
}

func (this *Event) InMovedTo() bool {
	return this.mask&IN_MOVED_TO == IN_MOVED_TO
}

func (this *Event) InDeleteSelf() bool {
	return this.mask&IN_DELETE_SELF == IN_DELETE_SELF
}

func (this *Event) InDelete() bool {
	return this.mask&IN_DELETE == IN_DELETE
}

func (this *Event) InCreate() bool {
	return this.mask&IN_CREATE == IN_CREATE
}
