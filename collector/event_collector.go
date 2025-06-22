package collector

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/config"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
)

type EventCollector struct {
	eventTriggered *prometheus.CounterVec
	eventTimestamp *prometheus.GaugeVec

	eventPath string
	mode      string
	lastPush  string
}

func NewEventCollector(cfg *config.Config) *EventCollector {
	var eventPath string
	switch cfg.Mode {
	case "test":
		eventPath = cfg.Test.EventPath
	default:
		switch runtime.GOOS {
		case "linux":
			eventPath = cfg.Paths.Logs.Linux.Event
		case "windows":
			eventPath = cfg.Paths.Logs.Windows.Event
		case "darwin":
			eventPath = cfg.Paths.Logs.Mac.Event
		}
	}

	return &EventCollector{
		eventPath: eventPath,
		mode:      cfg.Mode,
		lastPush:  "",
		eventTriggered: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "github_event_triggered_total",
			Help: "Number of GitHub workflow events triggered",
		}, []string{"repo", "org", "enterprise", "workflow", "run_id"}),

		eventTimestamp: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_event_triggered_timestamp_seconds",
			Help: "Unix timestamp of the last GitHub workflow trigger",
		}, []string{"repo", "org", "enterprise", "workflow", "run_id"}),
	}
}

func (c *EventCollector) Describe(ch chan<- *prometheus.Desc) {
	c.eventTriggered.Describe(ch)
	c.eventTimestamp.Describe(ch)
}

func (c *EventCollector) Collect(ch chan<- prometheus.Metric) {
	if _, err := os.Stat(c.eventPath); os.IsNotExist(err) {
		log.Printf("⚠️  event.json not found: %s", c.eventPath)

		c.emitDefault("none", "none", "none", "none", "unknown", ch)
		return
	}

	event, err := parser.ReadEventJSON(c.eventPath)
	if err != nil {
		log.Printf("❌ Failed to parse event.json: %v", err)
		c.emitDefault("none", "none", "none", "none", "unknown", ch)
		return
	}

	repo := event.Repository.RepoName
	org := "none"
	if event.Organization != nil {
		org = event.Organization.OrgName
	}
	ent := "none"
	if event.Enterprise != nil {
		ent = event.Enterprise.Slug
	}
	workflow := event.WorkflowName

	// ✅ Grab RunID from latest Worker log
	logDir := filepath.Dir(c.eventPath)
	fmt.Println("logggg path ", logDir)
	job, err := parser.ParseLatestWorkerLog(logDir)
	runID := "unknown"
	if err == nil && job != nil && job.RunID != "" {
		runID = job.RunID
	}

	ts, err := time.Parse(time.RFC3339, event.Repository.PushedAt)
	if err != nil {
		log.Printf("⚠️  Failed to parse pushed_at time: %v", err)
		return
	}

	if event.Repository.PushedAt != c.lastPush {
		log.Printf("📊 New GitHub event: repo=%s org=%s enterprise=%s workflow=%s run_id=%s time=%s",
			repo, org, ent, workflow, runID, ts)

		labels := []string{repo, org, ent, workflow, runID}
		c.eventTriggered.WithLabelValues(labels...).Inc()
		c.eventTimestamp.WithLabelValues(labels...).Set(float64(ts.Unix()))
		c.lastPush = event.Repository.PushedAt
	}

	c.eventTriggered.Collect(ch)
	c.eventTimestamp.Collect(ch)
}

func (c *EventCollector) emitDefault(repo, org, ent, workflow, runID string, ch chan<- prometheus.Metric) {
	labels := []string{repo, org, ent, workflow, runID}
	c.eventTriggered.WithLabelValues(labels...).Add(0)
	c.eventTimestamp.WithLabelValues(labels...).Set(0)
	c.eventTriggered.Collect(ch)
	c.eventTimestamp.Collect(ch)
}
