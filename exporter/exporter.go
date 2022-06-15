package exporter

import (
	"fmt"
	"strconv"

	"github.com/thoas/go-funk"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/gocd"
	"github.com/prometheus/client_golang/prometheus"
)

func NewExporter(logger log.Logger, client *gocd.Config, paths, skipMetrics []string) *Exporter {
	return &Exporter{
		logger:       logger,
		pipelinePath: paths,
		client:       client,
		skipMetrics:  skipMetrics,
		agentsCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentsCount,
			Help:      "number of GoCd agents",
		},
			[]string{"agents_count"},
		),
		agentDisk: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentDiskSpace,
			Help:      "information of GoCd agent's disk space availability",
		},
			[]string{"name", "id", "version", "os", "sandbox"},
		),
		agentDown: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentDown,
			Help:      "latest information on GoCd agent's state",
		},
			[]string{"name", "id", "version", "os", "sandbox", "state"},
		),
		pipelinesDiskUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricPipelineSize,
			Help:      "disk size that GoCd pipeline have occupied in bytes",
		},
			[]string{"pipeline_path", "type"},
		),
		serverHealth: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricServerHealth,
			Help:      "errors and warning ini GoCd server",
		},
			[]string{"type", "message"},
		),
		configRepoCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricConfigRepoCount,
			Help:      "number of config repos",
		},
			[]string{"repos"},
		),
		adminCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricSystemAdminsCount,
			Help:      "number users who are admins in gocd",
		},
			[]string{"users"},
		),
		backupConfigured: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricConfiguredBackup,
			Help:      "would be 1 if backup is enabled",
		}, []string{"success_email", "failure_email"}),
		pipelineGroupCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricPipelineGroupCount,
			Help:      "number of pipeline groups",
		}, []string{"pipeline_groups"}),
	}
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	if !funk.Contains(e.skipMetrics, common.MetricServerHealth) {
		// fetching server health status
		healthInfo, err := e.client.GetHealthInfo()
		if err != nil {
			level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving server health information errored with: %s", err.Error())) //nolint:errcheck
		}
		for _, health := range healthInfo {
			e.serverHealth.WithLabelValues(health.Level, health.Message).Set(1)
		}
		e.serverHealth.Collect(ch)
	}

	// fetching all node information
	nodes, err := e.client.GetNodesInfo()
	if err != nil {
		level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving agents information errored with: %s", err.Error())) //nolint:errcheck
	}

	if !funk.Contains(e.skipMetrics, common.MetricAgentsCount) {
		e.agentsCount.WithLabelValues("all").Set(float64(len(nodes)))
		e.agentsCount.Collect(ch)
	}

	if !funk.Contains(e.skipMetrics, common.MetricAgentDown) || !funk.Contains(e.skipMetrics, common.MetricAgentDiskSpace) {
		e.agentDown.Reset()
		for _, node := range nodes {
			if node.CurrentState == common.GoCdDisconnectedState {
				e.agentDown.WithLabelValues(node.Name, node.ID, node.Version, node.OS, node.Sandbox, node.CurrentState).Set(1)
			}
			if !funk.Contains(e.skipMetrics, common.MetricAgentDiskSpace) {
				e.agentDisk.WithLabelValues(node.Name, node.ID, node.Version, node.OS, node.Sandbox).Set(common.Float(node.DiskSpaceAvailable))
			}
		}
		e.agentDown.Collect(ch)
		e.agentDisk.Collect(ch)
	}

	if !funk.Contains(e.skipMetrics, common.MetricConfigRepoCount) {
		repos, err := e.client.GetConfigRepoInfo()
		if err != nil {
			level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving config repo information errored with: %s", err.Error())) //nolint:errcheck
		}
		e.configRepoCount.WithLabelValues("all").Set(float64(len(repos)))
		e.configRepoCount.Collect(ch)
	}

	if !funk.Contains(e.skipMetrics, common.MetricSystemAdminsCount) {
		admins, err := e.client.GetAdminsInfo()
		if err != nil {
			level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving system admin information errored with: %s", err.Error())) //nolint:errcheck
		}
		e.adminCount.WithLabelValues("all").Set(float64(len(admins.Users)))
		e.adminCount.Collect(ch)
	}

	if !funk.Contains(e.skipMetrics, common.MetricPipelineGroupCount) {
		groups, err := e.client.GetPipelineGroupInfo()
		if err != nil {
			level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving pipeline group information errored with: %s", err.Error())) //nolint:errcheck
		}
		e.pipelineGroupCount.WithLabelValues("all").Set(float64(len(groups)))
		e.pipelineGroupCount.Collect(ch)
	}

	if !funk.Contains(e.skipMetrics, common.MetricConfiguredBackup) {
		backup, err := e.client.GetBackupInfo()
		if err != nil {
			level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving backup information errored with: %s", err.Error())) //nolint:errcheck
		}
		var enabled float64
		if len(backup.Schedule) != 0 {
			enabled = 1.0
		}
		e.backupConfigured.WithLabelValues(strconv.FormatBool(backup.EmailOnSuccess), strconv.FormatBool(backup.EmailOnFailure)).Set(enabled)
		e.backupConfigured.Collect(ch)
	}

	if !funk.Contains(e.skipMetrics, common.MetricPipelineSize) {
		// fetching pipeline sizes
		for i, pipeline := range e.pipelinePath {
			if len(pipeline) == 0 {
				continue
			}
			pipelineSize, pathType, err := e.client.GetDiskSize(pipeline)
			if err != nil {
				level.Error(e.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
			}
			if i == 0 {
				pathType = "all"
			}
			level.Debug(e.logger).Log(common.LogCategoryMsg, fmt.Sprintf("pipeline present at %s would be scanned", pipeline)) //nolint:errcheck
			e.pipelinesDiskUsage.WithLabelValues(pipeline, pathType).Set(pipelineSize)
		}
		e.pipelinesDiskUsage.Collect(ch)
	}

	return nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.pipelinesDiskUsage.Describe(ch)
	e.agentsCount.Describe(ch)
	e.agentDown.Describe(ch)
	e.agentDisk.Describe(ch)
	e.serverHealth.Describe(ch)
	e.configRepoCount.Describe(ch)
	e.adminCount.Describe(ch)
	e.backupConfigured.Describe(ch)
	e.pipelineGroupCount.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		level.Error(e.logger).Log(common.LogCategoryErr, "Error scraping GoCd:", "err", err) //nolint:errcheck
		e.scrapeFailures.Inc()
		e.scrapeFailures.Collect(ch)
	}
}
