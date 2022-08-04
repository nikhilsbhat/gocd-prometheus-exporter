package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetBackupInfo fetches information of backup configured in GoCD server.
func (conf *client) GetBackupInfo() (BackupConfig, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("backup")) //nolint:errcheck

	var backUpConf BackupConfig
	resp, err := conf.client.R().SetResult(&backUpConf).Get(common.GoCdBackupConfigEndpoint)
	if err != nil {
		return BackupConfig{}, fmt.Errorf("call made to get backup information errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return BackupConfig{}, apiWithCodeError(resp.StatusCode())
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("backup")) //nolint:errcheck
	defer conf.lock.Unlock()

	return backUpConf, nil
}

func (conf *client) configureGetBackupInfo() {
	scheduleGetBackupInfo := cron.New(cron.WithLogger(getCronLogger(common.MetricConfiguredBackup)), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetBackupInfo.AddFunc(conf.getCron(common.MetricConfiguredBackup), func() {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricConfiguredBackup)) //nolint:errcheck

		backupInfo, err := conf.GetBackupInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, apiError("gocd backup", err.Error())) //nolint:errcheck
		}
		CurrentBackupConfig = backupInfo
	})
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}

	scheduleGetBackupInfo.Start()
}
