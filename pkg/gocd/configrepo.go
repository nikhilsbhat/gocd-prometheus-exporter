package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// GetConfigRepoInfo fetches information of all config-repos in GoCD server.
func (conf *client) GetConfigRepoInfo() ([]ConfigRepo, error) {
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

	return reposConf.ConfigRepos.ConfigRepos, nil
}

func (conf *client) updateConfigRepoInfo() {
	newConf := conf.getCronClient()
	newConf.lock.Lock()
	repos, err := newConf.GetConfigRepoInfo()
	if err != nil {
		level.Error(newConf.logger).Log(common.LogCategoryErr, apiError("config repo", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentConfigRepos = repos
	}
	defer newConf.lock.Unlock()
}
