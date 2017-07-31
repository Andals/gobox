package inotify

import (
	"fmt"
	"testing"
)

func TestInotify(t *testing.T) {
	path := "/tmp/a.log"
	watcher, _ := NewWatcher()
	watcher.AddWatch(path, IN_MODIFY|IN_DELETE_SELF)

	for {
		events, _ := watcher.ReadEvents()
		for _, event := range events {
			quit := showEvent(event)
			if quit {
				watcher.RmWatch(path)
				return
			}
		}
	}
}

func showEvent(event *Event) bool {
	fmt.Println(event)

	if event.InIgnored() {
		fmt.Println("IN_IGNORED")
		return true
	}

	if event.InModify() {
		fmt.Println("IN_MODIFY")
	}

	if event.InDeleteSelf() {
		fmt.Println("IN_DELETE_SELF")
	}

	return false
}
