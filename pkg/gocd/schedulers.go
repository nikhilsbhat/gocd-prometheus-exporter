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
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricPipelineSize)) //nolint:errcheck
		diskUsageJob, err := scheduler.Every(conf.diskCron).StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateDiskUsage)
		if err != nil {
			log.Fatal(err)
		}
		diskUsageJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage("pipeline size", diskUsageJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage("pipeline size", diskUsageJob.RunCount())) //nolint:errcheck
		})
		diskUsageJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricSystemAdminsCount) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricSystemAdminsCount)) //nolint:errcheck
		adminsInfoJob, err := scheduler.Every(conf.getCron(common.MetricSystemAdminsCount)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateAdminsInfo)
		if err != nil {
			log.Fatal(err)
		}
		adminsInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricSystemAdminsCount, adminsInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricSystemAdminsCount, adminsInfoJob.RunCount())) //nolint:errcheck
		})
		adminsInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricConfigRepoCount) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricConfigRepoCount)) //nolint:errcheck
		configRepoInfoJob, err := scheduler.Every(conf.getCron(common.MetricConfigRepoCount)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateConfigRepoInfo)
		if err != nil {
			log.Fatal(err)
		}
		configRepoInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricConfigRepoCount, configRepoInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricConfigRepoCount, configRepoInfoJob.RunCount())) //nolint:errcheck
		})
		configRepoInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricConfiguredBackup) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricConfiguredBackup)) //nolint:errcheck
		backupInfoJob, err := scheduler.Every(conf.getCron(common.MetricConfiguredBackup)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateBackupInfo)
		if err != nil {
			log.Fatal(err)
		}
		backupInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricConfiguredBackup, backupInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricConfiguredBackup, backupInfoJob.RunCount())) //nolint:errcheck
		})
		backupInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricAgentDown) &&
		!funk.Contains(conf.skipMetrics, common.MetricAgentDiskSpace) &&
		!funk.Contains(conf.skipMetrics, common.MetricAgentsCount) {
		metricName := fmt.Sprintf("%s/%s/%s", common.MetricAgentDown, common.MetricAgentDiskSpace, common.MetricAgentsCount)
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(metricName)) //nolint:errcheck
		agentsInfoJob, err := scheduler.Every(conf.getCron(common.MetricAgentDown)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateAgentsInfo)
		if err != nil {
			log.Fatal(err)
		}
		agentsInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(metricName, agentsInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(metricName, agentsInfoJob.RunCount())) //nolint:errcheck
		})
		agentsInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricServerHealth) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricServerHealth)) //nolint:errcheck
		healthInfoJob, err := scheduler.Every(conf.getCron(common.MetricServerHealth)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateHealthInfo)
		if err != nil {
			log.Fatal(err)
		}
		healthInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricServerHealth, healthInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricServerHealth, healthInfoJob.RunCount())) //nolint:errcheck
		})
		healthInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricPipelineGroupCount) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricPipelineGroupCount)) //nolint:errcheck
		pipelineGroupInfoJob, err := scheduler.Every(conf.getCron(common.MetricPipelineGroupCount)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updatePipelineGroupInfo)
		if err != nil {
			log.Fatal(err)
		}
		pipelineGroupInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricPipelineGroupCount, pipelineGroupInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricPipelineGroupCount, pipelineGroupInfoJob.RunCount())) //nolint:errcheck
		})
		pipelineGroupInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricEnvironmentCountAll) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricEnvironmentCountAll)) //nolint:errcheck
		groupInfoJob, err := scheduler.Every(conf.getCron(common.MetricEnvironmentCountAll)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateEnvironmentInfo)
		if err != nil {
			log.Fatal(err)
		}
		groupInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricEnvironmentCountAll, groupInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricEnvironmentCountAll, groupInfoJob.RunCount())) //nolint:errcheck
		})
		groupInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricVersion) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricVersion)) //nolint:errcheck
		versionInfoJob, err := scheduler.Every(conf.getCron(common.MetricVersion)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateVersionInfo)
		if err != nil {
			log.Fatal(err)
		}
		versionInfoJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricVersion, versionInfoJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricVersion, versionInfoJob.RunCount())) //nolint:errcheck
		})
		versionInfoJob.SingletonMode()
	}

	if !funk.Contains(conf.skipMetrics, common.MetricJobStatus) {
		level.Info(conf.logger).Log(common.LogCategoryMsg, getCronEnbaled(common.MetricJobStatus)) //nolint:errcheck
		agentJobRunHistoryJob, err := scheduler.Every(conf.getCron(common.MetricJobStatus)).
			StartAt(time.Now().Add(defaultStartAtSeconds * time.Second)).Do(conf.updateAgentJobRunHistory)
		if err != nil {
			log.Fatal(err)
		}
		agentJobRunHistoryJob.SetEventListeners(func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage(common.MetricJobStatus, agentJobRunHistoryJob.RunCount())) //nolint:errcheck
		}, func() {
			level.Info(conf.logger).Log(common.LogCategoryMsg, getCronCompleteMessage(common.MetricJobStatus, agentJobRunHistoryJob.RunCount())) //nolint:errcheck
		})
		agentJobRunHistoryJob.SingletonMode()
	}

	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		fmt.Printf("job %s\n just panicked", jobName)
		fmt.Printf("do something to handle the panic")
	})

	scheduler.StartAsync()
	scheduler.StartBlocking()
	scheduler.SingletonModeAll()
}
