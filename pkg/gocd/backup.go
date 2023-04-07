package gocd

import (
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

func (conf *client) updateBackupInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	backupInfo, err := client.GetBackupConfig()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("gocd backup", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentBackupConfig = backupInfo
	}

	defer conf.lock.Unlock()
}
