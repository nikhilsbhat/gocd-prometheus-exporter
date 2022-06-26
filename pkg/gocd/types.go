package gocd

var (
	// CurrentAgentsConfig holds updated GoCD agent config information.
	CurrentAgentsConfig []Agent
	// CurrentAgentJobRunHistory holds updated GoCD agent config information.
	CurrentAgentJobRunHistory = make([]AgentJobHistory, 0)
	// CurrentServerHealth holds updated GoCD server health information.
	CurrentServerHealth []ServerHealth
	// CurrentConfigRepos holds updated config repo information present GoCD server.
	CurrentConfigRepos []ConfigRepo
	// CurrentPipelineGroup holds updated pipeline group information present GoCD server.
	CurrentPipelineGroup []PipelineGroup
	// CurrentEnvironments holds updated environment information present GoCD server.
	CurrentEnvironments []Environment
	// CurrentPipelineCount holds updated pipeline count present in GoCD server.
	CurrentPipelineCount int
	// CurrentSystemAdmins holds updated information of the system admins present in GoCD server.
	CurrentSystemAdmins SystemAdmins
	// CurrentBackupConfig holds updated information of backups configured in GoCD server.
	CurrentBackupConfig BackupConfig
	// CurrentVersion holds updated GoCD server version information.
	CurrentVersion VersionInfo
	// CurrentPipelineSize holds updated information of disk sizes occupied by various GoCD pipelines.
	CurrentPipelineSize = make(map[string]PipelineSize)
)

const (
	defaultRetryCount    = 5
	defaultRetryWaitTime = 5
)

// AgentsConfig holds information of all agent of GoCD.
type AgentsConfig struct {
	Config Agents `json:"_embedded,omitempty"`
}

// Agents holds information of all agent of GoCD.
type Agents struct {
	Config []Agent `json:"agents,omitempty"`
}

// Agent holds information of a particular agent.
type Agent struct {
	Name               string      `json:"hostname,omitempty"`
	ID                 string      `json:"uuid,omitempty"`
	Version            string      `json:"agent_version,omitempty"`
	CurrentState       string      `json:"agent_state,omitempty"`
	OS                 string      `json:"operating_system,omitempty"`
	ConfigState        string      `json:"agent_config_state,omitempty"`
	Sandbox            string      `json:"sandbox,omitempty"`
	DiskSpaceAvailable interface{} `json:"free_space,omitempty"`
}

// ServerVersion holds version information GoCd server.
type ServerVersion struct {
	Version     string `json:"version,omitempty"`
	GitSha      string `json:"git_sha,omitempty"`
	FullVersion string `json:"full_version,omitempty"`
	CommitURL   string `json:"commit_url,omitempty"`
}

// ServerHealth holds information of GoCD server health.
type ServerHealth struct {
	Level   string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
}

// ConfigRepoConfig holds information of all config-repos present in GoCD.
type ConfigRepoConfig struct {
	ConfigRepos ConfigRepos `json:"_embedded,omitempty"`
}

// ConfigRepos holds information of all config-repos present in GoCD.
type ConfigRepos struct {
	ConfigRepos []ConfigRepo `json:"config_repos,omitempty"`
}

// ConfigRepo holds information of the specified config-repo.
type ConfigRepo struct {
	ID       string `json:"config_repos,omitempty"`
	Material struct {
		Type       string `json:"type,omitempty"`
		Attributes struct {
			URL        string `json:"url,omitempty"`
			Branch     string `json:"branch,omitempty"`
			AutoUpdate bool   `json:"auto_update,omitempty"`
		}
	}
}

// PipelineGroupsConfig holds information on the various pipeline groups present in GoCD.
type PipelineGroupsConfig struct {
	PipelineGroups PipelineGroups `json:"_embedded,omitempty"`
}

// PipelineGroups holds information on the various pipeline groups present in GoCD.
type PipelineGroups struct {
	PipelineGroups []PipelineGroup `json:"groups,omitempty"`
}

// PipelineGroup holds information of a specific pipeline group instance.
type PipelineGroup struct {
	Name          string `json:"name,omitempty"`
	PipelineCount int    `json:"pipeline_count,omitempty"`
	Pipelines     []struct {
		Name string `json:"name,omitempty"`
	}
}

// SystemAdmins holds information of the system admins present.
type SystemAdmins struct {
	Roles []string `json:"roles,omitempty"`
	Users []string `json:"users,omitempty"`
}

// BackupConfig holds information of the backup configured.
type BackupConfig struct {
	EmailOnSuccess bool   `json:"email_on_success,omitempty"`
	EmailOnFailure bool   `json:"email_on_failure,omitempty"`
	Schedule       string `json:"schedule,omitempty"`
}

// PipelineSize holds information of the pipeline size.
type PipelineSize struct {
	Size float64
	Type string
}

// Pipelines holds information of the pipelines present in GoCD.
type Pipelines struct {
	Pipelines []Pipeline `json:"pipelines,omitempty"`
}

// Pipeline holds information of a specific pipeline instance.
type Pipeline struct {
	Name string `json:"name,omitempty"`
}

// EnvironmentInfo holds information of all environments present in GoCD.
type EnvironmentInfo struct {
	Environments Environments `json:"_embedded,omitempty"`
}

// Environments holds information of all environments present in GoCD.
type Environments struct {
	Environments []Environment `json:"environments,omitempty"`
}

// Environment holds information of a specific environment present in GoCD.
type Environment struct {
	Name      string     `json:"name,omitempty"`
	Pipelines []Pipeline `json:"pipelines,omitempty"`
}

// VersionInfo holds version information of GoCD server.
type VersionInfo struct {
	Version     string `json:"version,omitempty"`
	FullVersion string `json:"full_version,omitempty"`
	GitSHA      string `json:"git_sha,omitempty"`
}

// AgentJobHistory holds information of pipeline run history of all GoCD agents.
type AgentJobHistory struct {
	Jobs       []JobRunHistory `json:"jobs,omitempty"`
	Pagination Pagination      `json:"pagination"`
}

// JobRunHistory holds information of pipeline run history of a specific GoCD agent.
type JobRunHistory struct {
	Name            string `json:"pipeline_name,omitempty"`
	JobName         string `json:"job_name,omitempty"`
	StageName       string `json:"stage_name,omitempty"`
	StageCounter    int64  `json:"stage_counter,string,omitempty"`
	PipelineCounter int64  `json:"pipeline_counter,omitempty"`
	Result          string `json:"result,omitempty"`
}

// Pagination holds information which is helpful in paginating the results of job run history.
type Pagination struct {
	PageSize int64 `json:"page_size,omitempty"`
	Offset   int64 `json:"offset,omitempty"`
	Total    int64 `json:"total,omitempty"`
}
