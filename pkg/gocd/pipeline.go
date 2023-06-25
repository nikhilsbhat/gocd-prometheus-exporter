package gocd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/nikhilsbhat/gocd-sdk-go"
)

type PipelineRunHistory struct {
	NotRunIn float64
	Usage    gocd.PipelineSchedules
}

func getPipelineCount(groups []gocd.PipelineGroup) int {
	var pipelines int
	for _, i := range groups {
		pipelines += i.PipelineCount
	}

	return pipelines
}

func (conf *client) updatePipelineGroupInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	pipelineInfo, err := goClient.GetPipelineGroups()
	if err != nil {
		conf.logger.Error(apiError("pipeline group", err.Error()))
	}

	if err == nil {
		CurrentPipelineCount = getPipelineCount(pipelineInfo)
		CurrentPipelineGroup = pipelineInfo
	}

	defer conf.lock.Unlock()
}

func (conf *client) updatePipelines() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	pipelinesInfo, err := goClient.GetPipelines()
	if err != nil {
		conf.logger.Error(apiError(common.MetricPipelines, err.Error()))
	}

	if err == nil {
		pipelines := make([]string, 0)
		for _, pipeline := range pipelinesInfo.Pipeline {
			pipelineName, err := gocd.GetPipelineName(pipeline.Href)
			if err != nil {
				conf.logger.Error(err.Error())
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
	goClient := conf.getCronClient()

	pipelinesStatus := make([]gocd.PipelineState, 0)
	for _, pipeline := range CurrentPipelines {
		pipelineStatus, err := goClient.GetPipelineState(pipeline)
		if err != nil {
			conf.logger.Error(apiError(common.MetricPipelineState, err.Error()))
		}
		pipelinesStatus = append(pipelinesStatus, pipelineStatus)
	}

	CurrentPipelineState = pipelinesStatus

	defer conf.lock.Unlock()
}

func (conf *client) updatePipelineRunInLastXDays() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	pipelinesStatus := make([]PipelineRunHistory, 0)
	for _, pipeline := range CurrentPipelines {
		pipelineStatus, err := goClient.GetPipelineSchedules(pipeline, "0", "1")
		if err != nil {
			conf.logger.Error(apiError(common.MetricPipelineState, err.Error()))
		}

		if pipelineStatus.Groups[0].History[0].ScheduledDate == "N/A" {
			continue
		}

		timeThen := time.UnixMilli(pipelineStatus.Groups[0].History[0].ScheduledTimestamp).UTC()
		timeNow := time.Now()

		timeDiff := timeNow.Sub(timeThen).Round(1).Hours()

		skipDays := os.Getenv("GOCD_PIPELINE_DAYS_TO_SKIP")
		if len(skipDays) != 0 {
			const bitSize = 64
			s, err := strconv.ParseFloat(skipDays, bitSize)
			if err == nil {
				fmt.Println(s)
			}

			if timeDiff < s {
				continue
			}
		}

		const hoursInaDay = 24
		pipelinesStatus = append(pipelinesStatus, PipelineRunHistory{
			NotRunIn: timeDiff / hoursInaDay,
			Usage:    pipelineStatus,
		})
	}

	CurrentPipelineNotRun = pipelinesStatus

	defer conf.lock.Unlock()
}
