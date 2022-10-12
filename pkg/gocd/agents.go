package gocd

import (
	"fmt"
	"strings"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/nikhilsbhat/gocd-sdk-go"
)

func (conf *client) updateAgentsInfo() {
	conf.lock.Lock()
	client := conf.getCronClient()

	agentsInfo, err := client.GetAgents()
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, apiError("agents", err.Error())) //nolint:errcheck
	}

	if err == nil {
		CurrentAgentsConfig = agentsInfo
	}

	defer conf.lock.Unlock()
}

func (conf *client) updateAgentJobRunHistory() {
	conf.lock.Lock()
	client := conf.getCronClient()

	agentsJobRunHistory := make([]gocd.AgentJobHistory, 0)
	var errors []string
	for _, agent := range CurrentAgentsConfig {
		agentJobRunHistory, err := client.GetAgentJobRunHistory(agent.ID)
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, apiError("agents", err.Error())) //nolint:errcheck
			errors = append(errors, err.Error())
		}
		agentsJobRunHistory = append(agentsJobRunHistory, agentJobRunHistory)
	}

	if len(errors) == 0 {
		CurrentAgentJobRunHistory = agentsJobRunHistory
	} else {
		fmt.Println(strings.Join(errors, ","))
	}

	defer conf.lock.Unlock()
}
