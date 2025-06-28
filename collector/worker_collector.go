// package collector

// import (
// 	"log"
// 	"path/filepath"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
// )

// type WorkerCollector struct {
// 	logPath         string
// 	workflowStart   *prometheus.GaugeVec
// 	workflowEnd     *prometheus.GaugeVec
// 	workflowRuntime *prometheus.GaugeVec
// }

// func NewWorkerCollector(path string) *WorkerCollector {
// 	return &WorkerCollector{
// 		logPath: path,
// 		workflowStart: prometheus.NewGaugeVec(prometheus.GaugeOpts{
// 			Name: "github_workflow_start_timestamp_seconds",
// 			Help: "Start time of latest GitHub workflow log",
// 		}, []string{"log_file", "run_id"}),

// 		workflowEnd: prometheus.NewGaugeVec(prometheus.GaugeOpts{
// 			Name: "github_workflow_end_timestamp_seconds",
// 			Help: "End time of latest GitHub workflow log",
// 		}, []string{"log_file", "run_id"}),

// 		workflowRuntime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
// 			Name: "github_workflow_duration_seconds",
// 			Help: "Total duration of latest GitHub workflow log",
// 		}, []string{"log_file", "run_id"}),
// 	}
// }

// func (c *WorkerCollector) Describe(ch chan<- *prometheus.Desc) {
// 	c.workflowStart.Describe(ch)
// 	c.workflowEnd.Describe(ch)
// 	c.workflowRuntime.Describe(ch)
// }

// func (c *WorkerCollector) Collect(ch chan<- prometheus.Metric) {
// 	Worker, err := parser.ParseLatestWorkerLog(c.logPath)
// 	if err != nil || Worker == nil {
// 		log.Printf("⚠️  Worker log not found or parse failed: %v", err)

// 		// Emit idle placeholders with static "none" log label
// 		labels := []string{"none", "unknown"}

// 		c.workflowStart.WithLabelValues(labels...).Set(0)
// 		c.workflowEnd.WithLabelValues(labels...).Set(0)
// 		c.workflowRuntime.WithLabelValues(labels...).Set(0)

// 		c.workflowStart.Collect(ch)
// 		c.workflowEnd.Collect(ch)
// 		c.workflowRuntime.Collect(ch)
// 		return
// 	}

// 	logLabel := Worker.LogFile

// 	runInfo, err := parser.ExtractRunAndWorkerIDFromLog(filepath.Join(c.logPath, logLabel))
// 	if err != nil {
// 		log.Printf("⚠️  Failed to extract RunId: %v", err)
// 		runInfo = &parser.RunWorkerInfo{RunID: "unknown"}
// 	}

// 	runID := runInfo.RunID
// 	labels := []string{logLabel, runID}

// 	c.workflowStart.WithLabelValues(labels...).Set(float64(Worker.StartTime.Unix()))
// 	c.workflowEnd.WithLabelValues(labels...).Set(float64(Worker.EndTime.Unix()))
// 	c.workflowRuntime.WithLabelValues(labels...).Set(Worker.TotalRuntime.Seconds())

//		c.workflowStart.Collect(ch)
//		c.workflowEnd.Collect(ch)
//		c.workflowRuntime.Collect(ch)
//	}
// VERSION @
// package collector

// import (
// 	"log"
// 	"strings"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
// )

// type WorkerCollector struct {
// 	logPath         string
// 	workflowStart   *prometheus.GaugeVec
// 	workflowEnd     *prometheus.GaugeVec
// 	workflowRuntime *prometheus.GaugeVec
// }

// func NewWorkerCollector(path string) *WorkerCollector {
// 	labelKeys := []string{
// 		"log_file",
// 		"run_id",
// 		"slug",
// 		"repository",
// 		"repository_owner",
// 		"workflow",
// 	}

// 	return &WorkerCollector{
// 		logPath: path,
// 		workflowStart: prometheus.NewGaugeVec(prometheus.GaugeOpts{
// 			Name: "github_workflow_start_timestamp_seconds",
// 			Help: "Start time of GitHub workflow run",
// 		}, labelKeys),

// 		workflowEnd: prometheus.NewGaugeVec(prometheus.GaugeOpts{
// 			Name: "github_workflow_end_timestamp_seconds",
// 			Help: "End time of GitHub workflow run",
// 		}, labelKeys),

// 		workflowRuntime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
// 			Name: "github_workflow_duration_seconds",
// 			Help: "Duration of GitHub workflow run",
// 		}, labelKeys),
// 	}
// }

// func (c *WorkerCollector) Describe(ch chan<- *prometheus.Desc) {
// 	c.workflowStart.Describe(ch)
// 	c.workflowEnd.Describe(ch)
// 	c.workflowRuntime.Describe(ch)
// }

// func (c *WorkerCollector) Collect(ch chan<- prometheus.Metric) {
// 	Worker, err := parser.ParseLatestWorkerLog(c.logPath)
// 	if err != nil || Worker == nil || Worker.RunID == "" {
// 		log.Printf("⚠️  Failed to parse Worker or missing run_id: %v", err)
// 		return
// 	}

// 	labels := []string{
// 		defaultIfEmpty(Worker.LogFile),
// 		defaultIfEmpty(Worker.RunID),
// 		defaultIfEmpty(Worker.Slug),
// 		defaultIfEmpty(Worker.Repo),
// 		defaultIfEmpty(Worker.Owner),
// 		defaultIfEmpty(Worker.Workflow),
// 	}

// 	log.Printf("📌 Labels: %#v", labels)
// 	log.Printf("✅ StartTime: %v (%d)", Worker.StartTime, Worker.StartTime.Unix())
// 	log.Printf("✅ EndTime  : %v (%d)", Worker.EndTime, Worker.EndTime.Unix())
// 	log.Printf("✅ Duration : %v (%.0f seconds)", Worker.TotalRuntime, Worker.TotalRuntime.Seconds())

// 	c.workflowStart.WithLabelValues(labels...).Set(float64(Worker.StartTime.Unix()))
// 	c.workflowEnd.WithLabelValues(labels...).Set(float64(Worker.EndTime.Unix()))
// 	c.workflowRuntime.WithLabelValues(labels...).Set(Worker.TotalRuntime.Seconds())

// 	c.workflowStart.Collect(ch)
// 	c.workflowEnd.Collect(ch)
// 	c.workflowRuntime.Collect(ch)
// }

// func defaultIfEmpty(s string) string {
// 	if strings.TrimSpace(s) == "" {
// 		return "unknown"
// 	}
// 	return strings.Trim(s, "{}\" ")
// }

package collector

import (
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/watcher"
)

type WorkerCollector struct {
	logPath         string
	workflowStart   *prometheus.GaugeVec
	workflowEnd     *prometheus.GaugeVec
	workflowRuntime *prometheus.GaugeVec
	runnerState     *prometheus.GaugeVec
}

func NewWorkerCollector(path string) *WorkerCollector {
	labelKeys := []string{
		"log_file",
		"run_id",
		"slug",
		"repository",
		"repository_owner",
		"workflow",
	}

	c := &WorkerCollector{
		logPath: path,
		workflowStart: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_workflow_start_timestamp_seconds",
			Help: "Start time of GitHub workflow run",
		}, labelKeys),

		workflowEnd: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_workflow_end_timestamp_seconds",
			Help: "End time of GitHub workflow run",
		}, labelKeys),

		workflowRuntime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_workflow_duration_seconds",
			Help: "Duration of GitHub workflow run",
		}, labelKeys),

		runnerState: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_runner_state",
			Help: "Runner state: 1=busy, 0=idle",
		}, []string{"hostname", "mode"}),
	}

	// Start watcher inline
	go func() {
		err := watcher.WatchLogDir(path, func(path string, event string) {
			switch event {
			case "created":
				log.Println("Runner active (event.json created)")
				c.runnerState.WithLabelValues("insight-development-lab", "prod").Set(1)
			case "deleted":
				log.Println("Runner idle (event.json deleted)")
				c.runnerState.WithLabelValues("insight-development-lab", "prod").Set(0)
			}
		})
		if err != nil {
			log.Printf(" Watcher error: %v", err)
		}
	}()

	return c
}

func (c *WorkerCollector) Describe(ch chan<- *prometheus.Desc) {
	c.workflowStart.Describe(ch)
	c.workflowEnd.Describe(ch)
	c.workflowRuntime.Describe(ch)
	c.runnerState.Describe(ch)
}

func (c *WorkerCollector) Collect(ch chan<- prometheus.Metric) {
	Worker, err := parser.ParseLatestWorkerLog(c.logPath)
	if err != nil || Worker == nil || Worker.RunID == "" {
		log.Printf("Failed to parse Worker or missing run_id: %v", err)
		return
	}

	labels := []string{
		defaultIfEmpty(Worker.LogFile),
		defaultIfEmpty(Worker.RunID),
		defaultIfEmpty(Worker.Slug),
		defaultIfEmpty(Worker.Repo),
		defaultIfEmpty(Worker.Owner),
		defaultIfEmpty(Worker.Workflow),
	}

	log.Printf("Labels: %#v", labels)
	log.Printf("StartTime: %v (%d)", Worker.StartTime, Worker.StartTime.Unix())
	log.Printf("EndTime  : %v (%d)", Worker.EndTime, Worker.EndTime.Unix())
	log.Printf("Duration : %v (%.0f seconds)", Worker.TotalRuntime, Worker.TotalRuntime.Seconds())

	c.workflowStart.WithLabelValues(labels...).Set(float64(Worker.StartTime.Unix()))
	c.workflowEnd.WithLabelValues(labels...).Set(float64(Worker.EndTime.Unix()))
	c.workflowRuntime.WithLabelValues(labels...).Set(Worker.TotalRuntime.Seconds())

	c.workflowStart.Collect(ch)
	c.workflowEnd.Collect(ch)
	c.workflowRuntime.Collect(ch)
	c.runnerState.Collect(ch)
}

func defaultIfEmpty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "unknown"
	}
	return strings.Trim(s, "{}\" ")
}
