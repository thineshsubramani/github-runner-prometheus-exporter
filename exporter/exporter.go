package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/thineshsubramani/github-runner-prometheus-exporter/collector"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/config"
)

type Exporter struct {
	Registry *prometheus.Registry
}

func New(cfg *config.Config) *Exporter {
	reg := prometheus.NewRegistry()

	// reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Register static info
	reg.MustRegister(collector.NewInfoCollector(cfg))
	reg.MustRegister(collector.NewDiskCollector())

	// Custom collectors
	// if cfg.Metrics.EnableRunner {
	// 	reg.MustRegister(collector.NewRunnerCollector(cfg.Paths.Logs.Linux.Worker)) // OS switch handled later
	// }
	// if cfg.Metrics.EnableJob {
	// 	reg.MustRegister(collector.NewJobCollector(cfg.Paths.Logs.Linux.Worker))
	// }
	if cfg.Metrics.EnableEvent {
		reg.MustRegister(collector.NewEventCollector(cfg))
	}

	return &Exporter{Registry: reg}
}
