package exporter

import (
	"sync"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	mutex              sync.Mutex
	logger             log.Logger
	skipMetrics        []string
	agentsCount        *prometheus.GaugeVec
	agentDisk          *prometheus.GaugeVec
	pipelinesDiskUsage *prometheus.GaugeVec
	agentDown          *prometheus.GaugeVec
	serverHealth       *prometheus.GaugeVec
	configRepoCount    *prometheus.GaugeVec
	pipelineGroupCount *prometheus.GaugeVec
	backupConfigured   *prometheus.GaugeVec
	adminCount         *prometheus.GaugeVec
}

func NewExporter(logger log.Logger, skipMetrics []string) *Exporter {
	return &Exporter{
		logger:      logger,
		skipMetrics: skipMetrics,
		agentsCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentsCount,
			Help:      "number of GoCd agents",
		}, []string{"agents_count"},
		),
		agentDisk: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentDiskSpace,
			Help:      "information of GoCd agent's disk space availability",
		}, []string{"name", "id", "version", "os", "sandbox"},
		),
		agentDown: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricAgentDown,
			Help:      "latest information on GoCd agent's state",
		}, []string{"name", "id", "version", "os", "sandbox", "state"},
		),
		pipelinesDiskUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricPipelineSize,
			Help:      "disk size that GoCd pipeline have occupied in bytes",
		}, []string{"pipeline_path", "type"},
		),
		serverHealth: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricServerHealth,
			Help:      "errors and warning ini GoCd server",
		}, []string{"type", "message"},
		),
		configRepoCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricConfigRepoCount,
			Help:      "number of config repos",
		}, []string{"repos"},
		),
		adminCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      common.MetricSystemAdminsCount,
			Help:      "number users who are admins in gocd",
		}, []string{"users"},
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
