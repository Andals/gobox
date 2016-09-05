package postoffice

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

const (
	AddressSplitChar = " "
)

type Postman struct {
}

func NewPostman() *Postman {
	return &Postman{}
}

func (pm *Postman) Send(em *Email) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("TODO: cannot send email on windows")
	}

	cmd := exec.Command("mail", "-s", em.Subject)
	cmd.Args = append(cmd.Args, "-c", strings.Join(em.Cc, AddressSplitChar))
	cmd.Args = append(cmd.Args, "-b", strings.Join(em.Bcc, AddressSplitChar))
	cmd.Args = append(cmd.Args, strings.Join(em.To, AddressSplitChar))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("StdinPipe failed to perform: %s (Command: %s, Arguments: %s)", err, cmd.Path, cmd.Args)
	}

	stdin.Write([]byte(em.Content))
	stdin.Close()

	_, err = cmd.Output()
	if err != nil || !cmd.ProcessState.Success() {
		return fmt.Errorf("Send mail error : %v ", err)
	}

	return nil
}
