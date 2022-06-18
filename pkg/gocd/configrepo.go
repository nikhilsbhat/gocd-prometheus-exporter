package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetConfigRepoInfo fetches information of all config-repos in GoCd server.
func (conf *Config) GetConfigRepoInfo() ([]ConfigRepo, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionFour,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "trying to retrieve config repos present in GoCd") //nolint:errcheck

	var reposConf ConfigRepoConfig
	resp, err := conf.client.R().SetResult(&reposConf).Get(common.GoCdConfigReposEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "successfully retrieved information of config repos configured in GoCd") //nolint:errcheck
	conf.lock.Unlock()
	return reposConf.ConfigRepos.ConfigRepos, nil
}

func (conf *Config) configureGetConfigRepo() {
	scheduleGetConfigRepo := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetConfigRepo.AddFunc(conf.otherCron, func() {
		repos, err := conf.GetConfigRepoInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving config repo information errored with: %s", err.Error())) //nolint:errcheck
		}
		CurrentConfigRepos = repos
	})

	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetConfigRepo.Start()
}
