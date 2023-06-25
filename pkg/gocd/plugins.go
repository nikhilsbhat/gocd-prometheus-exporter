package gocd

func (conf *client) updatePluginsInfo() {
	conf.lock.Lock()
	goClient := conf.getCronClient()

	plugins, err := goClient.GetPluginsInfo()
	if err != nil {
		conf.logger.Error(apiError("plugin information", err.Error()))
	}

	if err == nil {
		CurrentPluginInfo = plugins.Plugins
	}

	defer conf.lock.Unlock()
}
