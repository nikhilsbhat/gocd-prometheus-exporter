package gocd

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetNodesInfo implements method that fetches the details of all the agents present in GoCd server
func (conf *Config) GetNodesInfo() (NodesConfig, error) {
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionSeven,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "trying to retrieve nodes information present in GoCd") //nolint:errcheck

	var nodesConf NodesConfig
	resp, err := conf.client.R().SetResult(&nodesConf).Get(common.GoCdAgentsEndpoint)
	if err != nil {
		return NodesConfig{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return NodesConfig{}, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "successfully retrieved nodes information from GoCd") //nolint:errcheck
	return nodesConf, nil
}
