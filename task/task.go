/**
* @file task.go
* @brief tool for exec task
* @author huangqiuping
* @date 2016-09-07
 */

package task

import (
	"errors"
	"net/url"
)

type TaskFunc func(params url.Values)

type Task struct {
	taskList map[string]TaskFunc
}

func NewTask() *Task {
	this := new(Task)
	this.taskList = make(map[string]TaskFunc)

	return this
}

func (this *Task) FindTaskFunc(taskName string) TaskFunc {
	tf, ok := this.taskList[taskName]
	if ok {
		return tf
	}
	return nil
}

func (this *Task) Add(key string, tf TaskFunc) {
	this.taskList[key] = tf
}

func (this *Task) Run(taskName, paramStr string) error {
	if taskName == "" {
		return errors.New("need contain param taskName")
	}

	values, err := url.ParseQuery(paramStr)
	if err != nil {
		return err
	}

	tf := this.FindTaskFunc(taskName)
	if tf == nil {
		return errors.New("taskName invalid")
	}
	tf(values)

	return nil
}
