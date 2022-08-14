package gocd

import (
	"fmt"
	"log"
	"time"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/thoas/go-funk"

	"github.com/go-co-op/gocron"
	"github.com/go-kit/log/level"
)

const (
	defaultStartAtSeconds = 3
)

// CronSchedulers schedules all the jobs so that data will be available for the exporter to serve.
func (conf *client) CronSchedulers() { //nolint:funlen
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("api", conf.defaultAPICron)) //nolint:errcheck
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronMessages("disk", conf.diskCron))      //nolint:errcheck

	scheduler := gocron.NewScheduler(time.UTC)

	if !funk.Contains(conf.skipMetrics, common.MetricPipelineSize) {
		conf.schedule(scheduler, common.MetricPipelineSize, conf.updateDiskUsage)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricSystemAdminsCount) {
		conf.schedule(scheduler, common.MetricSystemAdminsCount, conf.updateAdminsInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricConfigRepoCount) {
		conf.schedule(scheduler, common.MetricConfigRepoCount, conf.updateConfigRepoInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricConfiguredBackup) {
		conf.schedule(scheduler, common.MetricConfiguredBackup, conf.updateBackupInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricAgentDown) &&
		!funk.Contains(conf.skipMetrics, common.MetricAgentDiskSpace) &&
		!funk.Contains(conf.skipMetrics, common.MetricAgentsCount) {
		metricName := fmt.Sprintf("%s/%s/%s", common.MetricAgentDown, common.MetricAgentDiskSpace, common.MetricAgentsCount)
		conf.schedule(scheduler, metricName, conf.updateAgentsInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricServerHealth) {
		conf.schedule(scheduler, common.MetricServerHealth, conf.updateHealthInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricPipelineGroupCount) {
		conf.schedule(scheduler, common.MetricPipelineGroupCount, conf.updatePipelineGroupInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricEnvironmentCountAll) {
		conf.schedule(scheduler, common.MetricEnvironmentCountAll, conf.updateEnvironmentInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricVersion) {
		conf.schedule(scheduler, common.MetricVersion, conf.updateVersionInfo)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricJobStatus) {
		conf.schedule(scheduler, common.MetricJobStatus, conf.updateAgentJobRunHistory)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricPipelines) {
		conf.schedule(scheduler, common.MetricPipelines, conf.updatePipelines)
	}

	if !funk.Contains(conf.skipMetrics, common.MetricPipelineState) {
		conf.schedule(scheduler, common.MetricPipelineState, conf.updatePipelineState)
	}

	scheduler.NextRun()
	scheduler.StartAsync()
	scheduler.SingletonMode()
	scheduler.StartBlocking()
	scheduler.SingletonModeAll()
}

func (conf *client) schedule(scheduler *gocron.Scheduler, metric string, taskFunc func()) {
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(metric)) //nolint:errcheck
	job, err := scheduler.Every(conf.getCron(metric)).StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(taskFunc)
	if err != nil {
		log.Fatal(err)
	}
	job.SetEventListeners(func() {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(metric, job.RunCount())) //nolint:errcheck
	}, func() {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(metric, job.RunCount())) //nolint:errcheck
	})

	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		fmt.Printf("job %s\n just panicked", jobName)
		fmt.Printf("do something to handle the panic")
	})
}
