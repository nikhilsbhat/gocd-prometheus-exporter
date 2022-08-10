package gocd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"
)

// client holds resty.Client which could be used for interacting with GoCD and other information.
type client struct {
	client             *resty.Client
	logger             log.Logger
	defaultAPICron     string
	diskCron           string
	metricSpecificCron map[string]string
	lock               sync.RWMutex
	paths              []string
	skipMetrics        []string
}

// GoCd implements methods to get various information regarding GoCD.
type GoCd interface {
	GetAgentsInfo() ([]Agent, error)
	GetAgentJobRunHistory() ([]AgentJobHistory, error)
	GetDiskSize(path string) (float64, string, error)
	GetHealthInfo() ([]ServerHealth, error)
	GetConfigRepoInfo() ([]ConfigRepo, error)
	GetAdminsInfo() (SystemAdmins, error)
	GetPipelineGroupInfo() ([]PipelineGroup, error)
	GetEnvironmentInfo() ([]Environment, error)
	GetVersionInfo() (VersionInfo, error)
	GetBackupInfo() (BackupConfig, error)
	CronSchedulers()
}

// NewClient returns new instance of client when invoked.
func NewClient(baseURL, userName, passWord, loglevel, defaultAPICron, diskCron string,
	metricSpecificCron map[string]string,
	caContent []byte,
	path, skipMetrics []string,
	logger log.Logger,
) GoCd {
	newClient := resty.New()
	newClient.SetRetryCount(defaultRetryCount)
	newClient.SetRetryWaitTime(defaultRetryWaitTime * time.Second)
	newClient.SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
		return 0, fmt.Errorf("quota exceeded") //nolint:goerr113
	})
	if loglevel == "debug" {
		newClient.SetDebug(true)
	}
	newClient.SetBaseURL(baseURL)
	newClient.SetBasicAuth(userName, passWord)
	if len(caContent) != 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caContent)
		newClient.SetTLSClientConfig(&tls.Config{RootCAs: certPool}) //nolint:gosec
	} else {
		newClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}

	return &client{
		client:             newClient,
		logger:             logger,
		lock:               sync.RWMutex{},
		defaultAPICron:     defaultAPICron,
		diskCron:           diskCron,
		metricSpecificCron: metricSpecificCron,
		paths:              path,
		skipMetrics:        skipMetrics,
	}
}

func (conf *client) getCron(metric string) string {
	if metric == common.MetricPipelineSize {
		return conf.diskCron
	}
	if val, ok := conf.metricSpecificCron[metric]; ok {
		level.Debug(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("the cron for metric %s would be %s", metric, val)) //nolint:errcheck

		return val
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("metric %s would be using default cron", metric)) //nolint:errcheck

	return conf.defaultAPICron
}

func (conf *client) getCronClient() *client {
	return &client{
		client:             conf.client,
		logger:             conf.logger,
		lock:               sync.RWMutex{},
		defaultAPICron:     conf.defaultAPICron,
		diskCron:           conf.diskCron,
		metricSpecificCron: conf.metricSpecificCron,
		paths:              conf.paths,
		skipMetrics:        conf.skipMetrics,
	}
}
