package gocd

// GoCd implements methods to get various information regarding GoCd.
type GoCd interface {
	GetNodesInfo() ([]Node, error)
	GetDiskSize(path string) (float64, string, error)
	GetHealthInfo() ([]ServerHealth, error)
	GetConfigRepoInfo() ([]ConfigRepo, error)
	GetAdminsInfo() (SystemAdmins, error)
	GetPipelineGroupInfo() ([]PipelineGroup, error)
	ScheDulers()
}
