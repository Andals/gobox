package inotify

import (
	"errors"
	"strings"
	"syscall"
	"unsafe"
)

type Watcher struct {
	fd  int
	wds map[string]uint32
}

func NewWatcher() (*Watcher, error) {
	fd, err := syscall.InotifyInit()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		fd:  fd,
		wds: make(map[string]uint32),
	}, nil
}

func (this *Watcher) AddWatch(path string, mask uint32) error {
	wd, err := syscall.InotifyAddWatch(this.fd, path, mask)
	if err != nil {
		return err
	}

	this.wds[path] = uint32(wd)
	return nil
}

func (this *Watcher) RmWatch(path string) {
	wd, ok := this.wds[path]
	if !ok {
		return
	}

	syscall.InotifyRmWatch(this.fd, wd)
	delete(this.wds, path)
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
			mask:   ie.Mask,
			cookie: ie.Cookie,
		}

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
