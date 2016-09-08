package task

import (
	"fmt"
	"net/url"
	"testing"
)

func TestAdd(t *testing.T) {
	testTask := NewTask()
	taskKey := "gen_task_contents"
	testTask.Add(taskKey, genTaskContents)
	tf := testTask.FindTaskFunc(taskKey)
	if tf == nil {
		t.Errorf("add task fail")
	}
}

func TestRun(t *testing.T) {
	testTask := NewTask()
	taskKey := "gen_task_contents"
	testTask.Add(taskKey, genTaskContents)

	taskParams := "taskName=gen_task_contents&clusters=prelease_clusters"
	err := testTask.Run(taskParams)
	if err != nil {
		t.Errorf("run task fail, error: %v", err)
	}
}

func genTaskContents(params *url.Values) {
	contents := "hello task!" + params.Get("clusters")
	fmt.Println(contents)
}
