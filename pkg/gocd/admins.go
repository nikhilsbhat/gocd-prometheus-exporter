package gocd

import (
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
)

func (conf *client) updateAdminsInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	admins, err := client.GetSystemAdmins()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("system admin", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentSystemAdmins = admins
	}

	defer conf.lock.Unlock()
}
