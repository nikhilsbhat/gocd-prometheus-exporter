package exporter

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/gocd"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	mutex              sync.Mutex
	client             *gocd.Config
	logger             log.Logger
	pipelinePath       []string
	scrapeFailures     prometheus.Counter
	agentsCount        *prometheus.GaugeVec
	agentDisk          *prometheus.GaugeVec
	pipelinesDiskUsage *prometheus.GaugeVec
	agentDown          *prometheus.GaugeVec
	serverHealth       *prometheus.GaugeVec
}

func NewExporter(logger log.Logger, client *gocd.Config, paths []string) *Exporter {
	return &Exporter{
		logger:       logger,
		pipelinePath: paths,
		client:       client,
		agentsCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      "agents_count",
			Help:      "number of GoCd agents",
		},
			[]string{"agents_count"},
		),
		agentDisk: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      "agent_disk_space",
			Help:      "information of GoCd agent's disk space availability",
		},
			[]string{"name", "id", "version", "os", "sandbox"},
		),
		agentDown: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      "agent_down",
			Help:      "latest information on GoCd agent's state",
		},
			[]string{"name", "id", "version", "os", "sandbox", "state"},
		),
		pipelinesDiskUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      "pipeline_size",
			Help:      "disk size that GoCd pipeline have occupied in bytes",
		},
			[]string{"pipeline_path", "type"},
		),
		serverHealth: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.Namespace,
			Name:      "server_health",
			Help:      "errors and warning ini GoCd server",
		},
			[]string{"type", "message"},
		),
	}
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	// fetching server health status
	healthInfo, err := e.client.GetHealthInfo()
	if err != nil {
		level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving server health information errored with: %s", err.Error())) //nolint:errcheck
	}
	for _, health := range healthInfo {
		e.serverHealth.WithLabelValues(health.Level, health.Message).Set(1)
	}

	// fetching all node information
	nodes, err := e.client.GetNodesInfo()
	if err != nil {
		level.Error(e.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving agents information errored with: %s", err.Error())) //nolint:errcheck
	}
	e.agentsCount.WithLabelValues("all").Set(float64(len(nodes.Config.Config)))
	e.agentDown.Reset()
	for _, node := range nodes.Config.Config {
		if node.CurrentState == common.GoCdDisconnectedState {
			e.agentDown.WithLabelValues(node.Name, node.ID, node.Version, node.OS, node.Sandbox, node.CurrentState).Set(1)
		}
		e.agentDisk.WithLabelValues(node.Name, node.ID, node.Version, node.OS, node.Sandbox).Set(common.Float(node.DiskSpaceAvailable))
	}

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

	e.agentsCount.Collect(ch)
	e.agentDisk.Collect(ch)
	e.agentDown.Collect(ch)
	e.pipelinesDiskUsage.Collect(ch)
	e.serverHealth.Collect(ch)
	return nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.pipelinesDiskUsage.Describe(ch)
	e.agentsCount.Describe(ch)
	e.agentDown.Describe(ch)
	e.agentDisk.Describe(ch)
	e.serverHealth.Describe(ch)
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
