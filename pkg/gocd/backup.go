package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetBackupInfo fetches information of backup configured in GoCd server.
func (conf *Config) GetBackupInfo() (BackupConfig, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("backup")) //nolint:errcheck

	var backUpConf BackupConfig
	resp, err := conf.client.R().SetResult(&backUpConf).Get(common.GoCdBackupConfigEndpoint)
	if err != nil {
		return BackupConfig{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return BackupConfig{}, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("backup")) //nolint:errcheck
	defer conf.lock.Unlock()
	return backUpConf, nil
}

func (conf *Config) configureGetBackupInfo() {
	scheduleGetBackupInfo := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetBackupInfo.AddFunc(conf.apiCron, func() {
		backupInfo, err := conf.GetBackupInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, getAPIErrMsg("gocd backup", err.Error())) //nolint:errcheck
		}
		CurrentBackupConfig = backupInfo
	})

	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}

	scheduleGetBackupInfo.Start()
}
