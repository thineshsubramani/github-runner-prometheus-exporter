package platform

import (
	"os/exec"
	"runtime"
	"strings"
)

func IsRunnerProcessRunning(processName string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	} else {
		cmd = exec.Command("pgrep", "-f", processName)
	}
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), processName)
}
