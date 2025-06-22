package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
)

type RunnerInfo struct {
	RunnerName  string `json:"agentName"`
	RunnerGroup string `json:"poolName"`
	WorkFolder  string `json:"workFolder"`
}

func ReadRunnerConfig(basePath string) (*RunnerInfo, error) {
	runnerPath := filepath.Join(basePath, ".runner")
	data, err := os.ReadFile(runnerPath)
	if err != nil {
		return nil, err
	}
	cleaned := cleanJSON(data)
	var info RunnerInfo
	if err := json.Unmarshal(cleaned, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func cleanJSON(input []byte) []byte {
	re := regexp.MustCompile(`[^\x20-\x7E\n\r\t]`)
	return re.ReplaceAll(input, []byte{})
}
