package exporter

import (
	"sync"

	"github.com/go-kit/log"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	mutex              sync.Mutex
	logger             log.Logger
	skipMetrics        []string
	agentsCount        *prometheus.GaugeVec
	agentDisk          *prometheus.GaugeVec
	agentDown          *prometheus.GaugeVec
	serverHealth       *prometheus.GaugeVec
	configRepoCount    *prometheus.GaugeVec
	pipelineGroupCount *prometheus.GaugeVec
	pipelineCount      *prometheus.GaugeVec
	backupConfigured   *prometheus.GaugeVec
	adminCount         *prometheus.GaugeVec
	environmentCount   *prometheus.GaugeVec
	versionInfo        *prometheus.GaugeVec
	jobStatus          *prometheus.GaugeVec
	pipelineState      *prometheus.GaugeVec
}

func NewExporter(logger log.Logger, skipMetrics []string) *Exporter {
	return &Exporter{
		logger:      logger,
		skipMetrics: skipMetrics,
		agentsCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentsCount,
			Help:      "number of GoCD agents",
		}, []string{"agents_count"},
		),
		agentDisk: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentDiskSpace,
			Help:      "information of GoCD agent's disk space availability",
		}, []string{"name", "id", "version", "os", "sandbox"},
		),
		agentDown: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentDown,
			Help:      "latest information on GoCD agent's state",
		}, []string{"name", "id", "version", "os", "sandbox", "state", "config_state"},
		),
		serverHealth: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricServerHealth,
			Help:      "errors and warning present in GoCD",
		}, []string{"type", "message"},
		),
		configRepoCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricConfigRepoCount,
			Help:      "number of config repos present in GoCD",
		}, []string{"repos"},
		),
		adminCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricSystemAdminsCount,
			Help:      "number users who are admins in GoCD",
		}, []string{"users"},
		),
		backupConfigured: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricConfiguredBackup,
			Help:      "would be 1 if backup is enabled",
		}, []string{"success_email", "failure_email", "scheduled"}),
		pipelineGroupCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricPipelineGroupCount,
			Help:      "number of pipeline groups present in GoCD",
		}, []string{"pipeline_groups"}),
		pipelineCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricPipelineCount,
			Help:      "total number of pipeline present in GoCD",
		}, []string{"pipeline_count"}),
		environmentCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricEnvironmentCountAll,
			Help:      "total number of environment present in GoCD",
		}, []string{"environment_count"}),
		versionInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricVersion,
			Help:      "GoCD server version",
		}, []string{"version", "git_sha", "full_version"}),
		jobStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricJobStatus,
			Help:      "GoCD pipeline status",
		}, []string{"name", "job", "stage", "state"}),
		pipelineState: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricPipelineState,
			Help:      "GoCD pipeline state",
		}, []string{"name", "paused", "locked", "schedulable", "paused_by"}),
	}
}
