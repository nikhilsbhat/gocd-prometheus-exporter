package gocd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kit/log/level"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/common"
)

// GetDiskSize retrieves size of the specified path along with type, it would be link if path is symlink
func (conf *Config) GetDiskSize(path string) (float64, string, error) {
	var pathType string
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, "", fmt.Errorf(fmt.Sprintf("stating path %s failed with error %s", path, err.Error()))
	}

	pathType = common.TypeDir
	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		pathType = common.TypeLink
		originPath, err := os.Readlink(path)
		if err != nil {
			return 0, "", err
		}
		level.Debug(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("path %s is link to %s so fetching size of destination", path, originPath)) //nolint:errcheck
		path = originPath
	}

	return diskSize(path), pathType, nil
}
func diskSize(path string) float64 {
	var dirSize int64 = 0

	sizes := make(chan int64)
	readSize := func(path string, file os.FileInfo, err error) error {
		if err != nil || file == nil {
			return err
		}

		if !file.IsDir() {
			sizes <- file.Size()
		}
		return nil
	}

	go func() {
		filepath.Walk(path, readSize) //nolint:errcheck
		close(sizes)
	}()

	for s := range sizes {
		dirSize += s
	}

	sizeMB := float64(dirSize)

	return sizeMB
}
