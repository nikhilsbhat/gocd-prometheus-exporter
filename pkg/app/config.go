package app

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

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
	ApiCron               string   `json:"api-cron-schedule,omitempty" yaml:"api-cron-schedule,omitempty"`
	DiskCron              string   `json:"disk-cron,omitempty" yaml:"disk-cron,omitempty"`
}

func GetConfig(conf Config, path string) (*Config, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Printf("config file %s not found, dropping configurations from file", path)
		return &conf, err
	}

	fileOUT, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("failed to read the config file, dropping configurations from file")
		return &conf, err
	}

	var newConfig Config
	if err = yaml.Unmarshal(fileOUT, &newConfig); err != nil {
		log.Println("failed to unmarshall configurations, dropping configurations from file")
		return &conf, err
	}
	if err = mergo.Merge(&newConfig, &conf, mergo.WithOverride); err != nil {
		log.Println("failed to merge configurations, dropping configurations from file")
		return &conf, err
	}
	return &newConfig, nil
}
