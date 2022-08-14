package gocd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// GetPipelineGroupInfo fetches information of backup configured in GoCD server.
func (conf *client) GetPipelineGroupInfo() ([]PipelineGroup, error) {
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

	return updatedGroupConf, nil
}

// GetPipelines fetches all pipelines configured in GoCD server.
func (conf *client) GetPipelines() (PipelinesInfo, error) {
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages(common.MetricPipelines)) //nolint:errcheck

	var pipelinesInfo PipelinesInfo
	resp, err := conf.client.R().SetResult(&pipelinesInfo).Get(common.GoCdAPIFeedPipelineEndpoint)
	if err != nil {
		return PipelinesInfo{}, fmt.Errorf("call made to get pipelines errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelinesInfo{}, apiWithCodeError(resp.StatusCode())
	}
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages(common.MetricPipelines)) //nolint:errcheck

	return pipelinesInfo, nil
}

func (conf *client) GetPipelineState() ([]PipelineState, error) {
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages(common.MetricPipelineState)) //nolint:errcheck

	pipelinesStatus := make([]PipelineState, 0)
	for _, pipeline := range CurrentPipelines {
		var pipelineState PipelineState
		resp, err := conf.client.R().SetResult(&pipelineState).Get(fmt.Sprintf(common.GoCdPipelineStatus, pipeline))
		if err != nil {
			return nil, fmt.Errorf("call made to get pipeline state errored with %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			return nil, apiWithCodeError(resp.StatusCode())
		}

		pipelineState.Name = pipeline
		pipelinesStatus = append(pipelinesStatus, pipelineState)
	}

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages(common.MetricJobStatus)) //nolint:errcheck

	return pipelinesStatus, nil
}

func (conf *client) getPipelineCount(groups []PipelineGroup) int {
	var pipelines int
	for _, i := range groups {
		pipelines += i.PipelineCount
	}

	return pipelines
}

func GetPipelineName(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("parsing URL errored with %w", err)
	}

	return strings.TrimSuffix(strings.TrimPrefix(parsedURL.Path, "/go/api/feed/pipelines/"), "/stages.xml"), nil
}

func (conf *client) updatePipelineGroupInfo() {
	newConf := conf.getCronClient()
	newConf.lock.Lock()
	pipelineInfo, err := newConf.GetPipelineGroupInfo()
	if err != nil {
		level.Error(newConf.logger).Log(common.LogCategoryErr, apiError("pipeline group", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentPipelineCount = newConf.getPipelineCount(pipelineInfo)
		CurrentPipelineGroup = pipelineInfo
	}
	defer newConf.lock.Unlock()
}

func (conf *client) updatePipelines() {
	newConf := conf.getCronClient()
	newConf.lock.Lock()
	pipelinesInfo, err := newConf.GetPipelines()
	if err != nil {
		level.Error(newConf.logger).Log(common.LogCategoryErr, apiError(common.MetricPipelines, err.Error())) //nolint:errcheck
	}

	fmt.Println("pipelinesInfo: ", pipelinesInfo)
	if err == nil {
		var pipelines []string
		for _, pipeline := range pipelinesInfo.Pipeline {
			pipelineName, err := GetPipelineName(pipeline.Href)
			if err != nil {
				level.Error(newConf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
			} else {
				pipelines = append(pipelines, pipelineName)
			}
		}
		fmt.Println("pipelines: ", pipelines)
		CurrentPipelines = pipelines
	}
	defer newConf.lock.Unlock()
}

func (conf *client) updatePipelineState() {
	newConf := conf.getCronClient()
	newConf.lock.Lock()
	pipelineStatus, err := newConf.GetPipelineState()
	if err != nil {
		level.Error(newConf.logger).Log(common.LogCategoryErr, apiError(common.MetricPipelineState, err.Error())) //nolint:errcheck
	}
	if err == nil {
		fmt.Println(pipelineStatus)
		CurrentPipelineState = pipelineStatus
	}
	defer newConf.lock.Unlock()
}
