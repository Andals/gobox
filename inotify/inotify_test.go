package inotify

import (
	"fmt"
	"testing"
)

func TestInotify(t *testing.T) {
	path := "/tmp/a.log"
	watcher, _ := NewWatcher()
	watcher.AddWatch(path, IN_ALL_EVENTS)

	events, _ := watcher.ReadEvents()
	for _, event := range events {
		showEvent(event)
	}

	watcher.RmWatch(path)
	watcher.AddWatch(path, IN_ALL_EVENTS)

	for {
		events, _ := watcher.ReadEvents()
		for _, event := range events {
			if watcher.IsUnreadEvent(event) {
				fmt.Println("it is a last remaining event")
			}
			showEvent(event)
		}
	}
}

func showEvent(event *Event) {
	fmt.Println(event)

	if event.InIgnored() {
		fmt.Println("IN_IGNORED")
	}

	if event.InModify() {
		fmt.Println("IN_MODIFY")
	}

	if event.InDeleteSelf() {
		fmt.Println("IN_DELETE_SELF")
	}

	if event.InMoveSelf() {
		fmt.Println("IN_MOVE_SELF")
	}
}
