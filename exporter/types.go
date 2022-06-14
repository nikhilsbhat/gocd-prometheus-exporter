package exporter

import (
	"sync"

	"github.com/go-kit/log"
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

type Config struct {
	GoCdBaseURL           string   `json:"gocd-base-url,omitempty" yaml:"gocd-base-url,omitempty"`
	GoCdUserName          string   `json:"gocd-username,omitempty" yaml:"gocd-username,omitempty"`
	GoCdPassword          string   `json:"gocd-password,omitempty" yaml:"gocd-password,omitempty"`
	InsecureTLS           bool     `json:"insecure-tls,omitempty" yaml:"insecure-tls,omitempty"`
	GoCdPipelinesPath     []string `json:"gocd-pipelines-path,omitempty" yaml:"gocd-pipelines-path,omitempty"`
	GoCdPipelinesRootPath string   `json:"gocd-pipca-pathelines-root-path,omitempty" yaml:"gocd-pipelines-root-path,omitempty"`
	CaPath                string   `json:"ca-path,omitempty" yaml:"ca-path,omitempty"`
	Port                  int      `json:"port,omitempty" yaml:"port,omitempty"`
	Endpoint              string   `json:"metric-endpoint,omitempty" yaml:"metric-endpoint,omitempty"`
	LogLevel              string   `json:"log-level,omitempty" yaml:"log-level,omitempty"`
}
