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

	taskName := "gen_task_contents"
	paramStr := "clusters=prelease_clusters&url=http://hao.360.cn"
	err := testTask.Run(taskName, paramStr)
	if err != nil {
		t.Errorf("run task fail, error: %v", err)
	}
}

func genTaskContents(params url.Values) {
	contents := "hello task! clusters: " + params.Get("clusters") + " url: " + params.Get("url")
	fmt.Println(contents)
}
