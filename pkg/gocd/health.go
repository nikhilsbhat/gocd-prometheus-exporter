package gocd

func (conf *client) updateHealthInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	healthInfo, err := goClient.GetServerHealthMessages()
	if err != nil {
		conf.logger.Error(apiError("server health", err.Error()))
	}

	if err == nil {
		CurrentServerHealth = healthInfo
	}

	defer conf.lock.Unlock()
}
