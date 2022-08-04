package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetConfigRepoInfo fetches information of all config-repos in GoCD server.
func (conf *client) GetConfigRepoInfo() ([]ConfigRepo, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionFour,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("config repos")) //nolint:errcheck

	var reposConf ConfigRepoConfig
	resp, err := conf.client.R().SetResult(&reposConf).Get(common.GoCdConfigReposEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get config repo errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("config repos")) //nolint:errcheck
	conf.lock.Unlock()

	return reposConf.ConfigRepos.ConfigRepos, nil
}

func (conf *client) configureGetConfigRepo() {
	scheduleGetConfigRepo := cron.New(cron.WithLogger(getCronLogger(common.MetricConfigRepoCount)), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetConfigRepo.AddFunc(conf.getCron(common.MetricConfigRepoCount), func() {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricConfigRepoCount)) //nolint:errcheck

		repos, err := conf.GetConfigRepoInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, apiError("config repo", err.Error())) //nolint:errcheck
		}
		CurrentConfigRepos = repos
	})
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetConfigRepo.Start()
}
