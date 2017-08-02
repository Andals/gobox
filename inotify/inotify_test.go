package inotify

import (
	"fmt"
	"testing"
	"os"
	"path/filepath"
)

func TestInotify(t *testing.T) {
	path := "/tmp/a.log"

	go watch(path)
	go watch(path)

	<-make(chan bool)
}

func watch(path string){
	watcher, _ := NewWatcher()
	watcher.AddWatch(path, IN_ALL_EVENTS)
	watcher.AddWatch(filepath.Dir(path), IN_ALL_EVENTS)

	events, _ := watcher.ReadEvents()
	for _, event := range events {
		showEvent(event,watcher.fd)
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
			showEvent(event,watcher.fd)
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

	if event.InModify() {
		fmt.Println(fd, event.wd, "IN_MODIFY")
	}

	if event.InDeleteSelf() {
		fmt.Println(fd, event.wd,"IN_DELETE_SELF")
	}

	if event.InMoveSelf() {
		fmt.Println(fd, event.wd,"IN_MOVE_SELF")
	}

	if event.InDelete(){
		fmt.Println(fd, event.wd,"IN_DELETE")
	}

	if event.InAttrib() {
		fmt.Println(fd, event.wd,"IN_ATTRIB")
	}
}
