package app_test

import (
	"path/filepath"
	"testing"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/app"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	config := app.Config{
		Port:         8090,
		GoCdUserName: "test",
		GoCdPassword: "password",
		APICron:      "@every 1m",
		Endpoint:     "/metrics",
	}
	t.Run("Should fetch config successfully, by merging both config", func(t *testing.T) {
		expected := app.Config{
			GoCdUserName:      "testing",
			GoCdPassword:      "password",
			InsecureTLS:       false,
			GoCdPipelinesPath: []string{"/path/to/pipeline/directory1", "/path/to/pipeline/directory2", "/path/to/pipeline/directory3"},
			Port:              9995,
			LogLevel:          "debug",
			SkipMetrics:       []string{"backup_configured", "admin_count"},
			APICron:           "@every @2m",
			MetricCron:        map[string]string{"agent_down": "@every 60s"},
			Endpoint:          "/new-metrics",
		}

		path, err := filepath.Abs("../../infrastructure/gocd-prometheus-exporter.fixture.yaml")
		assert.NoError(t, err)

		actual, err := app.GetConfig(config, path)
		assert.NoError(t, err)
		assert.Equal(t, &expected, actual)
	})

	t.Run("Should error out while fetching config due to wrong config file path", func(t *testing.T) {
		actual, err := app.GetConfig(config, "gocd-prometheus-exporter.fixture.yaml")
		assert.EqualError(t, err, "fetching config file information failed with: "+
			"stat gocd-prometheus-exporter.fixture.yaml: no such file or directory")
		assert.Equal(t, actual, &config)
	})

	t.Run("Should error out while decoding config from malformed config file", func(t *testing.T) {
		path, err := filepath.Abs("../../infrastructure/gocd-prometheus-exporter.fixture.malformed.yaml")
		assert.NoError(t, err)

		actual, err := app.GetConfig(config, path)
		assert.EqualError(t, err, "parsing config file errored with: yaml: line 13: could not find expected ':'")
		assert.Equal(t, actual, &config)
	})
}
