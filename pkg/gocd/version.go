package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
)

// GetVersionInfo fetches version information of the GoCD to which it is connected to.
func (conf *client) GetVersionInfo() (VersionInfo, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("version")) //nolint:errcheck

	var version VersionInfo
	resp, err := conf.client.R().Get(common.GoCdVersionEndpoint)
	if err != nil {
		return VersionInfo{}, fmt.Errorf("call made to get version information errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return VersionInfo{}, apiWithCodeError(resp.StatusCode())
	}
	if err := json.Unmarshal(resp.Body(), &version); err != nil {
		return VersionInfo{}, responseReadError(err.Error())
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("version")) //nolint:errcheck
	conf.lock.Unlock()

	return version, nil
}

func (conf *client) updateVersionInfo() {
	version, err := conf.GetVersionInfo()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("version", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentVersion = version
	}
}
