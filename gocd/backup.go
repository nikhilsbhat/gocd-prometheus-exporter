package gocd

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetBackupInfo fetches information of backup configured in GoCd server.
func (conf *Config) GetBackupInfo() (BackupConfig, error) {
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "trying to retrieve backup configurations present in GoCd") //nolint:errcheck

	var backUpConf BackupConfig
	resp, err := conf.client.R().SetResult(&backUpConf).Get(common.GoCdBackupConfigEndpoint)
	if err != nil {
		return BackupConfig{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return BackupConfig{}, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "successfully retrieved information of backup configurations in GoCd") //nolint:errcheck
	return backUpConf, nil
}
