package gocd

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"
)

// Config holds resty.Client which could be used for interacting with GoCd and other information
type Config struct {
	client *resty.Client
	logger log.Logger
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
	DiskSpaceAvailable interface{} `json:"free_space,string,omitempty"`
}

// ServerVersion holds information GoCd server
type ServerVersion struct {
	Version     string `json:"version,omitempty"`
	GitSha      string `json:"git_sha,omitempty"`
	FullVersion string `json:"full_version,omitempty"`
	CommitURL   string `json:"commit_url,omitempty"`
}

type ServerHealth struct {
	Level   string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
}

// NewConfig returns new instance of Config when invoked
func NewConfig(baseURL, userName, passWord string, caContent []byte, logger log.Logger) *Config {
	newClient := resty.New()
	newClient.SetBaseURL(baseURL)
	if len(caContent) != 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caContent)
		newClient.SetTLSClientConfig(&tls.Config{RootCAs: certPool})
		newClient.SetBasicAuth(userName, passWord)
	} else {
		newClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	return &Config{
		client: newClient,
		logger: logger,
	}
}
