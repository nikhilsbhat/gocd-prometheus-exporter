package gocd

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetAdminsInfo fetches information of all system admins present in GoCd server.
func (conf *Config) GetAdminsInfo() (SystemAdmins, error) {
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionTwo,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "trying to retrieve admins information present in GoCd") //nolint:errcheck

	var adminsConf SystemAdmins
	resp, err := conf.client.R().SetResult(&adminsConf).Get(common.GoCdSystemAdminEndpoint)
	if err != nil {
		return SystemAdmins{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return SystemAdmins{}, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "successfully retrieved information of config repos configured in GoCd") //nolint:errcheck
	return adminsConf, nil
}
