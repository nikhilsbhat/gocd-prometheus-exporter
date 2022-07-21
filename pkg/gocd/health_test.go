package gocd_test

import (
	"log"
	"testing"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"

	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetHealthInfo(t *testing.T) {
	t.Run("should be able to fetch the server health status", func(t *testing.T) {
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

		healthStatus, err := client.GetHealthInfo()
		assert.NoError(t, err)
		log.Println(healthStatus)
		assert.Equal(t, 5, len(healthStatus))
	})
}
