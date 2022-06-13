package common

const (
	Namespace                = "gocd"
	TypeLink                 = "link"
	TypeDir                  = "dir"
	LogCategoryMsg           = "msg"
	LogCategoryErr           = "err"
	GoCdHeaderVersionSeven   = "application/vnd.go.cd.v7+json"
	GoCdHeaderVersionOne     = "application/vnd.go.cd.v1+json"
	GoCdAgentsEndpoint       = "/api/agents"
	GoCdVersionEndpoint      = "/api/version"
	GoCdServerHealthEndpoint = "/api/server_health_messages"
	GoCdDisconnectedState    = "LostContact"
	GoCdReturnErrorMessage   = `gocd server returned code %d with message`
)

func Float(value interface{}) float64 {
	switch value.(type) {
	case int64:
		return value.(float64)
	case string:
		return float64(0)
	default:
		return value.(float64)
	}
}
