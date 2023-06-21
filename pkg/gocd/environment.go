package gocd

func (conf *client) updateEnvironmentInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	environmentInfo, err := goClient.GetEnvironments()
	if err != nil {
		conf.logger.Error(apiError("environment", err.Error()))
	}

	if err == nil {
		CurrentEnvironments = environmentInfo
	}

	defer conf.lock.Unlock()
}
