package gocd

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/thoas/go-funk"
)

const (
	defaultStartAtSeconds = 3
)

// CronSchedulers schedules all the jobs so that data will be available for the exporter to serve.
func (conf *client) CronSchedulers() { //nolint:funlen
	conf.logger.Infof(getCronMessages("api", conf.config.APICron))

	scheduler := gocron.NewScheduler(time.UTC)

	if !funk.Contains(conf.config.SkipMetrics, common.MetricSystemAdminsCount) {
		conf.schedule(scheduler, common.MetricSystemAdminsCount, conf.updateAdminsInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricConfigRepoCount) {
		conf.schedule(scheduler, common.MetricConfigRepoCount, conf.updateConfigRepoInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricConfiguredBackup) {
		conf.schedule(scheduler, common.MetricConfiguredBackup, conf.updateBackupInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricAgentDown) &&
		!funk.Contains(conf.config.SkipMetrics, common.MetricAgentDiskSpace) &&
		!funk.Contains(conf.config.SkipMetrics, common.MetricAgentsCount) {
		metricName := fmt.Sprintf("%s/%s/%s", common.MetricAgentDown, common.MetricAgentDiskSpace, common.MetricAgentsCount)
		conf.schedule(scheduler, metricName, conf.updateAgentsInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricServerHealth) {
		conf.schedule(scheduler, common.MetricServerHealth, conf.updateHealthInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricPipelineGroupCount) {
		conf.schedule(scheduler, common.MetricPipelineGroupCount, conf.updatePipelineGroupInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricEnvironmentCountAll) {
		conf.schedule(scheduler, common.MetricEnvironmentCountAll, conf.updateEnvironmentInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricVersion) {
		conf.schedule(scheduler, common.MetricVersion, conf.updateVersionInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricJobStatus) {
		conf.schedule(scheduler, common.MetricJobStatus, conf.updateAgentJobRunHistory)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricPipelines) {
		conf.schedule(scheduler, common.MetricPipelines, conf.updatePipelines)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricPipelineState) {
		conf.schedule(scheduler, common.MetricPipelineState, conf.updatePipelineState)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricElasticAgentProfileUsage) {
		conf.schedule(scheduler, common.MetricElasticAgentProfileUsage, conf.updateElasticProfileInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricPlugins) {
		conf.schedule(scheduler, common.MetricPlugins, conf.updatePluginsInfo)
	}

	if !funk.Contains(conf.config.SkipMetrics, common.MetricPipelineNotRun) {
		conf.schedule(scheduler, common.MetricPipelineNotRun, conf.updatePipelineRunInLastXDays)
	}

	scheduler.NextRun()
	scheduler.StartAsync()
	scheduler.SingletonMode()
	scheduler.StartBlocking()
	scheduler.SingletonModeAll()
}

func (conf *client) schedule(scheduler *gocron.Scheduler, metric string, taskFunc func()) {
	conf.logger.Infof(getCronEnbaled(metric))
	job, err := scheduler.Every(conf.getCron(metric)).StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(taskFunc)
	if err != nil {
		log.Fatal(err)
	}
	job.SetEventListeners(func() {
		conf.logger.Infof(getCronScheduledMessage(metric, job.RunCount()))
	}, func() {
		conf.logger.Infof(getCronCompleteMessage(metric, job.RunCount()))
	})

	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		fmt.Printf("job %s\n just panicked", jobName)
		fmt.Printf("do something to handle the panic")
	})
}
