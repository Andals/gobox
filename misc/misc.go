/**
* @file misc.go
* @brief misc supermarket
* @author ligang
* @date 2015-12-11
 */

package misc

import (
	"crypto/md5"
	"fmt"
	"os"
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

func CalMd5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
