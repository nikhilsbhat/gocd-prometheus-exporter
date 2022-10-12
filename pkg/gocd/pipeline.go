package gocd

import (
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/nikhilsbhat/gocd-sdk-go"
)

func getPipelineCount(groups []gocd.PipelineGroup) int {
	var pipelines int
	for _, i := range groups {
		pipelines += i.PipelineCount
	}

	return pipelines
}

func (conf *client) updatePipelineGroupInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	pipelineInfo, err := client.GetPipelineGroups()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("pipeline group", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentPipelineCount = getPipelineCount(pipelineInfo)
		CurrentPipelineGroup = pipelineInfo
	}

	defer conf.lock.Unlock()
}

func (conf *client) updatePipelines() {
	conf.lock.Lock()
	client := conf.getCronClient()

	pipelinesInfo, err := client.GetPipelines()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError(common.MetricPipelines, err.Error())) //nolint:errcheck
	}

	if err == nil {
		pipelines := make([]string, 0)
		for _, pipeline := range pipelinesInfo.Pipeline {
			pipelineName, err := gocd.GetPipelineName(pipeline.Href)
			if err != nil {
				level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
			} else {
				pipelines = append(pipelines, pipelineName)
			}
		}
		CurrentPipelines = pipelines
	}
	defer conf.lock.Unlock()
}

func (conf *client) updatePipelineState() {
	conf.lock.Lock()
	client := conf.getCronClient()

	pipelinesStatus := make([]gocd.PipelineState, 0)
	for _, pipeline := range CurrentPipelines {
		pipelineStatus, err := client.GetPipelineState(pipeline)
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, apiError(common.MetricPipelineState, err.Error())) //nolint:errcheck
		}
		pipelinesStatus = append(pipelinesStatus, pipelineStatus)
	}

	CurrentPipelineState = pipelinesStatus

	defer conf.lock.Unlock()
}
