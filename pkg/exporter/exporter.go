package exporter

import (
	"strconv"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"

	"github.com/thoas/go-funk"

	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collect(ch chan<- prometheus.Metric) {
	// fetching server health status
	if !funk.Contains(e.skipMetrics, common.MetricServerHealth) {
		for _, health := range gocd.CurrentServerHealth {
			e.serverHealth.WithLabelValues(health.Level, health.Message).Set(1)
		}
		e.serverHealth.Collect(ch)
	}

	// fetching agent count metrics
	if !funk.Contains(e.skipMetrics, common.MetricAgentsCount) {
		e.agentsCount.WithLabelValues("all").Set(float64(len(gocd.CurrentNodeConfig)))
		e.agentsCount.Collect(ch)
	}

	// fetching agent down and agent disk space metrics
	if !funk.Contains(e.skipMetrics, common.MetricAgentDown) || !funk.Contains(e.skipMetrics, common.MetricAgentDiskSpace) {
		e.agentDown.Reset()
		for _, node := range gocd.CurrentNodeConfig {
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

	// fetching repo count metrics
	if !funk.Contains(e.skipMetrics, common.MetricConfigRepoCount) {
		e.configRepoCount.WithLabelValues("all").Set(float64(len(gocd.CurrentConfigRepos)))
		e.configRepoCount.Collect(ch)
	}

	// fetching system admins metrics
	if !funk.Contains(e.skipMetrics, common.MetricSystemAdminsCount) {
		e.adminCount.WithLabelValues("all").Set(float64(len(gocd.CurrentSystemAdmins.Users)))
		e.adminCount.Collect(ch)
	}

	// fetching pipeline group metrics
	if !funk.Contains(e.skipMetrics, common.MetricPipelineGroupCount) {
		e.pipelineGroupCount.WithLabelValues("all").Set(float64(len(gocd.CurrentPipelineGroup)))
		e.pipelineGroupCount.Collect(ch)
	}

	// fetching backup metrics
	if !funk.Contains(e.skipMetrics, common.MetricConfiguredBackup) {
		var enabled float64
		if len(gocd.CurrentBackupConfig.Schedule) != 0 {
			enabled = 1.0
		}
		e.backupConfigured.WithLabelValues(
			strconv.FormatBool(gocd.CurrentBackupConfig.EmailOnSuccess),
			strconv.FormatBool(gocd.CurrentBackupConfig.EmailOnFailure)).Set(enabled)
		e.backupConfigured.Collect(ch)
	}

	// fetching pipeline sizes
	if !funk.Contains(e.skipMetrics, common.MetricPipelineSize) {
		for pipeline, pipelineInfo := range gocd.CurrentPipelineSize {
			e.pipelinesDiskUsage.WithLabelValues(pipeline, pipelineInfo.Type).Set(pipelineInfo.Size)
		}
		e.pipelinesDiskUsage.Collect(ch)
	}
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
	e.collect(ch)
}
