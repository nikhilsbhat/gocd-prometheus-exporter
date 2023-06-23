package gocd

import (
	"github.com/nikhilsbhat/gocd-sdk-go"
)

type ElasticProfileUsage struct {
	Name  string
	Usage []gocd.ElasticProfileUsage
}

func (conf *client) updateElasticProfileInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	profiles, err := goClient.GetElasticAgentProfiles()
	if err != nil {
		conf.logger.Error(apiError("elastic agent profiles", err.Error()))
	}

	elasticAgentProfilesUsage := make([]ElasticProfileUsage, 0)

	for _, profileID := range profiles.CommonConfigs {
		var elasticAgentProfileUsage []gocd.ElasticProfileUsage

		response, err := goClient.GetElasticAgentProfileUsage(profileID.ID)
		if err != nil {
			conf.logger.Error(apiError("elastic agent profiles", err.Error()))
		}

		elasticAgentProfileUsage = append(elasticAgentProfileUsage, response...)

		elasticAgentProfilesUsage = append(elasticAgentProfilesUsage, ElasticProfileUsage{
			Name:  profileID.ID,
			Usage: elasticAgentProfileUsage,
		})
	}

	if err == nil {
		CurrentElasticProfileUsage = elasticAgentProfilesUsage
	}

	defer conf.lock.Unlock()
}
