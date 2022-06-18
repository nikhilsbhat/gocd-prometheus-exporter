package gocd

import (
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
)

// ScheDulers schedules all the jobs so that data will be available for the exporter to serve
func (conf *Config) ScheDulers() {
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("api", conf.apiCron))   //nolint:errcheck
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("disk", conf.diskCron)) //nolint:errcheck
	conf.configureDiskUsage()
	conf.configureAdminsInfo()
	conf.configureGetConfigRepo()
	conf.configureGetBackupInfo()
	conf.configureGetNodesInfo()
	conf.configureGetHealthInfo()
	conf.configureGetPipelineGroupInfo()
}
