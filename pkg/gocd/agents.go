package gocd

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// GetAgentsInfo implements method that fetches the details of all the agents present in GoCD server.
func (conf *client) GetAgentsInfo() ([]Agent, error) {
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

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("agents")) //nolint:errcheck

	return agentsConf.Config.Config, nil
}

func (conf *client) GetAgentJobRunHistory() ([]AgentJobHistory, error) {
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionOne,
	})
	conf.client.SetQueryParam("sort_order", "DESC")

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getTryMessages("agent job run history")) //nolint:errcheck

	jobHistory := make([]AgentJobHistory, 0)
	for _, currentAgent := range CurrentAgentsConfig {
		var jobHistoryConf AgentJobHistory
		resp, err := conf.client.R().SetResult(&jobHistoryConf).Get(fmt.Sprintf(common.GoCdJobRunHistoryEndpoint, currentAgent.ID))
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

	level.Debug(conf.logger).Log(common.LogCategoryMsg, getSuccessMessages("agent job run history")) //nolint:errcheck

	return jobHistory, nil
}

func (conf *client) updateAgentsInfo() {
	newConf := conf.getCronClient()
	newConf.lock.Lock()
	agentsInfo, err := newConf.GetAgentsInfo()
	if err != nil {
		level.Error(newConf.logger).Log(common.LogCategoryErr, apiError("agents", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentAgentsConfig = agentsInfo
	}
	defer newConf.lock.Unlock()
}

func (conf *client) updateAgentJobRunHistory() {
	newConf := conf.getCronClient()
	newConf.lock.Lock()
	agentsJobRunHistory, err := newConf.GetAgentJobRunHistory()
	if err != nil {
		level.Error(newConf.logger).Log(common.LogCategoryErr, apiError("agents", err.Error())) //nolint:errcheck
	}
	if err == nil {
		CurrentAgentJobRunHistory = agentsJobRunHistory
	}
	defer newConf.lock.Unlock()
}
