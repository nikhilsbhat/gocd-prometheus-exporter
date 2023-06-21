package gocd

func (conf *client) updateBackupInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	backupInfo, err := goClient.GetBackupConfig()
	if err != nil {
		conf.logger.Error(apiError("gocd backup", err.Error()))
	}

	if err == nil {
		CurrentBackupConfig = backupInfo
	}

	defer conf.lock.Unlock()
}
