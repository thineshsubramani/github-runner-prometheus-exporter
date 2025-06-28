package watcher

import (
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type RunnerState struct {
	RunnerName string
	State      string
	JobID      string
}

var (
	stateMu sync.RWMutex
	runners = make(map[string]RunnerState)
)

func SetRunnerState(name string, state RunnerState) {
	stateMu.Lock()
	defer stateMu.Unlock()
	runners[name] = state
}

func GetRunnerState(name string) (RunnerState, bool) {
	stateMu.RLock()
	defer stateMu.RUnlock()
	s, ok := runners[name]
	return s, ok
}

// Track proccess state for Runner its child
// Support multi OS Linux, Window etc
// Runner Lifecycle tracking by working with event.js file tracking
// Get Proccess ID
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
