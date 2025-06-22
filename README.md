# GitHub Runner Prometheus Exporter

**🚧 Work in Progress – Some parts may be broken. Please wait for the official release.**

## Project Goal

GitHub's default monitoring exposes runner metrics via API but it **lose at scale** when managing **multiple self-hosted runners** across environments. This exporter is designed to run **on the same machine as the GitHub runner**, **scrape logs locally**, and expose **fine-grained Prometheus metrics** for real-time visibility, even in **air-gapped or restricted setups**.

---

## 🔥 Why This Exists

- **Existing GitHub API is limited**: It's not reliable or efficient for tracking multiple runners with real-time log-based metrics.
- **You own the runners, you should own the metrics.**
- Export **workflow logs** as Prometheus metrics: start time, end time, duration, errors, etc.
- Includes basic **system-level metrics** (disk usage, memory, etc.) for full runner health monitoring.
- Ideal for **bare metal**, **on-prem**, or **cloud** runners (GCP, AWS EC2, etc.).

---

## 🛠 How It Works

This exporter runs as a Go service:

1. Watches the GitHub runner's log directory.
2. Parses logs (`worker.log`, `runner.log`, etc.).
3. Exposes Prometheus-friendly metrics at `/metrics`.

Example metric:
```bash
# HELP github_workflow_duration_seconds Duration of GitHub workflow run
# TYPE github_workflow_duration_seconds gauge
github_workflow_duration_seconds{runner_name="gpu-runner",status="success"} 142.5
````

---

## 📦 Metrics Overview

```bash
# HELP github_runner_static_info Static config info: mode, runners, groups, hostname, os
# HELP github_workflow_duration_seconds Duration of GitHub workflow run
# HELP github_workflow_end_timestamp_seconds End time of GitHub workflow run
# HELP github_workflow_start_timestamp_seconds Start time of GitHub workflow run

```
---

| Metric                                    | Labels                                                  | Description                    |
| ----------------------------------------- | ------------------------------------------------------- | ------------------------------ |
| `github_runner_static_info`               | `runner_names`, `group_names`, `hostname`, `os`, `mode` | Static config about the runner |
| `github_workflow_start_timestamp_seconds` | `run_id`, `job_name`, etc.                              | Start time of a workflow       |
| `github_workflow_end_timestamp_seconds`   | `run_id`, `job_name`, etc.                              | End time of a workflow         |
| `github_workflow_duration_seconds`        | `run_id`, `status`, etc.                                | Duration of the workflow       |
| `disk_usage_bytes`                        | `mount`, `type` (`free`, `used`, `total`)               | System disk usage              |
| (More system metrics coming soon...)      |                                                         |                                |

---

## 🔧 Requirements

* GitHub Self-Hosted Runners (Linux tested)
* Prometheus
* Optional: Grafana dashboard (template WIP)

---

## 🚀 Getting Started

### 1. Clone & Build

```bash
git clone https://github.com/thineshsubramani/github-runner-prometheus-exporter
cd github-runner-prometheus-exporter
go build -o exporter cmd/main.go
```

### 2. Run Exporter

```bash
./exporter --log-path /path/to/github/runner/_diag
```

### 3. Prometheus Config

Add a job like:

```yaml
- job_name: 'github-runner-exporter'
  static_configs:
    - targets: ['<runner-host>:8080']
```

---

## 📍 Current Status

* ✅ Log scraping & basic metrics
* ✅ Prometheus exporter
* 🛠 Workflow grouping WIP
* 🛠 Idle/Active runner state detection coming
* 🛠 Grafana dashboard coming
* ⚠️ Logs might be incomplete or misgrouped in edge cases

---

## 🤝 Contributing
I'm working solo on this project, happy to collaborate if anyone wants to contribute! Pull requests are welcome. If you’ve got better parsing logic or want to help with Grafana dashboards, hit me up.

---

## 👀 Roadmap

* [ ] Better state monitoring of runners (avoid polling)
* [ ] Extension support like filebeat or fluent bit (log parsing --> exporter) 
* [ ] Multi-runner support on same host - (master/child architecture)

---

## 📢 Shoutout

This project is part of my journey transitioning from **DevOps → SRE**, focused on **deep observability, infrastructure monitoring**, and **log-to-metrics pipelines**. It’s built to solve real-world pain from GitHub runners used in production pipelines across cloud and hybrid setups. 
* Implemented in a real-world enterprise setup with 200+ GitHub runner servers distributed globally, helps trace which runners are most utilized across teams.

---

## 📬 Contact

Found a bug? Wanna contribute? Ping me via [GitHub issues](https://github.com/thineshsubramani/github-runner-prometheus-exporter/issues).
