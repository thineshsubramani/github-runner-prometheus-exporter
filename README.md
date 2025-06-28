# GitHub Runner Prometheus Exporter

**Work in Progress**: This project is under active development. Some features may be incomplete or subject to change. Please wait for the official release.

## Project Goal

The GitHub Runner Prometheus Exporter provides a solution for monitoring self-hosted GitHub runners at scale. The default GitHub API lacks efficiency for real-time, fine-grained metrics across multiple runners. This exporter runs locally on the same machine as the GitHub runner, scrapes logs, and exposes detailed Prometheus metrics for real-time observability. It is designed to work in diverse environments, including air-gapped or restricted setups, and supports bare metal, on-premises, or cloud-based runners (e.g., GCP, AWS EC2).

## Purpose

- Overcome limitations of the GitHub API for monitoring multiple self-hosted runners.
- Provide log-based Prometheus metrics for workflow execution (e.g., start time, end time, duration, status).
- Include system-level metrics (e.g., disk usage, CPU, memory, network) for comprehensive runner health monitoring.
- Enable full control over runner metrics in production environments.

## How It Works

The exporter is a Go-based service that:
1. Monitors the GitHub runner's log directory (e.g., `worker.log`, `runner.log`).
2. Parses logs to extract workflow and system metrics.
3. Exposes metrics in Prometheus format via an HTTP endpoint (`/metrics`).

Example metric:
```
# HELP github_workflow_duration_seconds Duration of GitHub workflow run
# TYPE github_workflow_duration_seconds gauge
github_workflow_duration_seconds{runner_name="gpu-runner",status="success"} 142.5
```

## Metrics

The exporter provides the following metrics, with additional system and process metrics for enhanced observability:

| Metric                                    | Labels                                                  | Description                                              |
|-------------------------------------------|---------------------------------------------------------|----------------------------------------------------------|
| `github_runner_static_info`               | `runner_names`, `group_names`, `hostname`, `os`, `mode` | Static configuration details of the runner               |
| `github_workflow_start_timestamp_seconds` | `run_id`, `job_name`, `runner_name`, `runner_group`, `org`, `repo`, `workflow` | Start time of a workflow run                             |
| `github_workflow_end_timestamp_seconds`   | `run_id`, `job_name`, `runner_name`, `runner_group`, `org`, `repo`, `workflow` | End time of a workflow run                               |
| `github_workflow_duration_seconds`        | `run_id`, `job_name`, `runner_name`, `runner_group`, `org`, `repo`, `workflow`, `status` | Duration of a workflow run                               |
| `github_event_triggered_total`            | `hostname`, `org`, `os`, `repo`, `runner_group`, `runner_name`, `workflow` | Total number of workflow events triggered                |
| `github_runner_state`                     | `hostname`, `os`, `runner_group`, `runner_name`, `state` | Runner state (1=busy, 0=idle)                           |
| `disk_usage_bytes`                        | `hostname`, `mount`, `os`, `runner_group`, `runner_name`, `type` | Disk usage for key mountpoints (`free`, `used`, `total`, `used_percent`) |
| `process_cpu_seconds_total`               | `hostname`, `os`, `runner_group`, `runner_name`         | Total user and system CPU time spent in seconds          |
| `process_max_fds`                         | `hostname`, `os`, `runner_group`, `runner_name`         | Maximum number of open file descriptors                  |
| `process_network_receive_bytes_total`     | `hostname`, `os`, `runner_group`, `runner_name`         | Total bytes received by the process over the network     |
| `process_network_transmit_bytes_total`    | `hostname`, `os`, `runner_group`, `runner_name`         | Total bytes sent by the process over the network         |
| `process_open_fds`                        | `hostname`, `os`, `runner_group`, `runner_name`         | Number of open file descriptors                         |
| `process_resident_memory_bytes`           | `hostname`, `os`, `runner_group`, `runner_name`         | Resident memory size in bytes                            |
| `process_start_time_seconds`              | `hostname`, `os`, `runner_group`, `runner_name`         | Process start time since Unix epoch in seconds           |
| `process_virtual_memory_bytes`            | `hostname`, `os`, `runner_group`, `runner_name`         | Virtual memory size in bytes                             |
| `process_virtual_memory_max_bytes`        | `hostname`, `os`, `runner_group`, `runner_name`         | Maximum virtual memory available in bytes                |

Example metrics output:
```
# HELP disk_usage_bytes Disk usage for key mountpoints and total
# TYPE disk_usage_bytes gauge
disk_usage_bytes{hostname="insight-development-lab",mount="/",os="linux",runner_group="sre-group",runner_name="gpu-runner",type="free"} 1.51457210368e+11
disk_usage_bytes{hostname="insight-development-lab",mount="/",os="linux",runner_group="sre-group",runner_name="gpu-runner",type="total"} 2.31907807232e+11
disk_usage_bytes{hostname="insight-development-lab",mount="/",os="linux",runner_group="sre-group",runner_name="gpu-runner",type="used"} 6.8595818496e+10
disk_usage_bytes{hostname="insight-development-lab",mount="/",os="linux",runner_group="sre-group",runner_name="gpu-runner",type="used_percent"} 31.172403692927343
# HELP github_event_triggered_total Number of GitHub workflow events triggered
# TYPE github_event_triggered_total counter
github_event_triggered_total{hostname="insight-development-lab",org="insight-dev",os="linux",repo="manage-shared-runners",runner_group="sre-group",runner_name="sharred-runner",workflow=".github/workflows/check-details.yml"} 1
# HELP github_runner_state State of the GitHub runner: 1=busy, 0=idle
# TYPE github_runner_state gauge
github_runner_state{hostname="insight-development-lab",os="linux",runner_group="sre-group",runner_name="gpu-runner",state="busy"} 1
github_runner_state{hostname="insight-development-lab",os="linux",runner_group="sre-group",runner_name="gpu-runner",state="idle"} 0
```

## Requirements

- GitHub Self-Hosted Runners (tested on Linux)
- Prometheus
- Optional: Grafana for visualization (dashboard template in progress)

## Getting Started

### 1. Clone and Build
```bash
git clone https://github.com/thineshsubramani/github-runner-prometheus-exporter
cd github-runner-prometheus-exporter
go build -o exporter cmd/main.go
```

### 2. Run the Exporter
```bash
./exporter --log-path /path/to/github/runner/_diag
```

### 3. Configure Prometheus
Add the following job to your Prometheus configuration:
```yaml
- job_name: 'github-runner-exporter'
  static_configs:
    - targets: ['<runner-host>:8080']
```

## Current Status

- Completed:
  - Log scraping and core metrics extraction
  - Prometheus exporter functionality
  - System and process metrics (disk, CPU, memory, network)
  - Runner state detection (busy/idle)
- In Progress:
  - Workflow grouping logic
  - Grafana dashboard template
- Known Issues:
  - Log parsing may be incomplete or misgrouped in edge cases

## Contributing

Contributions are welcome! Please submit pull requests or open issues for bugs, feature requests, or improvements. Areas for collaboration include:
- Enhanced log parsing logic
- Grafana dashboard development
- Support for multi-runner setups on a single host

File issues or contribute at: [GitHub Issues](https://github.com/thineshsubramani/github-runner-prometheus-exporter/issues)

## Roadmap

- Improve runner state monitoring (avoid polling)
- Add support for log ingestion via Filebeat or Fluent Bit
- Implement multi-runner support with master/child architecture
- Enhance Grafana dashboard for comprehensive visualization

## About

This project was developed to address real-world challenges in monitoring GitHub runners in production environments. It has been deployed in an enterprise setup with over 200 globally distributed runner servers, enabling insights into runner utilization across teams. The project aligns with the author's focus on observability, infrastructure monitoring, and log-to-metrics pipelines as part of a transition from DevOps to SRE.

## Contact

For bugs, feature requests, or contributions, please use [GitHub Issues](https://github.com/thineshsubramani/github-runner-prometheus-exporter/issues).
