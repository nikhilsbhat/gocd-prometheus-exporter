package gocd

import (
	"fmt"
	"strings"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func (conf *client) updateAgentsInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	agentsInfo, err := goClient.GetAgents()
	if err != nil {
		conf.logger.Error(apiError("agents", err.Error()))
	}

	if err == nil {
		CurrentAgentsConfig = agentsInfo
	}

	defer conf.lock.Unlock()
}

func (conf *client) updateAgentJobRunHistory() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	agentsJobRunHistory := make([]gocd.AgentJobHistory, 0)
	var errors []string
	for _, agent := range CurrentAgentsConfig {
		agentJobRunHistory, err := goClient.GetAgentJobRunHistory(agent.ID)
		if err != nil {
			conf.logger.Error(apiError("agents", err.Error()))
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
