package utils

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func GetVersion() string {
	return fmt.Sprintf(
		"Running info:\nVersion: %s\nCommit: %s\nBuild Time: %s\nGo Version: %s",
		Version,
		GitCommit,
		BuildTime,
		runtime.Version(),
	)
}
