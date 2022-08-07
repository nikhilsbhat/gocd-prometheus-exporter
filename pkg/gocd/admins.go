package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// GetAdminsInfo fetches information of all system admins present in GoCD server.
func (conf *client) GetAdminsInfo() (SystemAdmins, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionTwo,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("admins")) //nolint:errcheck

	var adminsConf SystemAdmins
	resp, err := conf.client.R().SetResult(&adminsConf).Get(common.GoCdSystemAdminEndpoint)
	if err != nil {
		return SystemAdmins{}, fmt.Errorf("call made to get system admin errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return SystemAdmins{}, apiWithCodeError(resp.StatusCode())
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("admins")) //nolint:errcheck
	defer conf.lock.Unlock()

	return adminsConf, nil
}

func (conf *client) updateAdminsInfo() {
	admins, err := conf.GetAdminsInfo()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("system admin", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentSystemAdmins = admins
	}
}
