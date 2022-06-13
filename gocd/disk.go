package gocd

import (
	"fmt"
	"os"
	"path/filepath"

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
		filepath.Walk(path, readSize)
		close(sizes)
	}()

	for s := range sizes {
		dirSize += s
	}

	sizeMB := float64(dirSize)

	return sizeMB
}
