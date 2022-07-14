package gocd

import (
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// CronSchedulers schedules all the jobs so that data will be available for the exporter to serve.
func (conf *client) CronSchedulers() {
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("api", conf.defaultAPICron)) //nolint:errcheck
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("disk", conf.diskCron))      //nolint:errcheck
	conf.configureDiskUsage()
	conf.configureAdminsInfo()
	conf.configureGetConfigRepo()
	conf.configureGetBackupInfo()
	conf.configureGetAgentsInfo()
	conf.configureGetHealthInfo()
	conf.configureGetPipelineGroupInfo()
	conf.configureGetEnvironmentInfo()
	conf.configureGetVersionInfo()
	conf.configureGetAgentJobRunHistory()
}
