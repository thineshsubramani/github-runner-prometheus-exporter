// This collector, Collect static informations such as OS, Runner Name, Runner Group (From YAML)
package collector

import (
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

	return &InfoCollector{
		mode:        cfg.Mode,
		runnerNames: cfg.Runners.Names,
		groupNames:  cfg.Runners.Groups,
		os:          runtime.GOOS,
		infoDesc: prometheus.NewDesc(
			"github_runner_static_info",
			"Static config info: mode, runners, groups, os",
			[]string{"mode", "runner_names", "group_names", "os"},
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
		c.os,
	)
}
