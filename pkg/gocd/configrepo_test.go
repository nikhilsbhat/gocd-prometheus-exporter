package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"

	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should be able retrieve config repo information", func(t *testing.T) {
		logger := promlog.New(&promlog.Config{}) //nolint:exhaustivestruct
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"",
			"",
			"info",
			"",
			"",
			nil,
			nil,
			nil,
			logger,
		)

		repos, err := client.GetConfigRepoInfo()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(repos))
	})
}
