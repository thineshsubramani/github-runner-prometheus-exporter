package platform

import (
	"runtime"
)

func GetOS() string {
	return runtime.GOOS // "linux", "windows", "darwin"
}

// We are still pulling info from config yaml like static infos
// TODO: more dynamic way to pull server metadata for labels
// EG. Server Pool, OS Version
//
