/**
* @file shell.go
* @brief tool for exec shell cmd
* @author ligang
* @date 2016-01-28
 */

package shell

import (
	//     "fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

type ShellResult struct {
	Ok     bool
	Output string
}

func NewCmd(cmdStr string) *exec.Cmd {
	return exec.Command("/bin/bash", "-c", cmdStr)
}

func RunCmd(cmdStr string) *ShellResult {
	result := &ShellResult{
		Ok:     true,
		Output: "",
	}

	cmd := NewCmd(cmdStr)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if nil != err {
		result.Ok = false
		result.Output += err.Error()
	}
	return result
}

func RunCmdBindTerminal(cmdStr string) {
	cmd := NewCmd(cmdStr)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func RunAsUser(cmdStr string, username string) *ShellResult {
	var cmd string

	curUser, _ := user.Current()
	if "root" == curUser.Username {
		cmd += "/sbin/runuser " + username + " -c \""
		cmd += strings.Replace(cmdStr, "\"", "\\\"", -1)
		cmd += "\""
	} else {
		cmd += "sudo -u " + username + " "
		cmd += cmdStr
	}

	return RunCmd(cmd)
}

func Rsync(host string, sou string, dst string, excludeFrom string, sshUser string, timeout int) *ShellResult {
	rsyncCmd := MakeRsyncCmd(host, sou, dst, excludeFrom, timeout)

	return RunAsUser(rsyncCmd, sshUser)
}

func MakeRsyncCmd(host string, sou string, dst string, excludeFrom string, timeout int) string {
	to := strconv.Itoa(timeout)
	rsyncCmd := "/usr/bin/rsync -av -e 'ssh -o StrictHostKeyChecking=no -o ConnectTimeout=" + to + "' --timeout=" + to + " --update "
	_, err := os.Stat(excludeFrom)
	if nil == err {
		rsyncCmd += "--exclude-from='" + excludeFrom + "' "
	}
	rsyncCmd += sou + " " + host + ":" + dst + " 2>&1"

	return rsyncCmd
}
