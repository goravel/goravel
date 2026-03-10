//go:build !windows

package fsutil

import (
	"path"

	"github.com/gookit/goutil/internal/comfunc"
)

// Realpath returns the shortest path name equivalent to path by purely lexical processing.
// Will expand ~ to home dir, and join with workdir if path is relative path.
func Realpath(pathStr string) string {
	pathStr = comfunc.ExpandHome(pathStr)

	if !IsAbsPath(pathStr) {
		pathStr = JoinSubPaths(comfunc.Workdir(), pathStr)
	}
	return path.Clean(pathStr)
}
