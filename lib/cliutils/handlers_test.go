package cliutils

import (
	"io/ioutil"
	"os"
	"testing"

	"strings"

	"fmt"

	"github.com/r3boot/rsoc/lib/config"
	"github.com/r3boot/rsoc/lib/logger"
)

const (
	TMP_DIR                = "/tmp"
	TMP_PREFIX             = "rsoc_lib_config_test"
	TEST_CLUSTER_VALID_CFG = `---
clusters:
  - name: test
    description: "Test cluster for unit testing"
    hosts:
      - localhost

commands:
  - name: true
    description: "Run the true command"
    command: "true"`
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

func CreateTestConfig(t *testing.T, content string) (*logger.Logger, *config.MainConfig, string) {
	log := logger.NewLogger(false, true)

	tmpFile := CreateTempFile(t)

	err := ioutil.WriteFile(tmpFile, []byte(content), 0400)
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

func TestShowClusters(t *testing.T) {
	testLogger, Config, cfgFile := CreateTestConfig(t, TEST_CLUSTER_VALID_CFG)
	defer CleanupTestconfig(t, cfgFile)

	tmpFile := CreateTempFile(t)
	defer CleanupTempFile(t, tmpFile)

	testFd, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("ShowClusters os.Create err != nil: %v", err)
	}

	Setup(testLogger)
	TestFd = testFd

	ShowClusters(Config)

	TestFd.Close()

	data, err := ioutil.ReadFile(tmpFile)
	lines := strings.Split(string(data), "\n")
	if !strings.HasPrefix(lines[0], "Name") {
		t.Errorf("ShowClusters: first line does not start with Name")
	}
	if !strings.HasPrefix(lines[1], "test") {
		t.Errorf("ShowClusters: test cluster not found")
	}
}

func TestShowNodes(t *testing.T) {
	testLogger, Config, cfgFile := CreateTestConfig(t, TEST_CLUSTER_VALID_CFG)
	defer CleanupTestconfig(t, cfgFile)

	tmpFileValidCluster := CreateTempFile(t)
	defer CleanupTempFile(t, tmpFileValidCluster)

	testFd, err := os.Create(tmpFileValidCluster)
	if err != nil {
		t.Fatalf("ShowNodes os.Create err != nil: %v", err)
	}

	Setup(testLogger)
	TestFd = testFd

	ShowNodes(Config, "test")

	TestFd.Close()

	data, err := ioutil.ReadFile(tmpFileValidCluster)
	lines := strings.Split(string(data), "\n")
	if !strings.HasPrefix(lines[0], "Nodes for") {
		t.Errorf("ShowNodes: first line does not start with 'Nodes for'")
	}
	if !strings.HasPrefix(lines[1], "localhost") {
		t.Errorf("ShowNodes: localhost node not found in test cluster")
	}

	tmpFileInvalidCluster := CreateTempFile(t)
	defer CleanupTempFile(t, tmpFileInvalidCluster)
	testFd, err = os.Create(tmpFileInvalidCluster)
	if err != nil {
		t.Fatalf("ShowNodes os.Create err != nil: %v", err)
	}

	Setup(testLogger)
	TestFd = testFd

	ShowNodes(Config, "nonexistingcluster")

	TestFd.Close()

	data, err = ioutil.ReadFile(tmpFileInvalidCluster)
	lines = strings.Split(string(data), "\n")
	if !strings.HasPrefix(lines[0], "ShowNodes: MainConfig.GetCluster: No such cluster:") {
		fmt.Printf("line: %s\n", lines[0])
		t.Errorf("ShowNodes: specified invalid cluster, but got wrong error")
	}
}

func TestShowCommands(t *testing.T) {
	testLogger, Config, cfgFile := CreateTestConfig(t, TEST_CLUSTER_VALID_CFG)
	defer CleanupTestconfig(t, cfgFile)

	tmpFile := CreateTempFile(t)
	defer CleanupTempFile(t, tmpFile)

	testFd, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("ShowNodes os.Create err != nil: %v", err)
	}

	Setup(testLogger)
	TestFd = testFd

	ShowCommands(Config)

	TestFd.Close()

	data, err := ioutil.ReadFile(tmpFile)
	lines := strings.Split(string(data), "\n")

	if !strings.HasPrefix(lines[0], "Name") {
		t.Errorf("ShowCommands: expected 'Name'")
	}

	if !strings.HasPrefix(lines[1], "true") {
		t.Errorf("ShowCommands: expected 'true' command")
	}
}
