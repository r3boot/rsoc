package jobrunner

import (
	"os"

	"github.com/r3boot/rsoc/lib/config"
)

const (
	KILL_PILL = "DIEDIEDIE"
	MAX_JOBS  = 1024
	MOD_GREP  = "grep"
	MOD_JSON  = "json"

	TEST_CLUSTER = "test"
	TEST_STDOUT  = "stdout"
	TEST_STDERR  = "stderr"
	TEST_ERR     = "err"

	TEST_STDOUT_MSG = "This is some stdout output\nwhich is split out over multiple lines\n"
	TEST_STDERR_MSG = "This is some stderr output\nwhich is split out over multiple lines\n"
	TEST_ERR_MSG    = "This is a test error"
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
	TestFd      *os.File
}
