package inotify

import (
	"errors"
	"strings"
	"syscall"
	"unsafe"
)

type Watcher struct {
	fd int

	pathToWdMap map[string]uint32
	wdToPathMap map[uint32]string
}

func NewWatcher() (*Watcher, error) {
	fd, err := syscall.InotifyInit()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		fd: fd,

		pathToWdMap: make(map[string]uint32),
		wdToPathMap: make(map[uint32]string),
	}, nil
}

func (this *Watcher) AddWatch(path string, mask uint32) error {
	wd, err := syscall.InotifyAddWatch(this.fd, path, mask)
	if err != nil {
		return err
	}

	uwd := uint32(wd)
	this.pathToWdMap[path] = uwd
	this.wdToPathMap[uwd] = path

	return nil
}

func (this *Watcher) RmWatch(path string) {
	wd, ok := this.pathToWdMap[path]
	if !ok {
		return
	}

	syscall.InotifyRmWatch(this.fd, wd)
	delete(this.pathToWdMap, path)
	delete(this.wdToPathMap, wd)
}

func (this *Watcher) ReadEvents() ([]*Event, error) {
	buf := make([]byte, syscall.SizeofInotifyEvent*4096)

	n, err := syscall.Read(this.fd, buf)
	if n == 0 {
		return nil, errors.New("Read 0 byte error")
	}
	if err != nil {
		return nil, err
	}

	var offset uint32
	var events []*Event
	for offset <= uint32(n-syscall.SizeofInotifyEvent) {
		ie := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[offset]))
		event := &Event{
			wd:     uint32(ie.Wd),
			mask:   ie.Mask,
			cookie: ie.Cookie,
		}
		event.path = this.wdToPathMap[event.wd]

		offset += syscall.SizeofInotifyEvent
		if ie.Len > 0 {
			nameBytes := (*[syscall.PathMax]byte)(unsafe.Pointer(&buf[offset]))
			event.Name = strings.TrimRight(string(nameBytes[0:ie.Len]), "\000")
			offset += ie.Len
		}

		events = append(events, event)
	}

	return events, nil
}

func (this *Watcher) IsUnreadEvent(event *Event) bool {
	if event.wd != this.pathToWdMap[event.path] {
		return true
	}

	return false
}

func (this *Watcher) Free() {
	for path, _ := range this.pathToWdMap {
		this.RmWatch(path)
	}
	syscall.Close(this.fd)
}
