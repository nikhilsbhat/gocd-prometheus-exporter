package gocd

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetConfigRepoInfo fetches information of all config-repos in GoCd server.
func (conf *Config) GetConfigRepoInfo() ([]ConfigRepo, error) {
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
	return reposConf.ConfigRepos.ConfigRepos, nil
}
