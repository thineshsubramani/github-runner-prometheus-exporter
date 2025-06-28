package platform

import (
	"log"
	"os"
	"runtime"

	"github.com/thineshsubramani/github-runner-prometheus-exporter/config"
)

func GetOS() string {
	return runtime.GOOS // "linux", "windows", "darwin"
}

// We are still pulling info from config yaml like static infos
// TODO: more dynamic way to pull server metadata for labels
// EG. Server Pool, OS Version
//

func DefaultPath(cfg *config.Config) string {
	runnerName := os.Getenv("RUNNER_NAME")
	if runnerName == "" {
		log.Println("⚠️  RUNNER_NAME env var not set, defaulting to first runner in config")
		if len(cfg.Runners) == 0 {
			log.Fatal("❌ no runners defined in config")
		}
		runnerName = cfg.Runners[0].Name
	}

	for _, runner := range cfg.Runners {
		if runner.Name == runnerName {
			// prefer config path if given
			if runner.Mode == "test" && runner.Test.EventPath != "" {
				return runner.Test.EventPath
			}
			if runner.Mode == "prod" && runner.Logs.Event != "" {
				return runner.Logs.Event
			}
			break
		}
	}

	// fallback defaults based on OS
	switch runtime.GOOS {
	case "linux":
		return "/var/log/github-runner/default-event.json"
	case "windows":
		return "C:\\github-runner\\default-event.json"
	case "darwin":
		return "/Users/Shared/github-runner/default-event.json"
	default:
		log.Fatal("unsupported OS for default path")
		return ""
	}
}
