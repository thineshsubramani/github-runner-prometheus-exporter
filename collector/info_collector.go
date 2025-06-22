// This collector, Collect static informations such as OS, Runner Name, Runner Group (From YAML)
package collector

import (
	"os"
	"runtime"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/config"
)

type InfoCollector struct {
	mode        string
	runnerNames []string
	groupNames  []string
	hostname    string
	os          string

	infoDesc *prometheus.Desc
}

func NewInfoCollector(cfg *config.Config) prometheus.Collector {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return &InfoCollector{
		mode:        cfg.Mode,
		runnerNames: cfg.Runners.Names,
		groupNames:  cfg.Runners.Groups,
		hostname:    hostname,
		os:          runtime.GOOS,
		infoDesc: prometheus.NewDesc(
			"github_runner_static_info",
			"Static config info: mode, runners, groups, hostname, os",
			[]string{"mode", "runner_names", "group_names", "hostname", "os"},
			nil,
		),
	}
}

func (c *InfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
}

func (c *InfoCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.infoDesc,
		prometheus.GaugeValue,
		1,
		c.mode,
		strings.Join(c.runnerNames, ","),
		strings.Join(c.groupNames, ","),
		c.hostname,
		c.os,
	)
}
