package gocd

import (
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
)

func (conf *client) updateVersionInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	version, err := client.GetVersionInfo()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("version", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentVersion = version
	}

	defer conf.lock.Unlock()
}
