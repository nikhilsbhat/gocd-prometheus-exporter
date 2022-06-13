package gocd

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetHealthInfo implements method that fetches the details of all warning and errors
func (conf *Config) GetHealthInfo() ([]ServerHealth, error) {
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "trying to retrieve GoCd server health status")

	var health []ServerHealth
	resp, err := conf.client.R().SetResult(&health).Get(common.GoCdServerHealthEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, "successfully retrieved GoCd server health status")
	return health, nil
}
