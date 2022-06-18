package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetAdminsInfo fetches information of all system admins present in GoCd server.
func (conf *Config) GetAdminsInfo() (SystemAdmins, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionTwo,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("admins")) //nolint:errcheck

	var adminsConf SystemAdmins
	resp, err := conf.client.R().SetResult(&adminsConf).Get(common.GoCdSystemAdminEndpoint)
	if err != nil {
		return SystemAdmins{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return SystemAdmins{}, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("admins")) //nolint:errcheck
	defer conf.lock.Unlock()
	return adminsConf, nil
}

func (conf *Config) configureAdminsInfo() {
	scheduleGetAdmins := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetAdmins.AddFunc(conf.apiCron, func() {
		admins, err := conf.GetAdminsInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, getAPIErrMsg("system admin", err.Error())) //nolint:errcheck
		}
		CurrentSystemAdmins = admins
	})

	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetAdmins.Start()
}
