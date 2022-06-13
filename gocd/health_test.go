package gocd

import (
	"log"
	"testing"

	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetHealthInfo(t *testing.T) {
	t.Run("should be able to fetch the server health status", func(t *testing.T) {
		logger := promlog.New(&promlog.Config{})
		client := NewConfig(
			"http://localhost:8153/go",
			"",
			"",
			nil,
			logger,
		)

		healthStatus, err := client.GetHealthInfo()
		assert.NoError(t, err)
		log.Println(healthStatus)
		assert.Equal(t, 5, len(healthStatus))
	})
}
