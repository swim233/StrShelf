package utils

import (
	"fmt"
	"runtime"
)

var (
	Version       = "dev"
	GitCommit     = "unknown"
	CommitMessage = "unknown"
	BuildTime     = "unknown"
	DebugModeStr  = "true"
	DebugMode     = DebugModeStr == "true"
)

func GetVersion() string {
	return fmt.Sprintf(
		"Running info:\nVersion: %s\nCommit: %s\nLast Commit: %s\nBuild Time: %s\nGo Version: %s",
		Version,
		GitCommit,
		CommitMessage,
		BuildTime,
		runtime.Version(),
	)
}
