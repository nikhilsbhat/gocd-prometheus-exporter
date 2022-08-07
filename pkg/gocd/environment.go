package gocd

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
)

// GetEnvironmentInfo fetches information of backup configured in GoCD server.
func (conf *client) GetEnvironmentInfo() ([]Environment, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionThree,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("environment")) //nolint:errcheck

	var envConf EnvironmentInfo
	resp, err := conf.client.R().SetResult(&envConf).Get(common.GoCdEnvironmentEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get environment errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("environment")) //nolint:errcheck

	conf.lock.Unlock()

	return envConf.Environments.Environments, nil
}

func (conf *client) updateEnvironmentInfo() {
	environmentInfo, err := conf.GetEnvironmentInfo()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("environment", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentEnvironments = environmentInfo
	}
}
