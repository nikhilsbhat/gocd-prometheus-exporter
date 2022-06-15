package common

const (
	Namespace              = "gocd"
	TypeLink               = "link"
	TypeDir                = "dir"
	LogCategoryMsg         = "msg"
	LogCategoryErr         = "err"
	GoCdDisconnectedState  = "LostContact"
	GoCdReturnErrorMessage = `gocd server returned code %d with message`
	ExporterConfigFileName = "gocd-prometheus-exporter"
	ExporterConfigFileExt  = "yaml"
)

const (
	GoCdAgentsEndpoint        = "/api/agents"
	GoCdVersionEndpoint       = "/api/version"
	GoCdServerHealthEndpoint  = "/api/server_health_messages"
	GoCdConfigReposEndpoint   = "/api/admin/config_repos"
	GoCdSystemAdminEndpoint   = "/api/admin/security/system_admins"
	GoCdBackupConfigEndpoint  = "/api/config/backup"
	GoCdPipelineGroupEndpoint = "/api/admin/pipeline_groups"
	GoCdHeaderVersionSeven    = "application/vnd.go.cd.v7+json"
	GoCdHeaderVersionOne      = "application/vnd.go.cd.v1+json"
	GoCdHeaderVersionFour     = "application/vnd.go.cd.v4+json"
	GoCdHeaderVersionTwo      = "application/vnd.go.cd.v2+json"
)
const (
	MetricAgentsCount        = "agents_count"
	MetricAgentDiskSpace     = "agent_disk_space"
	MetricAgentDown          = "agent_down"
	MetricPipelineSize       = "pipeline_size"
	MetricServerHealth       = "server_health"
	MetricConfigRepoCount    = "config_repo_count"
	MetricPipelineGroupCount = "pipeline_group_count"
	MetricConfiguredBackup   = "backup_configured"
	MetricSystemAdminsCount  = "admin_count"
)

func Float(value interface{}) float64 {
	switch value.(type) {
	case int64:
		return value.(float64)
	case string:
		return float64(0)
	default:
		return value.(float64)
	}
}
