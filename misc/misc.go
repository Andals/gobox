/**
* @file misc.go
* @brief misc supermarket
* @author ligang
* @date 2015-12-11
 */

package misc

import (
	"github.com/andals/gobox/color"

	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func IntSliceUnique(s []int) []int {
	m := make(map[int]bool)
	r := make([]int, 0, cap(s))

	for _, k := range s {
		_, ok := m[k]
		if !ok {
			r = append(r, k)
			m[k] = true
		}
	}

	return r
}

func StringSliceUnique(s []string) []string {
	m := make(map[string]bool)
	r := make([]string, 0, cap(s))

	for _, k := range s {
		_, ok := m[k]
		if !ok {
			r = append(r, k)
			m[k] = true
		}
	}

	return r
}

func FileExist(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	return true
}

func DirExist(path string) bool {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

func PrintCallerFuncNameForTest() {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	c := color.Yellow([]byte(f.Name()))
	fmt.Println(string(c))
}

func AppendBytes(b []byte, elems ...[]byte) []byte {
	buf := bytes.NewBuffer(b)
	for _, e := range elems {
		buf.Write(e)
	}

	return buf.Bytes()
}

func ListFilesInDir(rootDir string) ([]string, error) {
	rootDir = strings.TrimRight(rootDir, "/")
	if !DirExist(rootDir) {
		return nil, errors.New("Dir not exists")
	}

	var fileList []string
	dirList := []string{rootDir}

	for i := 0; i < len(dirList); i++ {
		curDir := dirList[i]
		file, err := os.Open(dirList[i])
		if err != nil {
			return nil, err
		}

		fis, err := file.Readdir(-1)
		if err != nil {
			return nil, err
		}

		for _, fi := range fis {
			path := curDir + "/" + fi.Name()
			if fi.IsDir() {
				dirList = append(dirList, path)
			} else {
				fileList = append(fileList, path)
			}
		}
	}

	return fileList, nil
}
