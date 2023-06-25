package gocd

import (
	"github.com/nikhilsbhat/gocd-sdk-go"
)

var (
	// CurrentAgentsConfig holds updated GoCD agent config information.
	CurrentAgentsConfig []gocd.Agent
	// CurrentAgentJobRunHistory holds updated GoCD agent config information.
	CurrentAgentJobRunHistory = make([]gocd.AgentJobHistory, 0)
	// CurrentServerHealth holds updated GoCD server health information.
	CurrentServerHealth []gocd.ServerHealth
	// CurrentConfigRepos holds updated config repo information present GoCD server.
	CurrentConfigRepos []gocd.ConfigRepo
	// CurrentPipelineGroup holds updated pipeline group information present GoCD server.
	CurrentPipelineGroup []gocd.PipelineGroup
	// CurrentEnvironments holds updated environment information present GoCD server.
	CurrentEnvironments []gocd.Environment
	// CurrentPipelineCount holds updated pipeline count present in GoCD server.
	CurrentPipelineCount int
	// CurrentSystemAdmins holds updated information of the system admins present in GoCD server.
	CurrentSystemAdmins gocd.SystemAdmins
	// CurrentBackupConfig holds updated information of backups configured in GoCD server.
	CurrentBackupConfig gocd.BackupConfig
	// CurrentVersion holds updated GoCD server version information.
	CurrentVersion gocd.VersionInfo
	// CurrentPipelines holds updated list of pipeline names that are present in GoCD.
	CurrentPipelines []string
	// CurrentPipelineState holds the information of the latest state of pipelines available in GoCD.
	CurrentPipelineState []gocd.PipelineState
	// CurrentElasticProfileUsage holds the information of the pipelines using various elastic agent profiles.
	CurrentElasticProfileUsage []ElasticProfileUsage
	// CurrentPluginInfo holds the information of the plugins installed in GoCD.
	CurrentPluginInfo []gocd.Plugin
	// CurrentPipelineNotRun holds the information of the pipelines not run in last X days.
	CurrentPipelineNotRun []PipelineRunHistory
)
