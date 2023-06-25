package common

const (
	Namespace             = "gocd"
	GoCdDisconnectedState = "LostContact"
	GoCDExporterName      = "gocd-prometheus-exporter"
	ExporterConfigFileExt = "yaml"
)

const (
	MetricAgentsCount              = "agents_count"
	MetricAgentDiskSpace           = "agent_disk_space"
	MetricAgentDown                = "agent_down"
	MetricPipelineSize             = "pipeline_size"
	MetricServerHealth             = "server_health"
	MetricConfigRepoCount          = "config_repo_count"
	MetricPipelineGroupCount       = "pipeline_group_count"
	MetricPipelineCount            = "pipeline_count"
	MetricConfiguredBackup         = "backup_configured"
	MetricSystemAdminsCount        = "admin_count"
	MetricEnvironmentCountAll      = "environment_count_all"
	MetricVersion                  = "version"
	MetricJobStatus                = "pipeline_status"
	MetricPipelines                = "pipelines"
	MetricPipelineState            = "pipeline_state"
	MetricElasticAgentProfileUsage = "elastic_agent_profile_usage"
	MetricPlugins                  = "plugins"
)

const (
	GoCdPipelineStatePass = "Passed"
	GoCdPipelineStateFail = "Failed"
)

func Float(value interface{}) float64 {
	switch value.(type) {
	case int64:
		return value.(float64) //nolint:forcetypeassert
	case string:
		return float64(0)
	default:
		return value.(float64) //nolint:forcetypeassert
	}
}
