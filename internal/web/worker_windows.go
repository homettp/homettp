//go:build windows

package web

import (
	"os/exec"
	"strings"
	"syscall"
)

const (
	shellName = "cmd.exe"
	shellArg  = "/C"

	payloadReference = `"%` + payloadEnvName + `%"`
)

func configureProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: shellName + ` /S ` + shellArg + ` "` + cmd.Args[2] + `"`,
	}
}

func payloadArgument(payload string) string {
	if payload == "" {
		return ""
	}

	return payloadReference
}

func isSafePayload(payload string) bool {
	return !strings.ContainsAny(payload, "\x00\r\n\"%!^&|<>()") && !strings.HasSuffix(payload, `\`)
}
