// package collector

// import (
// 	"log"
// 	"path/filepath"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
// )

// type JobCollector struct {
// 	logPath         string
// 	workflowStart   *prometheus.GaugeVec
// 	workflowEnd     *prometheus.GaugeVec
// 	workflowRuntime *prometheus.GaugeVec
// }

// func NewJobCollector(path string) *JobCollector {
// 	return &JobCollector{
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

// func (c *JobCollector) Describe(ch chan<- *prometheus.Desc) {
// 	c.workflowStart.Describe(ch)
// 	c.workflowEnd.Describe(ch)
// 	c.workflowRuntime.Describe(ch)
// }

// func (c *JobCollector) Collect(ch chan<- prometheus.Metric) {
// 	job, err := parser.ParseLatestWorkerLog(c.logPath)
// 	if err != nil || job == nil {
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

// 	logLabel := job.LogFile

// 	runInfo, err := parser.ExtractRunAndJobIDFromLog(filepath.Join(c.logPath, logLabel))
// 	if err != nil {
// 		log.Printf("⚠️  Failed to extract RunId: %v", err)
// 		runInfo = &parser.RunJobInfo{RunID: "unknown"}
// 	}

// 	runID := runInfo.RunID
// 	labels := []string{logLabel, runID}

// 	c.workflowStart.WithLabelValues(labels...).Set(float64(job.StartTime.Unix()))
// 	c.workflowEnd.WithLabelValues(labels...).Set(float64(job.EndTime.Unix()))
// 	c.workflowRuntime.WithLabelValues(labels...).Set(job.TotalRuntime.Seconds())

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

// type JobCollector struct {
// 	logPath         string
// 	workflowStart   *prometheus.GaugeVec
// 	workflowEnd     *prometheus.GaugeVec
// 	workflowRuntime *prometheus.GaugeVec
// }

// func NewJobCollector(path string) *JobCollector {
// 	labelKeys := []string{
// 		"log_file",
// 		"run_id",
// 		"slug",
// 		"repository",
// 		"repository_owner",
// 		"workflow",
// 	}

// 	return &JobCollector{
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

// func (c *JobCollector) Describe(ch chan<- *prometheus.Desc) {
// 	c.workflowStart.Describe(ch)
// 	c.workflowEnd.Describe(ch)
// 	c.workflowRuntime.Describe(ch)
// }

// func (c *JobCollector) Collect(ch chan<- prometheus.Metric) {
// 	job, err := parser.ParseLatestWorkerLog(c.logPath)
// 	if err != nil || job == nil || job.RunID == "" {
// 		log.Printf("⚠️  Failed to parse job or missing run_id: %v", err)
// 		return
// 	}

// 	labels := []string{
// 		defaultIfEmpty(job.LogFile),
// 		defaultIfEmpty(job.RunID),
// 		defaultIfEmpty(job.Slug),
// 		defaultIfEmpty(job.Repo),
// 		defaultIfEmpty(job.Owner),
// 		defaultIfEmpty(job.Workflow),
// 	}

// 	log.Printf("📌 Labels: %#v", labels)
// 	log.Printf("✅ StartTime: %v (%d)", job.StartTime, job.StartTime.Unix())
// 	log.Printf("✅ EndTime  : %v (%d)", job.EndTime, job.EndTime.Unix())
// 	log.Printf("✅ Duration : %v (%.0f seconds)", job.TotalRuntime, job.TotalRuntime.Seconds())

// 	c.workflowStart.WithLabelValues(labels...).Set(float64(job.StartTime.Unix()))
// 	c.workflowEnd.WithLabelValues(labels...).Set(float64(job.EndTime.Unix()))
// 	c.workflowRuntime.WithLabelValues(labels...).Set(job.TotalRuntime.Seconds())

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

type JobCollector struct {
	logPath         string
	workflowStart   *prometheus.GaugeVec
	workflowEnd     *prometheus.GaugeVec
	workflowRuntime *prometheus.GaugeVec
	runnerState     *prometheus.GaugeVec
}

func NewJobCollector(path string) *JobCollector {
	labelKeys := []string{
		"log_file",
		"run_id",
		"slug",
		"repository",
		"repository_owner",
		"workflow",
	}

	c := &JobCollector{
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

func (c *JobCollector) Describe(ch chan<- *prometheus.Desc) {
	c.workflowStart.Describe(ch)
	c.workflowEnd.Describe(ch)
	c.workflowRuntime.Describe(ch)
	c.runnerState.Describe(ch)
}

func (c *JobCollector) Collect(ch chan<- prometheus.Metric) {
	job, err := parser.ParseLatestWorkerLog(c.logPath)
	if err != nil || job == nil || job.RunID == "" {
		log.Printf("Failed to parse job or missing run_id: %v", err)
		return
	}

	labels := []string{
		defaultIfEmpty(job.LogFile),
		defaultIfEmpty(job.RunID),
		defaultIfEmpty(job.Slug),
		defaultIfEmpty(job.Repo),
		defaultIfEmpty(job.Owner),
		defaultIfEmpty(job.Workflow),
	}

	log.Printf("Labels: %#v", labels)
	log.Printf("StartTime: %v (%d)", job.StartTime, job.StartTime.Unix())
	log.Printf("EndTime  : %v (%d)", job.EndTime, job.EndTime.Unix())
	log.Printf("Duration : %v (%.0f seconds)", job.TotalRuntime, job.TotalRuntime.Seconds())

	c.workflowStart.WithLabelValues(labels...).Set(float64(job.StartTime.Unix()))
	c.workflowEnd.WithLabelValues(labels...).Set(float64(job.EndTime.Unix()))
	c.workflowRuntime.WithLabelValues(labels...).Set(job.TotalRuntime.Seconds())

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
