package gocd

import (
	"fmt"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/thoas/go-funk"

	"github.com/go-kit/log/level"
)

// CronSchedulers schedules all the jobs so that data will be available for the exporter to serve.
func (conf *client) CronSchedulers() {
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("api", conf.defaultAPICron)) //nolint:errcheck
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("disk", conf.diskCron))      //nolint:errcheck

	if !funk.Contains(conf.skipMetrics, common.MetricPipelineSize) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricPipelineSize))
		conf.configureDiskUsage()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricSystemAdminsCount) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricSystemAdminsCount))
		conf.configureAdminsInfo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricConfigRepoCount) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricConfigRepoCount))
		conf.configureGetConfigRepo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricConfiguredBackup) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricConfiguredBackup))
		conf.configureGetBackupInfo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricAgentDown) &&
		!funk.Contains(conf.skipMetrics, common.MetricAgentDiskSpace) &&
		!funk.Contains(conf.skipMetrics, common.MetricAgentsCount) {
		level.Info(conf.logger).Log(common.LogCategoryMsg,
			getCronEnbaled(
				fmt.Sprintf("%s/%s/%s", common.MetricAgentDown, common.MetricAgentDiskSpace, common.MetricAgentsCount)))
		conf.configureGetAgentsInfo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricServerHealth) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricServerHealth))
		conf.configureGetHealthInfo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricPipelineGroupCount) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricPipelineGroupCount))
		conf.configureGetPipelineGroupInfo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricEnvironmentCountAll) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricEnvironmentCountAll))
		conf.configureGetEnvironmentInfo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricVersion) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricVersion))
		conf.configureGetVersionInfo()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricJobStatus) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricJobStatus))
		conf.configureGetAgentJobRunHistory()
	}
}
