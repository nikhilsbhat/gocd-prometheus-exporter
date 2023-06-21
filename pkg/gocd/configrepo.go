package gocd

func (conf *client) updateConfigRepoInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	repos, err := goClient.GetConfigRepos()
	if err != nil {
		conf.logger.Error(apiError("config repo", err.Error()))
	}

	if err == nil {
		CurrentConfigRepos = repos
	}

	defer conf.lock.Unlock()
}
