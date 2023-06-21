package gocd

func (conf *client) updateAdminsInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	admins, err := goClient.GetSystemAdmins()
	if err != nil {
		conf.logger.Error(apiError("system admin", err.Error()))
	}

	if err == nil {
		CurrentSystemAdmins = admins
	}

	defer conf.lock.Unlock()
}
