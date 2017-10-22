package misc

import (
	"testing"
)

func TestBoolSet(t *testing.T) {
	bs := BoolSet(map[string]bool{
		"a": true,
		"b": false,
	})

	if !bs.IsTrue("a") {
		t.Error("key a error")
	}

	if bs.IsTrue("b") {
		t.Error("key b error")
	}
}
