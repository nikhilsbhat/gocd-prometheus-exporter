package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/app"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/exporter"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/urfave/cli/v2"
)

const (
	flagPipelinePath     = "pipeline-path"
	flagPipelinePathRoot = "pipeline-root-path"
	flagLogLevel         = "log-level"
	flagExporterPort     = "port"
	flagExporterEndpoint = "endpoint"
	flagGoCdBaseURL      = "goCd-server-url"
	flagGoCdUsername     = "goCd-username"
	flagGoCdPassword     = "goCd-password"
	flagInsecureTLS      = "insecure-tls"
	flagCaPath           = "ca-path"
	flagGraceDuration    = "grace-duration"
	flagConfigPath       = "config-file"
	flagSkipMetrics      = "skip-metrics"
	flagDiskCronSchedule = "disk-cron-schedule"
	flagAPICronSchedule  = "api-cron-schedule"
)

const (
	defaultAppPort = 8090
)

var (
	redirectData = `<html>
			 <head><title>GoCd Exporter</title></head>
			 <body>
			 <h1>GoCd Exporter</h1>
			 <p><a href='%s'>Metrics</a></p>
			 </body>
			 </html>`
	sigChan = make(chan os.Signal)
)

const (
	defaultGraceDuration = 5
)

// App returns the cli for gocd-prometheus-exporter.
func App() *cli.App {
	return &cli.App{
		Name:                 "gocd-prometheus-exporter",
		Usage:                "Utility to collect metrics of GoCd for prometheus",
		UsageText:            "gocd-prometheus-exporter [flags]",
		EnableBashCompletion: true,
		HideHelp:             false,
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "version of the gocd-prometheus-exporter",
				Action:  version.AppVersion,
			},
		},
		Flags:  registerFlags(),
		Action: goCdExport,
	}
}

func registerFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    flagPipelinePath,
			Usage:   "list of paths to pipelines that needs tp be monitored",
			Aliases: []string{"pt"},
		},
		&cli.StringFlag{
			Name:    flagPipelinePathRoot,
			Usage:   "root path of pipeline that needs to be monitored",
			Aliases: []string{"root-path"},
		},
		&cli.StringFlag{
			Name:    flagLogLevel,
			Usage:   "set log level for the GoCd exporter",
			Aliases: []string{"log"},
			Value:   "info",
		},
		&cli.IntFlag{
			Name:    flagExporterPort,
			Usage:   "port on which the metrics to be exposed",
			Value:   defaultAppPort,
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:    flagExporterEndpoint,
			Usage:   "path under which the metrics to be exposed",
			Value:   "/metrics",
			Aliases: []string{"e"},
		},
		&cli.StringFlag{
			Name:    flagGoCdBaseURL,
			Usage:   "GoCd server url to which the exporter needs to be connected",
			Aliases: []string{"server"},
		},
		&cli.StringFlag{
			Name:    flagGoCdUsername,
			Usage:   "required username for establishing connection to GoCd server if auth enabled",
			Aliases: []string{"user"},
			EnvVars: []string{"GOCD_USERNAME"},
		},
		&cli.StringFlag{
			Name:    flagGoCdPassword,
			Usage:   "required password for establishing connection to GoCd server if auth enabled",
			Aliases: []string{"password"},
			EnvVars: []string{"GOCD_PASSWORD"},
		},
		&cli.BoolFlag{
			Name:    flagInsecureTLS,
			Usage:   "enable insecure TLS if you wish to connect to GOCD insecurily",
			Value:   false,
			Aliases: []string{"insecure"},
		},
		&cli.StringFlag{
			Name:    flagCaPath,
			Usage:   "path to file containing CA information to make secure connections to GoCd",
			Aliases: []string{"ca"},
		},
		&cli.DurationFlag{
			Name:    flagGraceDuration,
			Usage:   "time duration to wait before stopping the service",
			Aliases: []string{"d"},
			Value:   time.Second * defaultGraceDuration,
		},
		&cli.StringFlag{
			Name:    flagConfigPath,
			Usage:   "path to file containing configurations for exporter",
			Aliases: []string{"c"},
			Value:   filepath.Join(os.Getenv("HOME"), fmt.Sprintf("%s.%s", common.ExporterConfigFileName, common.ExporterConfigFileExt)),
		},
		&cli.StringSliceFlag{
			Name:    flagSkipMetrics,
			Usage:   "list of metrics to be skipped",
			Aliases: []string{"sk"},
		},
		&cli.StringFlag{
			Name: flagAPICronSchedule,
			Usage: `cron expression to schedule the metric collection.
                    		- 'gocd-prometheus-exporter' schedules the job to collect the metrics in the specified intervals
                      			and stores the latest values in memory.
                      		- This is to reduce the load on the GoCd server when api requests are made to GoCd.
                      		- All expressions supported by https://github.com/robfig/cron will be supported`,
			Aliases: []string{"cron"},
			Value:   "@every 30s",
		},
		&cli.StringFlag{
			Name: flagDiskCronSchedule,
			Usage: `cron expression to schedule the pipeline disk usage metric collection.
                      		- This is to reduce the reduce resource consumed while computing the pipeline disk size.`,
			Aliases: []string{"disk-cron"},
			Value:   "@every 30s",
		},
	}
}

func goCdExport(context *cli.Context) error {
	config := app.Config{
		GoCdBaseURL:           context.String(flagGoCdBaseURL),
		GoCdUserName:          context.String(flagGoCdUsername),
		GoCdPassword:          context.String(flagGoCdPassword),
		InsecureTLS:           context.Bool(flagInsecureTLS),
		GoCdPipelinesPath:     context.StringSlice(flagPipelinePath),
		GoCdPipelinesRootPath: context.String(flagPipelinePathRoot),
		CaPath:                context.String(flagCaPath),
		Port:                  context.Int(flagExporterPort),
		Endpoint:              context.String(flagExporterEndpoint),
		LogLevel:              context.String(flagLogLevel),
		SkipMetrics:           context.StringSlice(flagSkipMetrics),
		APICron:               context.String(flagAPICronSchedule),
		DiskCron:              context.String(flagDiskCronSchedule),
		AppGraceDuration:      context.Duration(flagGraceDuration),
	}

	finalConfig, err := app.GetConfig(config, context.String(flagConfigPath))
	if err != nil {
		log.Println(err)
	}

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT) //nolint:govet

	promLogConfig := &promlog.Config{Level: &promlog.AllowedLevel{}, Format: &promlog.AllowedFormat{}}
	if err := promLogConfig.Level.Set(finalConfig.LogLevel); err != nil {
		return fmt.Errorf("configuring logger errored with: %w", err)
	}
	logger := promlog.New(promLogConfig)

	var caContent []byte
	if len(finalConfig.CaPath) != 0 {
		ca, err := ioutil.ReadFile(finalConfig.CaPath)
		if err != nil {
			level.Error(logger).Log(common.LogCategoryErr, getCAErrMsg(finalConfig.CaPath)) //nolint:errcheck
		}
		caContent = ca
	}

	pipelinePaths := make([]string, 0)
	pipelinePaths = append(pipelinePaths, finalConfig.GoCdPipelinesRootPath)
	pipelinePaths = append(pipelinePaths, finalConfig.GoCdPipelinesPath...)

	client := gocd.NewClient(
		finalConfig.GoCdBaseURL,
		finalConfig.GoCdUserName,
		finalConfig.GoCdPassword,
		finalConfig.LogLevel,
		finalConfig.APICron,
		finalConfig.DiskCron,
		finalConfig.MetricCron,
		caContent,
		pipelinePaths,
		finalConfig.SkipMetrics,
		logger,
	)

	// running schedules
	client.CronSchedulers()

	goCdExporter := exporter.NewExporter(logger, finalConfig.SkipMetrics)
	prometheus.MustRegister(goCdExporter)

	// listens to terminate signal
	go func() {
		sig := <-sigChan
		level.Info(logger).Log("msg", fmt.Sprintf("caught signal %v: terminating in %v", sig, finalConfig.AppGraceDuration)) //nolint:errcheck
		time.Sleep(context.Duration(flagGraceDuration))
		level.Info(logger).Log("msg", getAppTerminationMsg(finalConfig.Port)) //nolint:errcheck
		os.Exit(0)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(getRedirectData(finalConfig.Endpoint)))
	})

	level.Info(logger).Log(common.LogCategoryMsg, fmt.Sprintf("metrics will be exposed on port: %d", finalConfig.Port))         //nolint:errcheck
	level.Info(logger).Log(common.LogCategoryMsg, fmt.Sprintf("metrics will be exposed on endpoint: %s", finalConfig.Endpoint)) //nolint:errcheck
	http.Handle(finalConfig.Endpoint, promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", finalConfig.Port), nil); err != nil {
		return fmt.Errorf("starting server on specified port failed with: %w", err)
	}

	return nil
}

func getRedirectData(endpoint string) string {
	return fmt.Sprintf(redirectData, endpoint)
}
