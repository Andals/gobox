package query

import (
	"fmt"
	"net/url"
	"testing"

	"andals/gobox/misc"
)

func TestParse(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	qv, _ := url.ParseQuery("a=1&b=hello&c=64")
	qs := NewQuerySet()

	var a int
	var b string
	var c int64

	qs.IntVar(&a, "a", 101, "invalid a", CheckIntNotZero)
	qs.StringVar(&b, "b", 102, "invalid b", CheckStringNotEmpty)
	qs.Int64Var(&c, "c", 103, "invalid c", CheckInt64NotZero)

	e := qs.Parse(&qv)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(a, b, c)
	}
}
