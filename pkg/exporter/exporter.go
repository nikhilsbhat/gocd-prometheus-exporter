package exporter

import (
	"strconv"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"

	"github.com/thoas/go-funk"

	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collect(channel chan<- prometheus.Metric) { //nolint:funlen
	// fetching server health status
	if !funk.Contains(e.skipMetrics, common.MetricServerHealth) {
		for _, health := range gocd.CurrentServerHealth {
			e.serverHealth.WithLabelValues(health.Level, health.Message).Set(1)
		}
		e.serverHealth.Collect(channel)
	}

	// fetching agent count metrics
	if !funk.Contains(e.skipMetrics, common.MetricAgentsCount) {
		e.agentsCount.WithLabelValues("all").Set(float64(len(gocd.CurrentAgentsConfig)))
		e.agentsCount.Collect(channel)
	}

	// fetching agent down and agent disk space metrics
	if !funk.Contains(e.skipMetrics, common.MetricAgentDown) || !funk.Contains(e.skipMetrics, common.MetricAgentDiskSpace) {
		e.agentDown.Reset()
		for _, node := range gocd.CurrentAgentsConfig {
			if node.CurrentState == common.GoCdDisconnectedState {
				e.agentDown.WithLabelValues(node.Name, node.ID, node.Version, node.OS, node.Sandbox, node.CurrentState).Set(1)
			}
			if !funk.Contains(e.skipMetrics, common.MetricAgentDiskSpace) {
				e.agentDisk.WithLabelValues(node.Name, node.ID, node.Version, node.OS, node.Sandbox).Set(common.Float(node.DiskSpaceAvailable))
			}
		}
		e.agentDown.Collect(channel)
		e.agentDisk.Collect(channel)
	}

	// fetching repo count metrics
	if !funk.Contains(e.skipMetrics, common.MetricConfigRepoCount) {
		e.configRepoCount.WithLabelValues("all").Set(float64(len(gocd.CurrentConfigRepos)))
		e.configRepoCount.Collect(channel)
	}

	// fetching system admins metrics
	if !funk.Contains(e.skipMetrics, common.MetricSystemAdminsCount) {
		e.adminCount.WithLabelValues("all").Set(float64(len(gocd.CurrentSystemAdmins.Users)))
		e.adminCount.Collect(channel)
	}

	// fetching pipeline group metrics
	if !funk.Contains(e.skipMetrics, common.MetricPipelineGroupCount) {
		e.pipelineGroupCount.WithLabelValues("all").Set(float64(len(gocd.CurrentPipelineGroup)))
		e.pipelineGroupCount.Collect(channel)
	}

	// fetching pipeline count metrics
	if !funk.Contains(e.skipMetrics, common.MetricPipelineCount) {
		e.pipelineCount.WithLabelValues("all").Set(float64(gocd.CurrentPipelineCount))
		e.pipelineCount.Collect(channel)
	}

	if !funk.Contains(e.skipMetrics, common.MetricEnvironmentCountAll) {
		e.environmentCount.WithLabelValues("all").Set(float64(len(gocd.CurrentEnvironments)))
		e.environmentCount.Collect(channel)
	}

	// fetching backup metrics
	if !funk.Contains(e.skipMetrics, common.MetricConfiguredBackup) {
		var enabled float64
		if len(gocd.CurrentBackupConfig.Schedule) != 0 {
			enabled = 1.0
		}
		e.backupConfigured.WithLabelValues(
			strconv.FormatBool(gocd.CurrentBackupConfig.EmailOnSuccess),
			strconv.FormatBool(gocd.CurrentBackupConfig.EmailOnFailure),
			gocd.CurrentBackupConfig.Schedule).Set(enabled)
		e.backupConfigured.Collect(channel)
	}

	// fetching pipeline sizes
	if !funk.Contains(e.skipMetrics, common.MetricPipelineSize) {
		for pipeline, pipelineInfo := range gocd.CurrentPipelineSize {
			e.pipelinesDiskUsage.WithLabelValues(pipeline, pipelineInfo.Type).Set(pipelineInfo.Size)
		}
		e.pipelinesDiskUsage.Collect(channel)
	}

	if !funk.Contains(e.skipMetrics, common.MetricVersion) {
		if (gocd.CurrentVersion == gocd.VersionInfo{}) {
			e.versionInfo.WithLabelValues("", "", "").Set(0)
		} else {
			e.versionInfo.WithLabelValues(gocd.CurrentVersion.Version, gocd.CurrentVersion.GitSHA, gocd.CurrentVersion.FullVersion).Set(1)
		}
		e.versionInfo.Collect(channel)
	}

	if !funk.Contains(e.skipMetrics, common.MetricJobStatus) {
		for _, agentHistory := range gocd.CurrentAgentJobRunHistory {
			if len(agentHistory.Jobs) != 0 {
				switch agentHistory.Jobs[0].Result {
				case common.GoCdPipelineStateFail:
					e.jobStatus.WithLabelValues(agentHistory.Jobs[0].Name, agentHistory.Jobs[0].JobName,
						agentHistory.Jobs[0].StageName, agentHistory.Jobs[0].Result).Set(1)
				case common.GoCdPipelineStatePass:
					e.jobStatus.WithLabelValues(agentHistory.Jobs[0].Name, agentHistory.Jobs[0].JobName,
						agentHistory.Jobs[0].StageName, agentHistory.Jobs[0].Result).Set(1)
				default:
					e.jobStatus.WithLabelValues(agentHistory.Jobs[0].Name, agentHistory.Jobs[0].JobName,
						agentHistory.Jobs[0].StageName, agentHistory.Jobs[0].Result).Set(1)
				}
			}
		}
		e.jobStatus.Collect(channel)
	}
}

func (e *Exporter) Describe(channel chan<- *prometheus.Desc) {
	e.pipelinesDiskUsage.Describe(channel)
	e.agentsCount.Describe(channel)
	e.agentDown.Describe(channel)
	e.agentDisk.Describe(channel)
	e.serverHealth.Describe(channel)
	e.configRepoCount.Describe(channel)
	e.adminCount.Describe(channel)
	e.backupConfigured.Describe(channel)
	e.pipelineGroupCount.Describe(channel)
	e.pipelineCount.Describe(channel)
	e.environmentCount.Describe(channel)
	e.versionInfo.Describe(channel)
	e.jobStatus.Describe(channel)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	e.collect(ch)
}
