package collector

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/platform"
)

type RunnerCollector struct {
	actionRunnerPath string

	isActive prometheus.Gauge
	os       *prometheus.GaugeVec
	group    *prometheus.GaugeVec
	name     *prometheus.GaugeVec
	repo     *prometheus.GaugeVec
}

func NewRunnerCollector(path string) prometheus.Collector {

	return &RunnerCollector{
		actionRunnerPath: path,
		isActive: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "runner_status",
			Help: "Shows if the runner is running (1 = yes)",
		}),
		os: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "runner_os",
			Help: "OS of the runner",
		}, []string{"os"}),
		group: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "runner_group",
			Help: "Runner group",
		}, []string{"group"}),
		name: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "runner_name",
			Help: "Runner name",
		}, []string{"name"}),
		repo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "repo_name",
			Help: "Repository name",
		}, []string{"name"}),
	}
}

func (rc *RunnerCollector) Describe(ch chan<- *prometheus.Desc) {
	rc.isActive.Describe(ch)
	rc.os.Describe(ch)
	rc.group.Describe(ch)
	rc.name.Describe(ch)
	rc.repo.Describe(ch)
}

func (rc *RunnerCollector) Collect(ch chan<- prometheus.Metric) {
	if !platform.IsRunnerProcessRunning("Runner.Worker") {
		rc.isActive.Set(0)
		ch <- rc.isActive
		return
	}

	rc.isActive.Set(1)
	ch <- rc.isActive

	runner, err := parser.ReadRunnerConfig(rc.actionRunnerPath)
	if err == nil {
		rc.os.WithLabelValues(runtime.GOOS).Set(1)
		rc.group.WithLabelValues(runner.RunnerGroup).Set(1)
		rc.name.WithLabelValues(runner.RunnerName).Set(1)
	}

	eventPath := rc.actionRunnerPath + "/_temp/_github_workflow/event.json"
	event, err := parser.ReadEventJSON(eventPath)
	if err == nil {
		rc.repo.WithLabelValues(event.Repository.RepoName).Set(1)
	}

	rc.os.Collect(ch)
	rc.group.Collect(ch)
	rc.name.Collect(ch)
	rc.repo.Collect(ch)
}
