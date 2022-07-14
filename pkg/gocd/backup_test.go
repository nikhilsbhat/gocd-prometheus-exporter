package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"

	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetBackupInfo(t *testing.T) {
	t.Run("", func(t *testing.T) {
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

		backup, err := client.GetBackupInfo()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(backup.Schedule))
	})
}
