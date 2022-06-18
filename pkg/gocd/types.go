package gocd

import (
	"crypto/tls"
	"crypto/x509"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"
)

var (
	CurrentNodeConfig    []Node
	CurrentServerHealth  []ServerHealth
	CurrentConfigRepos   []ConfigRepo
	CurrentPipelineGroup []PipelineGroup
	CurrentSystemAdmins  SystemAdmins
	CurrentBackupConfig  BackupConfig
	CurrentPipelineSize  = make(map[string]PipelineSize)
)

const (
	defaultRetryCount    = 5
	defaultRetryWaitTime = 5
)

// Config holds resty.Client which could be used for interacting with GoCd and other information
type Config struct {
	client   *resty.Client
	logger   log.Logger
	apiCron  string
	diskCron string
	lock     sync.RWMutex
	paths    []string
}

// NodesConfig holds information of all agent of GoCd
type NodesConfig struct {
	Config Nodes `json:"_embedded,omitempty"`
}

// Nodes holds information of all agent of GoCd
type Nodes struct {
	Config []Node `json:"agents,omitempty"`
}

// Node holds information of a particular agent
type Node struct {
	Name               string      `json:"hostname,omitempty"`
	ID                 string      `json:"uuid,omitempty"`
	Version            string      `json:"agent_version,omitempty"`
	CurrentState       string      `json:"agent_state,omitempty"`
	OS                 string      `json:"operating_system,omitempty"`
	ConfigState        string      `json:"agent_config_state,omitempty"`
	Sandbox            string      `json:"sandbox,omitempty"`
	DiskSpaceAvailable interface{} `json:"free_space,omitempty"`
}

// ServerVersion holds version information GoCd server
type ServerVersion struct {
	Version     string `json:"version,omitempty"`
	GitSha      string `json:"git_sha,omitempty"`
	FullVersion string `json:"full_version,omitempty"`
	CommitURL   string `json:"commit_url,omitempty"`
}

// ServerHealth holds information of GoCd server health
type ServerHealth struct {
	Level   string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
}

// ConfigRepoConfig holds information of all config-repos present in GoCd
type ConfigRepoConfig struct {
	ConfigRepos ConfigRepos `json:"_embedded,omitempty"`
}

// ConfigRepos holds information of all config-repos present in GoCd
type ConfigRepos struct {
	ConfigRepos []ConfigRepo `json:"config_repos,omitempty"`
}

// ConfigRepo holds information of the specified config-repo
type ConfigRepo struct {
	ID       string `json:"config_repos,omitempty"`
	Material struct {
		Type       string `json:"type,omitempty"`
		Attributes struct {
			URL        string `json:"url,omitempty"`
			Branch     string `json:"branch,omitempty"`
			AutoUpdate bool   `json:"auto_update,omitempty"`
		}
	}
}

type PipelineGroupsConfig struct {
	PipelineGroups PipelineGroups `json:"_embedded,omitempty"`
}

type PipelineGroups struct {
	PipelineGroups []PipelineGroup `json:"groups,omitempty"`
}

type PipelineGroup struct {
	Name          string `json:"name,omitempty"`
	PipelineCount int    `json:"pipeline_count,omitempty"`
	Pipelines     []struct {
		Name string `json:"name,omitempty"`
	}
}

// SystemAdmins holds information of the system admins present
type SystemAdmins struct {
	Roles []string `json:"roles,omitempty"`
	Users []string `json:"users,omitempty"`
}

// BackupConfig holds information of the backup configured
type BackupConfig struct {
	EmailOnSuccess bool   `json:"email_on_success,omitempty"`
	EmailOnFailure bool   `json:"email_on_failure,omitempty"`
	Schedule       string `json:"schedule,omitempty"`
}

type PipelineSize struct {
	Size float64
	Type string
}

// NewConfig returns new instance of Config when invoked
func NewConfig(baseURL, userName, passWord, loglevel, cron, diskCron string, caContent []byte, path []string, logger log.Logger) *Config {
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
	return &Config{
		client:   newClient,
		logger:   logger,
		lock:     sync.RWMutex{},
		apiCron:  cron,
		diskCron: diskCron,
		paths:    path,
	}
}
