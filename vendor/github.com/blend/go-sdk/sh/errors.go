package sh

import (
	"os/exec"
	"syscall"
)

// IsEPIPE is the epipe erorr.
func IsEPIPE(err error) bool {
	if typed, ok := err.(*exec.ExitError); ok {
		status := typed.Sys().(syscall.WaitStatus)
		if status.Signal() == syscall.SIGPIPE {
			return true
		}
	}
	return false
}
