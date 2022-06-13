package cmd

import (
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/exporter"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/gocd"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
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
)

var (
	redirectData = `<html>
			 <head><title>GoCd Exporter</title></head>
			 <body>
			 <h1>GoCd Exporter</h1>
			 <p><a href='%s'>Metrics</a></p>
			 </body>
			 </html>`
)

// App returns the cli for gocd-prometheus-exporter
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
		&cli.StringFlag{
			Name:    flagExporterPort,
			Usage:   "port on which the metrics to be exposed",
			Value:   "8090",
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
	}
}

func getRedirectData(endpoint string) string {
	return fmt.Sprintf(redirectData, endpoint)
}

func goCdExport(context *cli.Context) error {
	promLogConfig := &promlog.Config{Level: &promlog.AllowedLevel{}, Format: &promlog.AllowedFormat{}}
	if err := promLogConfig.Level.Set(context.String(flagLogLevel)); err != nil {
		return err
	}
	logger := promlog.New(promLogConfig)

	caPath := context.String(flagCaPath)
	var caContent []byte
	if len(caPath) != 0 {
		ca, err := ioutil.ReadFile(caPath)
		if err != nil {
			level.Error(logger).Log(common.LogCategoryErr, fmt.Sprintf("an error occured while reading CA file: %s", caPath))
		}
		caContent = ca
	}

	client := gocd.NewConfig(
		context.String(flagGoCdBaseURL),
		context.String(flagGoCdUsername),
		context.String(flagGoCdPassword),
		caContent,
		logger,
	)

	pipelinePaths := context.StringSlice(flagPipelinePath)
	if len(context.String(flagPipelinePathRoot)) != 0 {
		pipelinePaths = append(pipelinePaths, context.String(flagPipelinePathRoot))
	}
	goCdExporter := exporter.NewExporter(logger, client, pipelinePaths)
	prometheus.MustRegister(goCdExporter)

	port := context.String(flagExporterPort)
	endpoint := context.String(flagExporterEndpoint)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(getRedirectData(endpoint)))
	})

	level.Debug(logger).Log(common.LogCategoryMsg, fmt.Sprintf("metrics will be exposed on port: %s", port))
	level.Debug(logger).Log(common.LogCategoryMsg, fmt.Sprintf("metrics will be exposed on endpoint: %s", endpoint))
	http.Handle(endpoint, promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		return err
	}
	return nil
}
