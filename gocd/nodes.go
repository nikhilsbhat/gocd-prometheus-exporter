package gocd

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetNodesInfo implements method that fetches the details of all the agents present in GoCd server
func (conf *Config) GetNodesInfo() ([]Node, error) {
	conf.lock.Lock()
	conf.client.SetHeaders(map[string]string{
		"Accept": common.GoCdHeaderVersionSeven,
	})
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "trying to retrieve nodes information present in GoCd") //nolint:errcheck

	var nodesConf NodesConfig
	resp, err := conf.client.R().SetResult(&nodesConf).Get(common.GoCdAgentsEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf(common.GoCdReturnErrorMessage, resp.StatusCode()))
	}

	conf.lock.Unlock()
	level.Debug(conf.logger).Log(common.LogCategoryMsg, "successfully retrieved nodes information from GoCd") //nolint:errcheck
	return nodesConf.Config.Config, nil
}

func (conf *Config) configureGetNodesInfo() {
	scheduleGetNodesInfo := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleGetNodesInfo.AddFunc(conf.otherCron, func() {
		nodesInfo, err := conf.GetNodesInfo()
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, fmt.Sprintf("retrieving agents information errored with: %s", err.Error())) //nolint:errcheck
		}
		CurrentNodeConfig = nodesInfo
	})

	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleGetNodesInfo.Start()
}
