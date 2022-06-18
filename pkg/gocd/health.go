package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetHealthInfo implements method that fetches the details of all warning and errors
func (conf *Config) GetHealthInfo() ([]ServerHealth, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("server health status")) //nolint:errcheck

	var health []ServerHealth
	resp, err := conf.client.R().SetResult(&health).Get(common.GoCdServerHealthEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("health status")) //nolint:errcheck
	conf.lock.Unlock()
	return health, nil
}

func (conf *Config) configureGetHealthInfo() {
	scheduleGetHealthInfo := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetHealthInfo.AddFunc(conf.apiCron, func() {
		healthInfo, err := conf.GetHealthInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, getAPIErrMsg("server health", err.Error())) //nolint:errcheck
		}
		CurrentServerHealth = healthInfo
	})

	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetHealthInfo.Start()
}
