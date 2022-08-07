package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// GetHealthInfo implements method that fetches the details of all warning and errors.
func (conf *client) GetHealthInfo() ([]ServerHealth, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("server health status")) //nolint:errcheck

	var health []ServerHealth
	resp, err := conf.client.R().SetResult(&health).Get(common.GoCdServerHealthEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get health info errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("health status")) //nolint:errcheck
	conf.lock.Unlock()

	return health, nil
}

func (conf *client) updateHealthInfo() {
	healthInfo, err := conf.GetHealthInfo()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("server health", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentServerHealth = healthInfo
	}
}
