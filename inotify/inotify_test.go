package inotify

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestInotify(t *testing.T) {
	path := "/tmp/a.log"

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go watch(path, wg)

	wg.Add(1)
	go watch(path, wg)

	wg.Wait()
}

func watch(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	watcher, _ := NewWatcher()
	watcher.AddWatch(path, IN_ALL_EVENTS)
	watcher.AddWatch(filepath.Dir(path), IN_ALL_EVENTS)

	events, _ := watcher.ReadEvents()
	for _, event := range events {
		showEvent(event, watcher.fd)
	}

	os.OpenFile(path, os.O_RDONLY, 0)
	watcher.RmWatch(path)
	watcher.AddWatch(path, IN_ALL_EVENTS)

	for i := 0; i < 50; i++ {
		events, _ := watcher.ReadEvents()
		for _, event := range events {
			if watcher.IsUnreadEvent(event) {
				fmt.Println("it is a last remaining event")
			}
			showEvent(event, watcher.fd)
		}
	}

	watcher.Free()
	fmt.Println("bye")
}

func showEvent(event *Event, fd int) {
	fmt.Println(fd, event)

	if event.InIgnored() {
		fmt.Println(fd, event.wd, "IN_IGNORED")
	}

	if event.InAttrib() {
		fmt.Println(fd, event.wd, "IN_ATTRIB")
	}

	if event.InModify() {
		fmt.Println(fd, event.wd, "IN_MODIFY")
	}

	if event.InMoveSelf() {
		fmt.Println(fd, event.wd, "IN_MOVE_SELF")
	}

	if event.InMovedFrom() {
		fmt.Println(fd, event.wd, "IN_MOVED_FROM")
	}

	if event.InMovedTo() {
		fmt.Println(fd, event.wd, "IN_MOVED_TO")
	}

	if event.InDeleteSelf() {
		fmt.Println(fd, event.wd, "IN_DELETE_SELF")
	}

	if event.InDelete() {
		fmt.Println(fd, event.wd, "IN_DELETE")
	}

	if event.InCreate() {
		fmt.Println(fd, event.wd, "IN_CREATE")
	}
}
