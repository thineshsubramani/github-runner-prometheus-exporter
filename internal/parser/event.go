// This if for parse event.json file in runner directory
// Disable by default, Enable this in YAML Configuration if needed
package parser

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// EventInfo holds parsed GitHub event data
type EventInfo struct {
	WorkflowName string `json:"workflow"`

	Repository struct {
		RepoName     string `json:"name"`
		RepoFullName string `json:"full_name"`
		PushedAt     string `json:"pushed_at"` // RFC3339 format
	} `json:"repository"`

	Organization *struct {
		OrgName string `json:"login"`
	} `json:"organization,omitempty"`

	Enterprise *struct {
		Slug string `json:"slug"`
	} `json:"enterprise,omitempty"`
}

func ReadEventJSON(path string) (*EventInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var event EventInfo
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(event, "", "  ")
	log.Println("✅ Parsed Event JSON:")
	log.Println(string(j))

	return &event, nil
}

// GetPushedAtUnix returns the pushed_at field as a Unix timestamp
func (e *EventInfo) GetPushedAtUnix() (int64, error) {
	t, err := time.Parse(time.RFC3339, e.Repository.PushedAt)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
