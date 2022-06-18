package gocd

import "fmt"

var (
	goCdAPIError       = `retrieving %s information errored with: %s`
	goCdTryMessage     = `trying to retrieve %s information present in GoCd`
	goCdSuccessMessage = `successfully retrieved information of %s information from GoCd`
	cronMessage        = `%s cron will be scheduled for %s as specified`
	symlinkMessage     = `path %s is link to %s so fetching size of destination`
)

func getAPIErrMsg(component, errMsg string) string {
	return fmt.Sprintf(goCdAPIError, component, errMsg)
}

func getTryMessages(component string) string {
	return fmt.Sprintf(goCdTryMessage, component)
}

func getSuccessMessages(component string) string {
	return fmt.Sprintf(goCdSuccessMessage, component)
}

func getCronMessages(component, schedule string) string {
	return fmt.Sprintf(cronMessage, component, schedule)
}

func getLinkMessage(link, path string) string {
	return fmt.Sprintf(symlinkMessage, link, path)
}
