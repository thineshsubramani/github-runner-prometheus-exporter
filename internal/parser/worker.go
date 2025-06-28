// package parser

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"strings"
// 	"time"
// )

// type WorkerTimestamps struct {
// 	LogFile      string        `json:"log_file"`
// 	StartTime    time.Time     `json:"start_time"`
// 	EndTime      time.Time     `json:"end_time"`
// 	TotalRuntime time.Duration `json:"duration"`
// 	RunID        string        `json:"run_id"`
// }

// // extractTimestamp parses [YYYY-MM-DD HH:MM:SSZ into time.Time
// func extractTimestamp(line string) (time.Time, error) {
// 	start := strings.Index(line, "[")
// 	end := strings.Index(line, "Z")
// 	if start == -1 || end == -1 || end <= start+1 {
// 		return time.Time{}, fmt.Errorf("timestamp not found in: %s", line)
// 	}

// 	raw := line[start+1 : end] // e.g. 2025-05-20 04:14:12
// 	parsed, err := time.Parse("2006-01-02 15:04:05", raw)
// 	if err != nil {
// 		return time.Time{}, fmt.Errorf("failed to parse timestamp: %v", err)
// 	}
// 	return parsed.UTC(), nil
// }

// func ParseLatestWorkerLog(dir string) (*WorkerTimestamps, error) {
// 	logs, err := filepath.Glob(filepath.Join(dir, "Worker_*.log"))
// 	if err != nil || len(logs) == 0 {
// 		return nil, fmt.Errorf("❌ no Worker_*.log found in: %s", dir)
// 	}
// 	sort.Slice(logs, func(i, j int) bool {
// 		return logs[i] > logs[j]
// 	})
// 	latestLog := logs[0]
// 	log.Printf("📄 Latest log selected: %s", latestLog)

// 	file, err := os.Open(latestLog)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	// Read first line
// 	var firstLine string
// 	for scanner.Scan() {
// 		firstLine = scanner.Text()
// 		if strings.TrimSpace(firstLine) != "" {
// 			break
// 		}
// 	}
// 	log.Printf("🔍 First line: %s", firstLine)

// 	// Read last line
// 	var lastLine string
// 	for scanner.Scan() {
// 		text := scanner.Text()
// 		if strings.TrimSpace(text) != "" {
// 			lastLine = text
// 		}
// 	}
// 	log.Printf("🔍 Last line : %s", lastLine)

// 	startTS, err1 := extractTimestamp(firstLine)
// 	endTS, err2 := extractTimestamp(lastLine)

// 	if err1 != nil || err2 != nil {
// 		if err1 != nil {
// 			log.Printf("🧨 Failed to parse start timestamp: %v", err1)
// 		}
// 		if err2 != nil {
// 			log.Printf("🧨 Failed to parse end timestamp: %v", err2)
// 		}
// 		return nil, fmt.Errorf("❌ could not parse timestamps from log")
// 	}

// 	duration := endTS.Sub(startTS)
// 	log.Printf("✅ Parsed start: %s", startTS.Format(time.RFC3339))
// 	log.Printf("✅ Parsed end  : %s", endTS.Format(time.RFC3339))
// 	log.Printf("🕒 Duration     : %s", duration)

// 	runInfo, _ := ExtractRunAndWorkerIDFromLog(latestLog)
// 	runID := "unknown"
// 	if runInfo != nil && runInfo.RunID != "" {
// 		runID = runInfo.RunID
// 	}
// 	fmt.Println(runID)

// 	return &WorkerTimestamps{
// 		LogFile:      filepath.Base(latestLog),
// 		StartTime:    startTS,
// 		EndTime:      endTS,
// 		TotalRuntime: duration,
// 		RunID:        runID,
// 	}, nil
// }

// type RunWorkerInfo struct {
// 	RunID string `json:"run_id"`
// }

// type KeyVal struct {
// 	K string `json:"k"`
// 	V string `json:"v"`
// }

// func ExtractRunAndWorkerIDFromLog(logPath string) (*RunWorkerInfo, error) {
// 	file, err := os.Open(logPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open log file: %w", err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	foundKey := false

// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())

// 		if foundKey {
// 			if strings.Contains(line, `"v":`) {
// 				parts := strings.Split(line, ":")
// 				if len(parts) < 2 {
// 					return nil, fmt.Errorf("invalid run_id value line: %s", line)
// 				}
// 				// Strip quotes and comma if exists
// 				value := strings.Trim(parts[1], `" ,`)
// 				return &RunWorkerInfo{RunID: value}, nil
// 			}
// 		}

// 		if strings.Contains(line, `"k": "run_id"`) {
// 			foundKey = true
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return nil, fmt.Errorf("scanner error: %w", err)
// 	}

//		return nil, fmt.Errorf("run_id not found in log")
//	}
//
// This is custom Worker log parsing logics
// Will add more advanced parsing capability soon
package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type WorkerTimestamps struct {
	LogFile      string        `json:"log_file"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	TotalRuntime time.Duration `json:"duration"`
	RunID        string        `json:"run_id"`
	Slug         string        `json:"slug"`
	Repo         string        `json:"repository"`
	Owner        string        `json:"repository_owner"`
	Workflow     string        `json:"workflow"`
}

// extractTimestamp parses [YYYY-MM-DD HH:MM:SSZ into time.Time
func extractTimestamp(line string) (time.Time, error) {
	start := strings.Index(line, "[")
	end := strings.Index(line, "Z")
	if start == -1 || end == -1 || end <= start+1 {
		return time.Time{}, fmt.Errorf("timestamp not found in: %s", line)
	}

	raw := line[start+1 : end] // e.g. 2025-05-20 04:14:12
	parsed, err := time.Parse("2006-01-02 15:04:05", raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %v", err)
	}
	return parsed.UTC(), nil
}

func ParseLatestWorkerLog(dir string) (*WorkerTimestamps, error) {
	logs, err := filepath.Glob(filepath.Join(dir, "Worker_*.log"))
	if err != nil || len(logs) == 0 {
		return nil, fmt.Errorf("❌ no Worker_*.log found in: %s", dir)
	}
	sort.Slice(logs, func(i, j int) bool {
		return logs[i] > logs[j]
	})
	latestLog := logs[0]
	log.Printf("📄 Latest log selected: %s", latestLog)

	file, err := os.Open(latestLog)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read first line
	var firstLine string
	for scanner.Scan() {
		firstLine = scanner.Text()
		if strings.TrimSpace(firstLine) != "" {
			break
		}
	}
	log.Printf("🔍 First line: %s", firstLine)

	// Read last line
	var lastLine string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) != "" {
			lastLine = text
		}
	}
	log.Printf("🔍 Last line : %s", lastLine)

	startTS, err1 := extractTimestamp(firstLine)
	endTS, err2 := extractTimestamp(lastLine)

	if err1 != nil || err2 != nil {
		if err1 != nil {
			log.Printf("🧨 Failed to parse start timestamp: %v", err1)
		}
		if err2 != nil {
			log.Printf("🧨 Failed to parse end timestamp: %v", err2)
		}
		return nil, fmt.Errorf("❌ could not parse timestamps from log")
	}

	duration := endTS.Sub(startTS)
	log.Printf(" Parsed start: %s", startTS.Format(time.RFC3339))
	log.Printf(" Parsed end  : %s", endTS.Format(time.RFC3339))
	log.Printf(" Duration     : %s", duration)

	runInfo, _ := ExtractJSONFromLog(latestLog)
	runID := "unknown"
	slug := ""
	repo := ""
	owner := ""
	workflow := ""
	if runInfo != nil {
		runID = runInfo.RunID
		slug = runInfo.Slug
		repo = runInfo.Repository
		owner = runInfo.RepositoryOwner
		workflow = runInfo.Workflow
	}
	fmt.Println("RunID:", runID)

	return &WorkerTimestamps{
		LogFile:      filepath.Base(latestLog),
		StartTime:    startTS,
		EndTime:      endTS,
		TotalRuntime: duration,
		RunID:        runID,
		Slug:         slug,
		Repo:         repo,
		Owner:        owner,
		Workflow:     workflow,
	}, nil
}

type RunWorkerInfo struct {
	RunID           string `json:"run_id"`
	Slug            string `json:"slug"`
	Repository      string `json:"repository"`
	RepositoryOwner string `json:"repository_owner"`
	Workflow        string `json:"workflow"`
}

func ExtractJSONFromLog(logPath string) (*RunWorkerInfo, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentKey string
	info := &RunWorkerInfo{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, `"k":`) {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				continue
			}
			currentKey = strings.Trim(parts[1], `" ,`)
		} else if strings.HasPrefix(line, `"v":`) && currentKey != "" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				continue
			}
			value := strings.Trim(parts[1], `" ,{}"`)

			if value == "" {
				value = "unknown"
			}

			switch currentKey {
			case "run_id":
				if info.RunID == "" {
					info.RunID = value
				}
			case "slug":
				if info.Slug == "" {
					info.Slug = value
				}
			case "repository":
				if info.Repository == "" {
					info.Repository = value
				}
			case "repository_owner":
				if info.RepositoryOwner == "" {
					info.RepositoryOwner = value
				}
			case "workflow":
				if info.Workflow == "" {
					info.Workflow = value
				}
			}

			currentKey = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	if info.RunID == "" {
		return nil, fmt.Errorf("run_id not found in log")
	}

	// Fallback defaults
	if info.Slug == "" {
		info.Slug = "unknown"
	}
	if info.Repository == "" {
		info.Repository = "unknown"
	}
	if info.RepositoryOwner == "" {
		info.RepositoryOwner = "unknown"
	}
	if info.Workflow == "" {
		info.Workflow = "unknown"
	}

	return info, nil
}

// TODO: Use for live log parsing line by line
