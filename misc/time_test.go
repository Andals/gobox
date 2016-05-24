package misc

import (
	"fmt"
	"testing"
)

func TestRandByTime(t *testing.T) {
	PrintCallerFuncNameForTest()

	fmt.Println(RandByTime(), RandByTime(), RandByTime())
}
