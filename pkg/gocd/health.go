package gocd

import (
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
)

func (conf *client) updateHealthInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	healthInfo, err := client.GetServerHealthMessages()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("server health", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentServerHealth = healthInfo
	}

	defer conf.lock.Unlock()
}
