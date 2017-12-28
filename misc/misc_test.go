package misc

import (
	"fmt"
	"testing"
)

func TestIntSliceUnique(t *testing.T) {
	PrintCallerFuncNameForTest()

	s := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}

	fmt.Println("origin slice is: ", s)

	s = IntSliceUnique(s)

	fmt.Println("after call slice is: ", s)
}

func TestStringSliceUnique(t *testing.T) {
	PrintCallerFuncNameForTest()

	s := []string{"a", "ab", "ab", "abc", "abc", "abc", "abcd", "abcd", "abcd", "abcd", "abcd"}

	fmt.Println("origin slice is: ", s)

	s = StringSliceUnique(s)

	fmt.Println("after call slice is: ", s)
}

func TestFileExist(t *testing.T) {
	PrintCallerFuncNameForTest()

	f := "/etc/passwd"

	r := FileExist(f)
	if r {
		fmt.Println(f, "is exist")
	} else {
		fmt.Println(f, "is not exist")
	}
}

func TestDirExist(t *testing.T) {
	PrintCallerFuncNameForTest()

	d := "/home/ligang/devspace"

	r := DirExist(d)
	if r {
		fmt.Println(d, "is exist")
	} else {
		fmt.Println(d, "is not exist")
	}
}

func TestAppendBytes(t *testing.T) {
	PrintCallerFuncNameForTest()

	b := []byte("abc")
	b = AppendBytes(b, []byte("def"), []byte("ghi"))

	fmt.Println(string(b))
}

func TestListFilesInDir(t *testing.T) {
	PrintCallerFuncNameForTest()

	fileList, err := ListFilesInDir("/home/ligang/tmp")
	if err != nil {
		t.Log(err)
		return
	}

	for _, path := range fileList {
		t.Log(path)
	}
}
