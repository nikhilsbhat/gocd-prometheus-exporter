package gocd

import (
	"fmt"
)

var (
	cronMessage          = `%s cron will be scheduled for %s as specified`
	cronScheduledMessage = `%d instance of %s cron got scheduled`
	cronCompleteMessage  = `scheduled %d instance of %s cron was completed`
	cronEnabled          = `cron is enabled for %s metric collection`
)

func apiError(component, errMsg string) error {
	return fmt.Errorf("retrieving %s information errored with: %s", component, errMsg) //nolint:goerr113
}

func getCronMessages(component, schedule string) string {
	return fmt.Sprintf(cronMessage, component, schedule)
}

func getCronEnbaled(component string) string {
	return fmt.Sprintf(cronEnabled, component)
}

func getCronScheduledMessage(component string, instance int) string {
	return fmt.Sprintf(cronScheduledMessage, instance, component)
}

func getCronCompleteMessage(component string, instance int) string {
	return fmt.Sprintf(cronCompleteMessage, instance, component)
}
