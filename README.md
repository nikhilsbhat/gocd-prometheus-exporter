# GoCD prometheus exporter


[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/gocd-prometheus-exporter)](https://goreportcard.com/report/github.com/nikhilsbhat/gocd-prometheus-exporter)
[![shields](https://img.shields.io/badge/license-MIT-blue)](https://github.com/nikhilsbhat/gocd-prometheus-exporter/blob/master/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/gocd-prometheus-exporter?status.svg)](https://godoc.org/github.com/nikhilsbhat/gocd-prometheus-exporter)
[![shields](https://img.shields.io/github/v/tag/nikhilsbhat/gocd-prometheus-exporter.svg)](https://github.com/nikhilsbhat/gocd-prometheus-exporter/tags)
[![shields](https://img.shields.io/github/downloads/nikhilsbhat/gocd-prometheus-exporter/total.svg)](https://github.com/nikhilsbhat/gocd-prometheus-exporter/releases)


prometheus exporter for `GoCD` that helps in collecting few metrics from [GoCD](https://www.gocd.org/) server.

## Introduction

A Prometheus exporter that helps in collecting various metadata and other information from GoCD and exposes it as metrics.

And this interacts with the `GoCD` server's API to collect information and expose it as a metric.

It schedules both metric collections from the `GoCD` server and pipeline size as cron to reduce resource spikes when /metrics is invoked.
Not all the data in the CI tool changes that frequently, so we do not need to overburden it with an enormous number of API calls (most of which turn out to be expensive).</br>
And the cron would address this issue by invoking these APIs on a schedule basis, cache them locally, and serve them when they are scraped.

## Requirements

* [Go](https://golang.org/dl/) 1.19 or above . Installing go can be found [here](https://golang.org/doc/install).
* Basic understanding of prometheus exporter and its golang [client](https://github.com/prometheus/client_golang.git) libraries and [building](https://prometheus.io/docs/guides/go-application/) them.


## Installation

* Recommend installing released versions. Release binaries are available on the [releases](https://github.com/nikhilsbhat/gocd-prometheus-exporter/releases) page and docker from [here](https://hub.docker.com/repository/docker/basnik/gocd-prometheus-exporter).
* Can always build it locally by running `go build` against cloned repo.

## Usage
```shell
NAME:
   gocd-prometheus-exporter - Utility to collect metrics of GoCd for prometheus

USAGE:
   gocd-prometheus-exporter [flags]

COMMANDS:
   version, v  version of the gocd-prometheus-exporter
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-cron-schedule value, --cron value        cron expression to schedule the metric collection.
                                                    - 'gocd-prometheus-exporter' schedules the job to collect the metrics in the specified intervals
                                                      and stores the latest values in memory.
                                                    - This is to reduce the load on the GoCd server when api requests are made to GoCd.
                                                    - All expressions supported by https://github.com/robfig/cron will be supported (default: "@every 30s")
   --ca-path value, --ca value                    path to file containing CA information to make secure connections to GoCd
   --config-file value, -c value                  path to file containing configurations for exporter (default: "$HOME/gocd-prometheus-exporter.yaml")
   --disk-cron-schedule value, --disk-cron value  cron expression to schedule the pipeline disk usage metric collection.
                                                    - This is to reduce the reduce resource consumed while computing the pipeline disk size. (default: "@every 30s")
   --endpoint value, -e value                     path under which the metrics to be exposed (default: "/metrics")
   --goCd-password value, --password value        required password for establishing connection to GoCd server if auth enabled [$GOCD_PASSWORD]
   --goCd-server-url value, --server value        GoCd server url to which the exporter needs to be connected
   --goCd-username value, --user value            required username for establishing connection to GoCd server if auth enabled [$GOCD_USERNAME]
   --grace-duration value, -d value               time duration to wait before stopping the service (default: 5s)
   --help, -h                                     show help (default: false)
   --insecure-tls, --insecure                     enable insecure TLS if you wish to connect to GOCD insecurily (default: false)
   --log-level value, --log value                 set log level for the GoCd exporter (default: "info")
   --pipeline-path value, --pt value              list of paths to pipelines that needs tp be monitored  (accepts multiple inputs)
   --pipeline-root-path value, --root-path value  root path of pipeline that needs to be monitored
   --port value, -p value                         port on which the metrics to be exposed (default: 8090)
   --skip-metrics value, --sk value               list of metrics to be skipped  (accepts multiple inputs)
```

### Run

```shell
gocd-prometheus-exporter --goCd-server-url http://gocdurl:8153/go --goCd-password username --goCd-password password
```

Prometheus scrape config, so that prometheus starts scraping metrics from the exporter.
```
scrape_configs:
  - job_name: gocd-artifact-monitor
    scrape_interval: 5s
    static_configs:
      - targets:
        - "<gocd-prometheus-exporter-domain/IP>:8090"
```
All configurations required by exporter can be passed via configuration file, config file is picked automatically if it is placed at `$HOME/gocd-prometheus-exporter.yaml`

Sample config file can be found [here](https://github.com/nikhilsbhat/gocd-prometheus-exporter/blob/master/gocd-prometheus-exporter.sample.yaml). Check for all yaml fields [here](https://github.com/nikhilsbhat/gocd-prometheus-exporter/blob/master/pkg/app/config.go#L15) (check yaml tags)

## Metrics supported by the `gocd-prometheus-exporter` currently.

```
# HELP gocd_agent_disk_space information of GoCd agent's disk space availability
# TYPE gocd_agent_disk_space gauge
gocd_agent_disk_space{id="",name="",os="",sandbox="",version=""} 0
# HELP gocd_admin_count number users who are admins in gocd
# TYPE gocd_admin_count gauge
gocd_admin_count{users="all"} 0
# HELP gocd_agent_down latest information on GoCd agent's state
# TYPE gocd_agent_down gauge
gocd_agent_down{id="",name="",os="",sandbox="",version=""} 0
# HELP gocd_agents_count number of GoCd agents
# TYPE gocd_agents_count gauge
gocd_agents_count{agents_count="all"} 0
# HELP gocd_backup_configured would be 1 if backup is enabled
# TYPE gocd_backup_configured gauge
gocd_backup_configured{failure_email="false",success_email="false"} 0
# HELP gocd_config_repo_count number of config repos
# TYPE gocd_config_repo_count gauge
gocd_config_repo_count{repos="all"} 0
# HELP gocd_pipeline_group_count number of pipeline groups
# TYPE gocd_pipeline_group_count gauge
gocd_pipeline_group_count{pipeline_groups="all"} 0
# HELP gocd_pipeline_size disk size that GoCd pipeline have occupied in bytes
# TYPE gocd_pipeline_size gauge
gocd_pipeline_size{pipeline_path="",type="dir"} 0
gocd_pipeline_size{pipeline_path="",type="link"} 0
# HELP gocd_plugins GoCD elastic agents plugin usage
# TYPE gocd_plugins gauge
gocd_plugins{bundled="true",id="yaml.config.plugin",state="active"} 1
gocd_plugins{bundled="true",id="yum",state="invalid"} 1
# HELP gocd_pipeline_state GoCD pipeline state
# TYPE gocd_pipeline_state gauge
gocd_pipeline_state{locked="false",name="action-movies",paused="false",paused_by="",schedulable="true"} 0
gocd_pipeline_state{locked="false",name="action-movies-manual",paused="false",paused_by="",schedulable="true"} 0
# HELP gocd_elastic_agent_profile_usage GoCD elastic agents profile usage
# TYPE gocd_elastic_agent_profile_usage gauge
gocd_elastic_agent_profile_usage{job_name="lint",name="sample_ec2",pipeline_config_origin="gocd",pipeline_name="helm-drift",stage_name="lint"} 1
# HELP gocd_pipeline_not_run GoCD pipeline not run in last X days
# TYPE gocd_pipeline_not_run gauge
gocd_pipeline_not_run{pipeline="gocd-prometheus-exporter",scheduled_date="25 Jun, 2023 at 19:47:02 [+0530]",scheduled_timestamp="1687702622382"} 0.12195940645833332
gocd_pipeline_not_run{pipeline="helm-images",scheduled_date="25 Jun, 2023 at 19:48:16 [+0530]",scheduled_timestamp="1687702696672"} 0.12109960001157406
# HELP gocd_server_health errors and warning present in GoCD
# TYPE gocd_server_health gauge
gocd_server_health{message="Message from plugin: com.thoughtworks.gocd.analytics",type="ERROR"} 1
gocd_server_health{message="Message from plugin: com.thoughtworks.gocd.analytics",type="WARNING"} 1
# HELP gocd_version GoCD server version
# TYPE gocd_version gauge
gocd_version{full_version="22.1.0 (13913-f4c9c1650e2e27fe0a9962faa39536f94f57e297)",git_sha="f4c9c1650e2e27fe0a9962faa39536f94f57e297",version="22.1.0"} 1
```

## Documentation

Updated documentation on all available commands and flags can be found [here](https://github.com/nikhilsbhat/gocd-prometheus-exporter/blob/master/docs/doc/gocd_prometheus_exporter.md).

### Note
The environment variable `GOCD_PIPELINE_DAYS_TO_SKIP` can be set with a value to drop pipeline from consideration for the metric `gocd_pipeline_not_run`.</br>
**Ex:** Setting `GOCD_PIPELINE_DAYS_TO_SKIP` to 10 would make one drop all pipelines whose last run is not within 10 days

## TODO
* [ ] Expose more metrics.
* [ ] Increase test coverage.