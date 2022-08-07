package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/gocd"
	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
)

func Test_client_GetAgentJobRunHistory(t *testing.T) {
	t.Run("should be able to fetch the agent job run history", func(t *testing.T) {
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
			nil,
			logger,
		)

		gocd.CurrentAgentsConfig = []gocd.Agent{{ID: "6132c45f-9818-42c9-9bd1-154132bd265f"}}
		history, err := client.GetAgentJobRunHistory()
		assert.NoError(t, err)
		assert.Equal(t, "Passed", history[0].Jobs[0].Result)
		assert.Equal(t, "", history)
	})
}
