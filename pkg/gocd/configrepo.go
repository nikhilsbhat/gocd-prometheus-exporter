package gocd

import (
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
)

func (conf *client) updateConfigRepoInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	repos, err := client.GetConfigRepos()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("config repo", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentConfigRepos = repos
	}

	defer conf.lock.Unlock()
}
