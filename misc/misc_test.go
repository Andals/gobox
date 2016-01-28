package misc

import (
	"fmt"
	"testing"
)

func TestIntSliceUnique(t *testing.T) {
	s := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}

	fmt.Println("origin slice is: ", s)

	s = IntSliceUnique(s)

	fmt.Println("after call slice is: ", s)
}

func TestStringSliceUnique(t *testing.T) {
	s := []string{"a", "ab", "ab", "abc", "abc", "abc", "abcd", "abcd", "abcd", "abcd", "abcd"}

	fmt.Println("origin slice is: ", s)

	s = StringSliceUnique(s)

	fmt.Println("after call slice is: ", s)
}

func TestFileExist(t *testing.T) {
	f := "/etc/passwd"

	r := FileExist(f)
	if r {
		fmt.Println(f, "is exist")
	} else {
		fmt.Println(f, "is not exist")
	}
}

func TestDirExist(t *testing.T) {
	d := "/home/ligang/devspace"

	r := DirExist(d)
	if r {
		fmt.Println(d, "is exist")
	} else {
		fmt.Println(d, "is not exist")
	}
}
