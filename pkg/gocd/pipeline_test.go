package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"
	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func Test_client_GetPipelines(t *testing.T) {
	t.Run("", func(t *testing.T) {
		logger := promlog.New(&promlog.Config{}) //nolint:exhaustivestruct
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			"",
			"",
			nil,
			nil,
			nil,
			nil,
			logger,
		)

		pipelines, err := client.GetPipelines()
		assert.NoError(t, err)
		assert.Equal(t, 6, len(pipelines.Pipeline))
	})
}

func Test_client_getPipelineName(t *testing.T) {
	t.Run("should be able to fetch pipeline name from the href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "animation-and-action-movies", name)
	})
	t.Run("should fetch malformed pipeline name as malformed/invalid (prefix) href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8153/go/api/feed/pipelinesss/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "/go/api/feed/pipelinesss/animation-and-action-movies", name)
	})
	t.Run("should fetch malformed pipeline name as malformed/invalid (suffix) href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/test/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "animation-and-action-movies/test", name)
	})
	t.Run("should fail while fetching pipeline name from malformed/invalid href", func(t *testing.T) {
		pipelineLink := "://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.EqualError(t, err, "parse \"://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml\": missing protocol scheme")
		assert.Equal(t, "", name)
	})
}

func Test_client_GetPipelineStatus(t *testing.T) {
	t.Run("should be able to get pipeline statuses of all pipeline present in gocd", func(t *testing.T) {
		logger := promlog.New(&promlog.Config{}) //nolint:exhaustivestruct
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			"",
			"",
			nil,
			nil,
			nil,
			nil,
			logger,
		)

		gocd.CurrentPipelines = []string{"action-movies-manual", "gocd-prometheus-exporter"}
		piplineStates, err := client.GetPipelineState()
		assert.NoError(t, err)
		assert.Equal(t, "action-movies-manual", piplineStates[0].Name)
	})
}
