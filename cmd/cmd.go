package cmd

import (
	"fmt"
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
	"github.com/nikhilsbhat/gocd-prometheus-exporter/version"
	goCDLogger "github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/sirupsen/logrus"
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
	flagGoCDBearerToken  = "goCd-bearer-token" //nolint:gosec
	flagInsecureTLS      = "insecure-tls"
	flagCaPath           = "ca-path"
	flagGraceDuration    = "grace-duration"
	flagConfigPath       = "config-file"
	flagSkipMetrics      = "skip-metrics"
	flagAPICronSchedule  = "api-cron-schedule"
)

const (
	defaultAppPort = 8090
	defaultTimeout = 30
)

var sigChan = make(chan os.Signal)

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
		Authors: []*cli.Author{
			{
				Name:  "Nikhil Bhat",
				Email: "nikhilsbhat93@gmail.com",
			},
		},
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
		&cli.StringFlag{
			Name:    flagGoCDBearerToken,
			Usage:   "required bearer-token for establishing connection to GoCd server if auth enabled",
			Aliases: []string{"token"},
			EnvVars: []string{"GOCD_BEARER_TOKEN"},
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
			Value:   filepath.Join(os.Getenv("HOME"), fmt.Sprintf("%s.%s", common.GoCDExporterName, common.ExporterConfigFileExt)),
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
                      		- All expressions supported by github.com/go-co-op/gocron will be supported`,
			Aliases: []string{"cron"},
			Value:   "30s",
		},
	}
}

func goCdExport(context *cli.Context) error {
	logger := logrus.New()
	logger.SetLevel(goCDLogger.GetLoglevel(context.String(flagLogLevel)))
	logger.WithField(common.GoCDExporterName, true)
	logger.SetFormatter(&logrus.JSONFormatter{})

	config := app.Config{
		GoCdBaseURL:           context.String(flagGoCdBaseURL),
		GoCdUserName:          context.String(flagGoCdUsername),
		GoCdPassword:          context.String(flagGoCdPassword),
		GoCDBearerToken:       context.String(flagGoCDBearerToken),
		InsecureTLS:           context.Bool(flagInsecureTLS),
		GoCdPipelinesPath:     context.StringSlice(flagPipelinePath),
		GoCdPipelinesRootPath: context.String(flagPipelinePathRoot),
		CaPath:                context.String(flagCaPath),
		Port:                  context.Int(flagExporterPort),
		Endpoint:              context.String(flagExporterEndpoint),
		LogLevel:              context.String(flagLogLevel),
		SkipMetrics:           context.StringSlice(flagSkipMetrics),
		APICron:               context.String(flagAPICronSchedule),
		AppGraceDuration:      context.Duration(flagGraceDuration),
	}

	finalConfig, err := app.GetConfig(config, context.String(flagConfigPath))
	if err != nil {
		log.Println(err)
	}

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT) //nolint:govet

	var caContent []byte
	if len(finalConfig.CaPath) != 0 {
		ca, err := os.ReadFile(finalConfig.CaPath)
		if err != nil {
			logger.Fatal(getCAErrMsg(finalConfig.CaPath))
		}
		caContent = ca
	}

	client := gocd.NewClient(*finalConfig, logger, caContent)

	// running schedules
	go func() {
		client.CronSchedulers()
	}()

	goCdExporter := exporter.NewExporter(logger, finalConfig.SkipMetrics)
	prometheus.MustRegister(goCdExporter)

	// listens to terminate signal
	go func() {
		sig := <-sigChan
		logger.Infof("caught signal %v: terminating in %v", sig, finalConfig.AppGraceDuration)
		time.Sleep(context.Duration(flagGraceDuration))
		logger.Info(getAppTerminationMsg(finalConfig.Port))
		os.Exit(0)
	}()

	logger.Infof("metrics will be exposed on port: %d", finalConfig.Port)
	logger.Infof("metrics will be exposed on endpoint: %s", finalConfig.Endpoint)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
	})

	landingConfig := web.LandingConfig{
		Name:        "GoCD Exporter",
		Description: "Prometheus Exporter for CI/CD server GoCD",
		Version:     version.GetAppVersion(),
		Links: []web.LandingLinks{
			{
				Address: finalConfig.Endpoint,
				Text:    "Metrics",
			},
			{
				Address: "/health",
				Text:    "Health",
			},
		},
	}

	landingPage, err := web.NewLandingPage(landingConfig)
	if err != nil {
		logger.Fatal(err)
	}

	http.Handle("/", landingPage)
	http.Handle(finalConfig.Endpoint, promhttp.Handler())

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", finalConfig.Port),
		ReadHeaderTimeout: defaultTimeout * time.Second,
	}

	if err = server.ListenAndServe(); err != nil {
		return fmt.Errorf("starting server on specified port failed with: %w", err)
	}

	return nil
}
