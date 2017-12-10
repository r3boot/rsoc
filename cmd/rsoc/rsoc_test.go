package main

import (
	"testing"
)

func TestConfigureApp(t *testing.T) {
	ConfigureApp()

	if Logger == nil {
		t.Errorf("init: Logger == nil")
	}

	if Config == nil {
		t.Errorf("init: config == nil")
	}

	if JobRunner == nil {
		t.Errorf("init: JobRunner == nil")
	}
}

func TestRunApplication(t *testing.T) {
	ConfigureApp()

	statusCode := RunApplication()
	if statusCode != 1 {
		t.Errorf("RunApplication no options: status code != 1: %d", statusCode)
	}

	*showClusters = true
	statusCode = RunApplication()
	if statusCode != 0 {
		t.Errorf("RunApplication -C: status code != 0: %d", statusCode)
	}
	*showClusters = false

	*showNodes = "test"
	statusCode = RunApplication()
	if statusCode != 0 {
		t.Errorf("RunApplication -N: status code != 0: %d", statusCode)
	}
	*showNodes = ""

	*showCommands = true
	statusCode = RunApplication()
	if statusCode != 0 {
		t.Errorf("RunApplication -L: status code != 0: %d", statusCode)
	}
	*showCommands = false

	*useCluster = "test"
	statusCode = RunApplication()
	if statusCode != 1 {
		t.Errorf("RunApplication -c: status code != 1: %d", statusCode)
	}
	*useCluster = ""
}
