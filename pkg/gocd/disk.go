package gocd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/pkg/common"

	"github.com/go-kit/log/level"
	"github.com/robfig/cron/v3"
)

// GetDiskSize retrieves size of the specified path along with type, it would be link if path is symlink.
func (conf *client) GetDiskSize(path string) (float64, string, error) {
	var pathType string
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, "", fmt.Errorf("stating path %s failed with error %w", path, err) //nolint:goerr113
	}

	pathType = common.TypeDir
	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		pathType = common.TypeLink
		originPath, err := os.Readlink(path)
		if err != nil {
			return 0, "", fmt.Errorf("reading link errored with: %w", err)
		}
		level.Debug(conf.logger).Log(common.LogCategoryMsg, getLinkMessage(path, originPath)) //nolint:errcheck
		path = originPath
	}

	return diskSize(path), pathType, nil
}

func diskSize(path string) float64 {
	var dirSize int64

	sizes := make(chan int64)
	readSize := func(path string, file os.FileInfo, err error) error { //nolint:unparam
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

func (conf *client) configureDiskUsage() {
	scheduleDiskUsage := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := scheduleDiskUsage.AddFunc(conf.diskCron, conf.getAndUpdateDiskSize)
	if err != nil {
		level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
	}
	scheduleDiskUsage.Start()
}

func (conf *client) getAndUpdateDiskSize() {
	conf.lock.Lock()
	level.Info(conf.logger).Log(common.LogCategoryMsg, getCronScheduledMessage("pipeline size")) //nolint:errcheck
	for _, path := range conf.paths {
		level.Debug(conf.logger).Log(common.LogCategoryMsg, fmt.Sprintf("pipeline present at %s would be scanned", path)) //nolint:errcheck
		size, pathType, err := conf.GetDiskSize(path)
		if err != nil {
			level.Error(conf.logger).Log(common.LogCategoryErr, err.Error()) //nolint:errcheck
		}
		CurrentPipelineSize[path] = PipelineSize{
			Size: size,
			Type: pathType,
		}
	}
	defer conf.lock.Unlock()
}
