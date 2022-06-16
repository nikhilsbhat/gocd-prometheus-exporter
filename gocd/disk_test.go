package gocd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_diskSize(t *testing.T) {
	t.Run("", func(t *testing.T) {
		actual := diskSize("/path/to/dir")
		assert.Equal(t, 0, actual/1024/1024/1024)
	})
}
