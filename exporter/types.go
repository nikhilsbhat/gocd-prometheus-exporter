package exporter

import (
	"sync"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	mutex              sync.Mutex
	logger             log.Logger
	skipMetrics        []string
	scrapeFailures     prometheus.Counter
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

type Config struct {
	GoCdBaseURL           string   `json:"gocd-base-url,omitempty" yaml:"gocd-base-url,omitempty"`
	GoCdUserName          string   `json:"gocd-username,omitempty" yaml:"gocd-username,omitempty"`
	GoCdPassword          string   `json:"gocd-password,omitempty" yaml:"gocd-password,omitempty"`
	InsecureTLS           bool     `json:"insecure-tls,omitempty" yaml:"insecure-tls,omitempty"`
	GoCdPipelinesPath     []string `json:"gocd-pipelines-path,omitempty" yaml:"gocd-pipelines-path,omitempty"`
	GoCdPipelinesRootPath string   `json:"gocd-pipelines-root-path,omitempty" yaml:"gocd-pipelines-root-path,omitempty"`
	CaPath                string   `json:"ca-path,omitempty" yaml:"ca-path,omitempty"`
	Port                  int      `json:"port,omitempty" yaml:"port,omitempty"`
	Endpoint              string   `json:"metric-endpoint,omitempty" yaml:"metric-endpoint,omitempty"`
	LogLevel              string   `json:"log-level,omitempty" yaml:"log-level,omitempty"`
	SkipMetrics           []string `json:"skip-metrics,omitempty" yaml:"skip-metrics,omitempty"`
	OtherCron             string   `json:"cron,omitempty" yaml:"cron,omitempty"`
	DiskCron              string   `json:"disk-cron,omitempty" yaml:"disk-cron,omitempty"`
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
