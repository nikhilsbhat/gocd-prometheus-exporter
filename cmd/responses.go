package cmd

import "fmt"

func getCAErrMsg(path string) string {
	return fmt.Sprintf("an error occurred while reading CA file at %s", path)
}

func getAppTerminationMsg(port int) string {
	return fmt.Sprintf("terminate gocd-prometheus-exporter running on port: %d", port)
}
