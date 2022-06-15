package gocd

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetPipelineGroupInfo fetches information of backup configured in GoCd server.
func (conf *Config) GetPipelineGroupInfo() ([]PipelineGroup, error) {
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "trying to retrieve pipeline groups from GoCd") //nolint:errcheck

	var groupConf PipelineGroupsConfig
	resp, err := conf.client.R().SetResult(&groupConf).Get(common.GoCdPipelineGroupEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "successfully retrieved pipeline groups information from GoCd") //nolint:errcheck

	updatedGroupConf := make([]PipelineGroup, 0)
	for _, group := range groupConf.PipelineGroups.PipelineGroups {
		updatedGroupConf = append(updatedGroupConf, PipelineGroup{
			Name:          group.Name,
			PipelineCount: len(group.Pipelines),
		})
	}
	return updatedGroupConf, nil
}
