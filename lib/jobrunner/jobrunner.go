package jobrunner

import (
	"github.com/r3boot/rsoc/lib/config"
	"github.com/r3boot/rsoc/lib/logger"
)

var log *logger.Logger

func NewJobRunner(l *logger.Logger, cfg *config.MainConfig) *JobRunner {
	log = l
	runner := &JobRunner{
		config:      cfg,
		JobQueue:    make(chan NodeJob, MAX_JOBS),
		ResultQueue: make(chan Result, MAX_JOBS),
	}

	return runner
}
