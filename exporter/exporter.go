package exporter

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/thineshsubramani/github-runner-prometheus-exporter/collector"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/config"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/platform"
)

type Exporter struct {
	Registry *prometheus.Registry
}

func New(cfg *config.Config) *Exporter {
	reg := prometheus.NewRegistry()

	hostname, _ := os.Hostname()
	wrappedReg := prometheus.WrapRegistererWith(
		prometheus.Labels{"hostname": hostname, "os": platform.GetOS()},
		reg,
	)

	runner_name := cfg.Runners.Names[0]
	group_name := cfg.Runners.Groups[0]
	runnerWrappedReg := prometheus.WrapRegistererWith(
		prometheus.Labels{"runner_name": runner_name, "runner_group": group_name},
		wrappedReg,
	)
	// reg.MustRegister(collectors.NewGoCollector())

	wrappedReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	// wrappedReg.MustRegister(collector.NewInfoCollector(cfg))
	wrappedReg.MustRegister(collector.NewDiskCollector())
	// Custom collectors
	// if cfg.Metrics.EnableRunner {
	// 	reg.MustRegister(collector.NewRunnerCollector(cfg.Paths.Logs.Linux.Worker)) // OS switch handled later
	// }
	// if cfg.Metrics.EnableJob {
	// 	reg.MustRegister(collector.NewJobCollector(cfg.Paths.Logs.Linux.Worker))
	// }
	if cfg.Metrics.EnableEvent {
		runnerWrappedReg.MustRegister(collector.NewEventCollector(cfg))
	}

	return &Exporter{Registry: reg}
}
