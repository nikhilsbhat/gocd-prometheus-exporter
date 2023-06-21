package gocd

func (conf *client) updateVersionInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	version, err := goClient.GetVersionInfo()
	if err != nil {
		conf.logger.Error(apiError("version", err.Error()))
	}

	if err == nil {
		CurrentVersion = version
	}

	defer conf.lock.Unlock()
}
