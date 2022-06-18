package gocd

import (
	"testing"

	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should be able retrieve config repo information", func(t *testing.T) {
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

		repos, err := client.GetConfigRepoInfo()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(repos))
	})
}
