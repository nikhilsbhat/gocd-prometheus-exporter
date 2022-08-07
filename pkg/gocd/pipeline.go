package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// GetPipelineGroupInfo fetches information of backup configured in GoCD server.
func (conf *client) GetPipelineGroupInfo() ([]PipelineGroup, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("pipeline groups")) //nolint:errcheck

	var groupConf PipelineGroupsConfig
	resp, err := conf.client.R().SetResult(&groupConf).Get(common.GoCdPipelineGroupEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get pipeline group information errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("pipeline groups")) //nolint:errcheck

	updatedGroupConf := make([]PipelineGroup, 0)
	for _, group := range groupConf.PipelineGroups.PipelineGroups {
		updatedGroupConf = append(updatedGroupConf, PipelineGroup{
			Name:          group.Name,
			PipelineCount: len(group.Pipelines),
		})
	}

	conf.lock.Unlock()

	return updatedGroupConf, nil
}

func (conf *client) getPipelineCount(groups []PipelineGroup) int {
	var pipelines int
	for _, i := range groups {
		pipelines += i.PipelineCount
	}

	return pipelines
}

func (conf *client) updatePipelineGroupInfo() {
	pipelineInfo, err := conf.GetPipelineGroupInfo()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("pipeline group", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentPipelineCount = conf.getPipelineCount(pipelineInfo)
		CurrentPipelineGroup = pipelineInfo
	}
}
