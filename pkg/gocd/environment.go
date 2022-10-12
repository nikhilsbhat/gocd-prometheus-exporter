package gocd

import (
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
)

func (conf *client) updateEnvironmentInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	environmentInfo, err := client.GetEnvironments()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("environment", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentEnvironments = environmentInfo
	}

	defer conf.lock.Unlock()
}
