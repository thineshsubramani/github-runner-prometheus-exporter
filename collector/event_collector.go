package collector

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/config"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/platform"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/watcher"
)

type EventCollector struct {
	eventTriggered *prometheus.CounterVec
	eventTimestamp *prometheus.GaugeVec
	runnerState    *prometheus.GaugeVec

	eventPath    string
	lastPush     string
	activeLabels []string
	runnerIdle   bool
	mu           sync.Mutex
}

func NewEventCollector(cfg *config.Config) *EventCollector {
	var eventPath string
	eventPath = platform.DefaultPath(cfg)
	fmt.Println("-----------------------------------", eventPath)

	c := &EventCollector{
		eventPath:  eventPath,
		lastPush:   "",
		runnerIdle: true,
		eventTriggered: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "github_event_triggered_total",
			Help: "Number of GitHub workflow events triggered",
		}, []string{"repo", "org", "workflow"}),

		eventTimestamp: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_event_triggered_timestamp_seconds",
			Help: "Unix timestamp of the last GitHub workflow trigger",
		}, []string{"repo", "org", "workflow"}),

		runnerState: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_runner_state",
			Help: "State of the GitHub runner: 1=busy, 0=idle",
		}, []string{"state"}),
	}

	// Initial state
	if _, err := os.Stat(eventPath); err == nil {
		c.setRunnerState(false)
	} else {
		c.setRunnerState(true)
	}

	// Watch the parent directory of event.json
	eventDir := filepath.Dir(eventPath)
	go func() {
		err := watcher.WatchLogDir(eventDir, func(path string, event string) {
			if filepath.Base(path) != filepath.Base(eventPath) {
				return
			}
			c.mu.Lock()
			defer c.mu.Unlock()

			switch event {
			case "created":
				c.setRunnerState(false)
			case "deleted":
				c.setRunnerState(true)
				c.lastPush = ""
				if c.activeLabels != nil {
					c.eventTriggered.DeleteLabelValues(c.activeLabels...)
					// c.eventTimestamp.DeleteLabelValues(c.activeLabels...)
					c.activeLabels = nil
				}
			}
		})
		if err != nil {
			log.Printf("Watcher error: %v", err)
		}
	}()

	return c
}

func (c *EventCollector) setRunnerState(idle bool) {
	c.runnerIdle = idle
	if idle {
		c.runnerState.WithLabelValues("idle").Set(1)
		c.runnerState.WithLabelValues("busy").Set(0)
	} else {
		c.runnerState.WithLabelValues("idle").Set(0)
		c.runnerState.WithLabelValues("busy").Set(1)
	}
}

func (c *EventCollector) Describe(ch chan<- *prometheus.Desc) {
	c.eventTriggered.Describe(ch)
	c.runnerState.Describe(ch)
}

func (c *EventCollector) Collect(ch chan<- prometheus.Metric) {
	c.runnerState.Collect(ch)

	if c.runnerIdle {
		return
	}

	if _, err := os.Stat(c.eventPath); os.IsNotExist(err) {
		return
	}

	event, err := parser.ReadEventJSON(c.eventPath)
	if err != nil {
		log.Printf("Failed to parse event.json: %v", err)
		return
	}

	repo := event.Repository.RepoName
	org := ""
	if event.Organization != nil {
		org = event.Organization.OrgName
	}
	// ent := ""
	// if event.Enterprise != nil {
	// 	ent = event.Enterprise.Slug
	// }
	workflow := event.WorkflowName

	// logDir := filepath.Dir(c.eventPath)
	// job, err := parser.ParseLatestWorkerLog(logDir)
	// runID := ""
	// if err == nil && job != nil && job.RunID != "" {
	// 	runID = job.RunID
	// }

	// ts, err := time.Parse(time.RFC3339, event.Repository.PushedAt)
	// if err != nil {
	// 	log.Printf("Failed to parse pushed_at: %v", err)
	// 	return
	// }

	labels := []string{repo, org, workflow}

	if event.Repository.PushedAt != c.lastPush {
		c.eventTriggered.WithLabelValues(labels...).Inc()
		c.lastPush = event.Repository.PushedAt

		c.mu.Lock()
		c.activeLabels = labels
		c.mu.Unlock()
	}

	c.eventTriggered.Collect(ch)
	c.eventTimestamp.Collect(ch)
}
