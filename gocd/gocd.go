package gocd

// GoCd implements methods to get various information regarding GoCd.
type GoCd interface {
	GetNodesInfo() (NodesConfig, error)
	GetDiskSize(path string) (float64, string, error)
	GetHealthInfo() ([]ServerHealth, error)
}
