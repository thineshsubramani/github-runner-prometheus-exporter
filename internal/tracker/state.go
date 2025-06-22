package tracker

import "sync"

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
