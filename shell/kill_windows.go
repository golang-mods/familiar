//go:build windows

package shell

import (
	"os/exec"
	"strconv"
	"syscall"
)

func prepareCommand(command *exec.Cmd) {}

func killChildren(command *exec.Cmd, signal syscall.Signal) error {
	kill := exec.Command("taskkill", "/t", "/f", "/pid", strconv.Itoa(command.Process.Pid))
	kill.Stderr = command.Stderr
	kill.Stdout = command.Stdout
	return kill.Run()
}
