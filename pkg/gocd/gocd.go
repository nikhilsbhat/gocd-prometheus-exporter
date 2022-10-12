package gocd

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/app"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/nikhilsbhat/gocd-sdk-go"
)

// client holds resty.Client which could be used for interacting with GoCD and other information.
type client struct {
	logger    log.Logger
	config    app.Config
	lock      sync.RWMutex
	caContent []byte
}

// GoCd implements methods to get various information regarding GoCD.
type GoCd interface {
	CronSchedulers()
}

// NewClient returns new instance of client when invoked.
func NewClient(config app.Config, logger log.Logger, ca []byte) GoCd {
	return &client{
		config:    config,
		logger:    logger,
		lock:      sync.RWMutex{},
		caContent: ca,
	}
}

func (conf *client) getCron(metric string) string {
	if metric == common.MetricPipelineSize {
		return conf.config.DiskCron
	}
	if val, ok := conf.config.MetricCron[metric]; ok {
		level.Debug(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("the cron for metric %s would be %s", metric, val)) //nolint:errcheck

		return val
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("metric %s would be using default cron", metric)) //nolint:errcheck

	return conf.config.APICron
}

func (conf *client) getCronClient() gocd.GoCd {
	return gocd.NewClient(
		conf.config.GoCdBaseURL,
		conf.config.GoCdUserName,
		conf.config.GoCdPassword,
		conf.config.LogLevel,
		conf.caContent,
	)
}
