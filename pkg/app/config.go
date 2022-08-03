package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// Config holds the information of the app gocd-prometheus-exporter.
type Config struct {
	GoCdBaseURL           string            `json:"gocd-base-url,omitempty" yaml:"gocd-base-url,omitempty"`
	GoCdUserName          string            `json:"gocd-username,omitempty" yaml:"gocd-username,omitempty"`
	GoCdPassword          string            `json:"gocd-password,omitempty" yaml:"gocd-password,omitempty"`
	InsecureTLS           bool              `json:"insecure-tls,omitempty" yaml:"insecure-tls,omitempty"`
	GoCdPipelinesPath     []string          `json:"gocd-pipelines-path,omitempty" yaml:"gocd-pipelines-path,omitempty"`
	GoCdPipelinesRootPath string            `json:"gocd-pipelines-root-path,omitempty" yaml:"gocd-pipelines-root-path,omitempty"`
	CaPath                string            `json:"ca-path,omitempty" yaml:"ca-path,omitempty"`
	Port                  int               `json:"port,omitempty" yaml:"port,omitempty"`
	Endpoint              string            `json:"metric-endpoint,omitempty" yaml:"metric-endpoint,omitempty"`
	LogLevel              string            `json:"log-level,omitempty" yaml:"log-level,omitempty"`
	SkipMetrics           []string          `json:"skip-metrics,omitempty" yaml:"skip-metrics,omitempty"`
	APICron               string            `json:"api-cron-schedule,omitempty" yaml:"api-cron-schedule,omitempty"`
	DiskCron              string            `json:"disk-cron-schedule,omitempty" yaml:"disk-cron-schedule,omitempty"`
	MetricCron            map[string]string `json:"metric-cron,omitempty" yaml:"metric-cron,omitempty"`
	AppGraceDuration      time.Duration     `json:"grace-duration,omitempty" yaml:"grace-duration,omitempty"`
}

// GetConfig returns the new instance of Config.
func GetConfig(conf Config, path string) (*Config, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Printf("config file %s not found, dropping configurations from file", path)

		return &conf, fmt.Errorf("fetching config file information failed with: %w", err)
	}

	fileOUT, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("failed to read the config file, dropping configurations from file")

		return &conf, fmt.Errorf("reading config file errored with: %w", err)
	}

	var newConfig Config
	if err = yaml.Unmarshal(fileOUT, &newConfig); err != nil {
		log.Println("failed to unmarshall configurations, dropping configurations from file")

		return &conf, fmt.Errorf("parsing config file errored with: %w", err)
	}
	if err = mergo.Merge(&newConfig, &conf); err != nil {
		log.Println("failed to merge configurations, dropping configurations from file")

		return &conf, fmt.Errorf("merging config errored with: %w", err)
	}

	return &newConfig, nil
}
