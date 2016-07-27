package shell

import (
	"fmt"
	"testing"
)

func TestRunCmd(t *testing.T) {
	result := RunCmd("ls -l")
	fmt.Println(result.Ok, string(result.Output))
}

func TestRunAsUser(t *testing.T) {
	result := RunAsUser("ls -l", "root")
	fmt.Println(result.Ok, string(result.Output))
}

func TestRsync(t *testing.T) {
	host := "10.16.57.92"
	sou := "/home/ligang/tmp/go/*"
	dst := "/home/ligang/tmp/go"
	sshUser := "ligang"

	result := Rsync(host, sou, dst, "", sshUser, 3)
	fmt.Println(result.Ok, string(result.Output))
}
