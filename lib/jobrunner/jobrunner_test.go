package jobrunner

import (
	"io/ioutil"
	"testing"

	"os"

	"github.com/r3boot/rsoc/lib/config"
	"github.com/r3boot/rsoc/lib/logger"
)

const (
	TMP_DIR     = "/tmp"
	TMP_PREFIX  = "rsoc_lib_config_test"
	TEST_CONFIG = `---
clusters:
  - name: test
    description: "Test cluster for unit testing"
    hosts:
      - localhost
  - name: local
    description: "Test Cluster for local ssh access"
    hosts:
      - localhost

commands:
  - name: true
    description: "Run the true command"
    command: "true"
  - name: uname
    description: "run uname -s on a node"
    command: "uname -s"
  - name: sshtrigger
    description: "force nonexisting parameters into ssh"
    command: "--nonexisting --commandline --options"`
)

func CreateTempFile(t *testing.T) string {
	tmpFile, err := ioutil.TempFile(TMP_DIR, TMP_PREFIX)
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
	}
	fname := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		t.Errorf("Failed to close tempfile: %v", err)
	}
	if err := os.Remove(fname); err != nil {
		t.Errorf("Failed to remove tempfile: %v", err)
	}

	return fname
}

func CleanupTempFile(t *testing.T, fname string) {
	if err := os.Remove(fname); err != nil {
		t.Errorf("Failed to remove tempfile: %v", err)
	}
}

func CleanupTestconfig(t *testing.T, fname string) {
	if err := os.Remove(fname); err != nil {
		t.Errorf("Failed to remove tempfile: %v", err)
	}
}

func CreateTestConfig(t *testing.T) (*logger.Logger, *config.MainConfig, string) {
	log := logger.NewLogger(false, true)

	tmpFile := CreateTempFile(t)

	err := ioutil.WriteFile(tmpFile, []byte(TEST_CONFIG), 0400)
	if err != nil {
		t.Fatalf("config.NewConfig ioutil.WriteFile err != nil: %v", err)
	}

	config, err := config.NewConfig(log, tmpFile)
	if err != nil {
		t.Fatalf("config.NewConfig err != nil: %v", err)
	}

	if config == nil {
		t.Fatalf("config.NewConfig config == nil")
	}

	return log, config, tmpFile
}

func TestNewJobRunner(t *testing.T) {
	log, config, tmpFile := CreateTestConfig(t)
	defer CleanupTestconfig(t, tmpFile)

	jobRunner := NewJobRunner(log, config)
	if jobRunner == nil {
		t.Errorf("NewJobRunner jobRunner == nil")
	}
}
