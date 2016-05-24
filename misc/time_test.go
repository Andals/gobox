package misc

import (
	"fmt"
	"testing"
	"time"
)

func TestRandByTime(t *testing.T) {
	PrintCallerFuncNameForTest()

	tm := time.Now()
	fmt.Println(RandByTime(&tm), RandByTime(&tm), RandByTime(nil))
}
