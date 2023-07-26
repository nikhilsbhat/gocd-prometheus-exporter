package gocd

import (
	"sync"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/app"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/sirupsen/logrus"
)

// client holds resty.Client which could be used for interacting with GoCD and other information.
type client struct {
	logger    *logrus.Logger
	config    app.Config
	lock      sync.RWMutex
	caContent []byte
}

// GoCd implements methods to get various information regarding GoCD.
type GoCd interface {
	CronSchedulers()
}

// NewClient returns new instance of client when invoked.
func NewClient(config app.Config, logger *logrus.Logger, ca []byte) GoCd {
	return &client{
		config:    config,
		logger:    logger,
		lock:      sync.RWMutex{},
		caContent: ca,
	}
}

func (conf *client) getCron(metric string) string {
	if val, ok := conf.config.MetricCron[metric]; ok {
		conf.logger.Debugf("the cron for metric %s would be %s", metric, val)

		return val
	}
	conf.logger.Debugf("metric %s would be using default cron", metric)

	return conf.config.APICron
}

func (conf *client) getCronClient() gocd.GoCd {
	auth := gocd.Auth{
		UserName:    conf.config.GoCdUserName,
		Password:    conf.config.GoCdPassword,
		BearerToken: conf.config.GoCDBearerToken,
	}

	return gocd.NewClient(
		conf.config.GoCdBaseURL,
		auth,
		conf.config.LogLevel,
		conf.caContent,
	)
}
