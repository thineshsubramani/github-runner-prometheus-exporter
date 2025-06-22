package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/platform"
)

// Check if process exists
func ValidateRunnerProcess(processName string) error {
	if !platform.IsRunnerProcessRunning(processName) {
		return fmt.Errorf("runner process %q not running", processName)
	}
	return nil
}

// Validate directory and required files exist
func ValidatePaths(basePath string) error {
	// Check if base path exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return fmt.Errorf("base path does not exist: %s", basePath)
	}

	// Validate .runner
	runnerPath := filepath.Join(basePath, ".runner")
	if _, err := os.Stat(runnerPath); os.IsNotExist(err) {
		return fmt.Errorf("missing .runner config at: %s", runnerPath)
	}

	// Validate _temp/_github_workflow/event.json
	// eventPath := filepath.Join(basePath, "_temp/_github_workflow/event.json")
	// if _, err := os.Stat(eventPath); os.IsNotExist(err) {
	// 	return fmt.Errorf("missing event.json at: %s", eventPath)
	// }

	return nil
}
