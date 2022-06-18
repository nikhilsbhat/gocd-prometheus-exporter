package gocd

import (
	"testing"

	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetBackupInfo(t *testing.T) {
	t.Run("", func(t *testing.T) {
		logger := promlog.New(&promlog.Config{})
		client := NewConfig(
			"http://localhost:8153/go",
			"",
			"",
			"info",
			"",
			"",
			nil,
			nil,
			logger,
		)

		backup, err := client.GetBackupInfo()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(backup.Schedule))
	})
}
