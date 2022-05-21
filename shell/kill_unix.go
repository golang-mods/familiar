//go:build !windows

package shell

import (
	"os/exec"
	"syscall"
)

func prepareCommand(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killChildren(command *exec.Cmd, signal syscall.Signal) error {
	err := syscall.Kill(-command.Process.Pid, sig)
	if err == nil && sig != syscall.SIGKILL && sig != syscall.SIGCONT {
		err = syscall.Kill(-command.Process.Pid, syscall.SIGCONT)
	}
	return err
}
