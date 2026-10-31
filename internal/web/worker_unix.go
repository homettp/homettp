//go:build unix

package web

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	shellName = "/bin/sh"
	shellArg  = "-c"

	payloadReference = `${` + payloadEnvName + `:+"$` + payloadEnvName + `"}`
)

func configureProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}

		return err
	}
}

func payloadArgument(string) string {
	return payloadReference
}

func isSafePayload(payload string) bool {
	return !strings.ContainsRune(payload, 0)
}
