package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetPipelineGroupInfo fetches information of backup configured in GoCd server.
func (conf *Config) GetPipelineGroupInfo() ([]PipelineGroup, error) {
	conf.lock.Lock()
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

	conf.lock.Unlock()
	return updatedGroupConf, nil
}

func (conf *Config) configureGetPipelineGroupInfo() {
	scheduleGetPipelineGroupInfo := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetPipelineGroupInfo.AddFunc(conf.otherCron, func() {
		pipelineInfo, err := conf.GetPipelineGroupInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving pipeline group information errored with: %s", err.Error())) //nolint:errcheck
		}
		CurrentPipelineGroup = pipelineInfo
	})

	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetPipelineGroupInfo.Start()
}