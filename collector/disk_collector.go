package collector

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskCollector struct {
	desc *prometheus.Desc
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{
		desc: prometheus.NewDesc(
			"disk_usage_bytes",
			"Disk usage for key mountpoints and total",
			[]string{"mount", "type"}, // type: total, used, free, used_percent
			nil,
		),
	}
}

func (c *DiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *DiskCollector) Collect(ch chan<- prometheus.Metric) {
	var mounts []string

	switch runtime.GOOS {
	case "windows":
		mounts = []string{"C:\\", "D:\\"}
	default:
		mounts = []string{"/", "/tmp"}
	}

	var totalTotal, totalUsed, totalFree uint64

	for _, m := range mounts {
		usage, err := disk.Usage(m)
		if err != nil {
			continue
		}

		totalTotal += usage.Total
		totalUsed += usage.Used
		totalFree += usage.Free

		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(usage.Total), m, "total")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(usage.Used), m, "used")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(usage.Free), m, "free")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, usage.UsedPercent, m, "used_percent")
	}

	// Add system-wide total (sum of selected mounts)
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(totalTotal), "all", "total")
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(totalUsed), "all", "used")
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(totalFree), "all", "free")
	if totalTotal > 0 {
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, (float64(totalUsed)/float64(totalTotal))*100, "all", "used_percent")
	}
}
