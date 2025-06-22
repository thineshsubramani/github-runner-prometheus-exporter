package platform

import (
	"runtime"
)

func GetOS() string {
	return runtime.GOOS // "linux", "windows", "darwin"
}
