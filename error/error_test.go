package error

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	e := NewError(101, "test error")

	fmt.Println(e.GetErrno(), e.GetMsg())
	fmt.Println(e.Error())
}
