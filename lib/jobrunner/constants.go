package jobrunner

import "github.com/r3boot/rsoc/lib/config"

const (
	KILL_PILL = "DIEDIEDIE"
	MAX_JOBS  = 1024
	MOD_GREP  = "grep"
	MOD_JSON  = "json"
)

type NodeJob struct {
	Node    string
	Command string
}

type Job struct {
	Cluster string
	Command string
}

type Result struct {
	Node   string `json:"node"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	err    error
}

type JobRunner struct {
	config      *config.MainConfig
	numWorkers  int
	JobQueue    chan NodeJob
	ResultQueue chan Result
}
