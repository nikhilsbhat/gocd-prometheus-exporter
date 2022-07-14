package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
)

// GetAgentsInfo implements method that fetches the details of all the agents present in GoCD server.
func (conf *client) GetAgentsInfo() ([]Agent, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionSeven,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("agents")) //nolint:errcheck

	var agentsConf AgentsConfig
	resp, err := conf.client.R().SetResult(&agentsConf).Get(common.GoCdAgentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get agents information errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}

	conf.lock.Unlock()
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("agents")) //nolint:errcheck

	return agentsConf.Config.Config, nil
}

func (conf *client) GetAgentJobRunHistory() ([]AgentJobHistory, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	conf.client.SetQueryParam("sort_order", "DESC")

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("agent job run history")) //nolint:errcheck

	jobHistory := make([]AgentJobHistory, 0)
	for _, currentAgent := range CurrentAgentsConfig {
		var jobHistoryConf AgentJobHistory
		resp, err := conf.client.R().SetResult(&jobHistoryConf).Get(fmt.Sprintf(common.GoCdJobRunHistoryEndpoint, currentAgent.ID))
		level.Debug(conf.logger).Log(common.LogCategoryMsg, resp.Request.URL) //nolint:errcheck

		if err != nil {
			return nil, fmt.Errorf("call made to get agent job run history errored with %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			return nil, apiWithCodeError(resp.StatusCode())
		}
		jobHistory = append(jobHistory, jobHistoryConf)
	}

	if len(jobHistory) == 0 {
		level.Debug(conf.logger).Log(common.LogCategoryMsg, "no history found") //nolint:errcheck
	}

	conf.lock.Unlock()
	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("agent job run history")) //nolint:errcheck

	return jobHistory, nil
}

func (conf *client) configureGetAgentsInfo() {
	scheduleGetAgentsInfo := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetAgentsInfo.AddFunc(conf.getCron(common.MetricAgentDown), func() {
		agentsInfo, err := conf.GetAgentsInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, apiError("agents", err.Error())) //nolint:errcheck
		}
		CurrentAgentsConfig = agentsInfo
	})
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetAgentsInfo.Start()
}

func (conf *client) configureGetAgentJobRunHistory() {
	scheduleGetAgentJobRunHistory := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetAgentJobRunHistory.AddFunc(conf.getCron(common.MetricJobStatus), func() {
		agentsJobRunHistory, err := conf.GetAgentJobRunHistory()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, apiError("agents", err.Error())) //nolint:errcheck
		}
		CurrentAgentJobRunHistory = agentsJobRunHistory
	})
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetAgentJobRunHistory.Start()
}
