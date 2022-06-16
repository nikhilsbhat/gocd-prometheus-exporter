package gocd

import (
	"fmt"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// ScheDulers schedules all the jobs so that data will be available for the exporter to serve
func (conf *Config) ScheDulers() {
	level.Info(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("otherCron will be scheduled for %s as specified", conf.otherCron)) //nolint:errcheck
	level.Info(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("diskCron will be scheduled for %s as specified", conf.diskCron))   //nolint:errcheck
	conf.configureDiskUsage()
	conf.configureAdminsInfo()
	conf.configureGetConfigRepo()
	conf.configureGetBackupInfo()
	conf.configureGetNodesInfo()
	conf.configureGetHealthInfo()
	conf.configureGetPipelineGroupInfo()
}
