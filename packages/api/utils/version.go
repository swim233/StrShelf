package utils

import (
	"runtime"

	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
)

func DisplayVersion() {
	logger.Suger.Infof("List running info")
	logger.Suger.Infof("Version: %s", lib.Version)
	logger.Suger.Infof("Commit: %s", lib.GitCommit)
	logger.Suger.Infof("Last Commit: %s", lib.CommitMessage)
	logger.Suger.Infof("Build Time: %s", lib.BuildTime)
	logger.Suger.Infof("Go Version: %s", runtime.Version())
}
