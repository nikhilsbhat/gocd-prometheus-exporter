package gocd

import (
	"crypto/tls"
	"crypto/x509"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"
)

// client holds resty.Client which could be used for interacting with GoCD and other information.
type client struct {
	client   *resty.Client
	logger   log.Logger
	apiCron  string
	diskCron string
	lock     sync.RWMutex
	paths    []string
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
func NewClient(baseURL, userName, passWord, loglevel, cron, diskCron string, caContent []byte, path []string, logger log.Logger) GoCd {
	newClient := resty.New()
	newClient.SetRetryCount(defaultRetryCount)
	newClient.SetRetryWaitTime(defaultRetryWaitTime * time.Second)
	if loglevel == "debug" {
		newClient.SetDebug(true)
	}
	newClient.SetBaseURL(baseURL)
	if len(caContent) != 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caContent)
		newClient.SetTLSClientConfig(&tls.Config{RootCAs: certPool}) //nolint:gosec
		newClient.SetBasicAuth(userName, passWord)
	} else {
		newClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}

	return &client{
		client:   newClient,
		logger:   logger,
		lock:     sync.RWMutex{},
		apiCron:  cron,
		diskCron: diskCron,
		paths:    path,
	}
}
