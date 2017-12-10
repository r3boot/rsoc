package config

import (
	"testing"

	"io/ioutil"

	"os"

	"fmt"

	"strings"

	"github.com/r3boot/rsoc/lib/logger"
)

func createNewConfig(t *testing.T) (*MainConfig, string) {
	tmpFile := CreateTempFile(t)

	log := logger.NewLogger(false, true)

	config, err := NewConfig(log, tmpFile)
	if err != nil {
		t.Errorf("NewConfig failed: %v", err)
	}

	return config, tmpFile
}

func TestMainConfig_HasCluster(t *testing.T) {
	config, tmpFile := createNewConfig(t)
	defer CleanupTempFile(t, tmpFile)

	if !config.HasCluster("webservers") {
		t.Errorf("config.HasCluster: config does not contain webservers cluster")
	}

	if config.HasCluster("nonexistent") {
		t.Errorf("config.HasCluster: config contains nonexistent cluster")
	}
}

func TestMainConfig_GetCluster(t *testing.T) {
	config, tmpFile := createNewConfig(t)
	defer CleanupTempFile(t, tmpFile)

	cluster, err := config.GetCluster("databases")
	if err != nil {
		t.Errorf("config.GetCluster: err != nil: %v", err)
	} else if cluster.Name == "" {
		t.Errorf("config.GetCluster: cluster.Name is empty")
	}

	cluster, err = config.GetCluster("nonexistent")
	if err == nil {
		t.Errorf("config.GetCluster: nonexistent cluster but err == nil")
	} else if cluster.Name != "" {
		t.Errorf("config.GetCluster: nonexistent cluster but cluster.Name is not empty")
	}
}

func TestMainConfig_GetAllClusters(t *testing.T) {
	config, tmpFile := createNewConfig(t)
	defer CleanupTempFile(t, tmpFile)

	allClusters := config.GetAllClusters()
	if len(allClusters) != 2 {
		t.Errorf("config.GetAllClusters: expected 2 clusters, got %d", len(allClusters))
	}
}

func TestMainConfig_HasCommand(t *testing.T) {
	config, tmpFile := createNewConfig(t)
	defer CleanupTempFile(t, tmpFile)

	if !config.HasCommand("uname") {
		t.Errorf("config.HasCommand: config does not contain uname command")
	}

	if config.HasCommand("nonexistent") {
		t.Errorf("config.HasCommand: config contains nonexistent command")
	}
}

func TestMainConfig_GetCommand(t *testing.T) {
	config, tmpFile := createNewConfig(t)
	defer CleanupTempFile(t, tmpFile)

	command, err := config.GetCommand("df")
	if err != nil {
		t.Errorf("config.GetCommand: err != nil: %v", err)
	} else if command.Name == "" {
		t.Errorf("config.GetCommand: command.Name is empty")
	}

	command, err = config.GetCommand("nonexistent")
	if err == nil {
		t.Errorf("config.GetCommand: nonexistent command but err == nil")
	} else if command.Name != "" {
		t.Errorf("config.GetCommand: nonexistent command but cluster.Name is not empty")
	}
}

func TestMainConfig_GetAllCommands(t *testing.T) {
	config, tmpFile := createNewConfig(t)
	defer CleanupTempFile(t, tmpFile)

	allCommands := config.GetAllCommands()
	if len(allCommands) != 2 {
		t.Errorf("config.GetAllCommands: expected 2 commands, got %d", len(allCommands))
	}
}

func TestMainConfig_CreateExample(t *testing.T) {
	log := logger.NewLogger(false, true)

	// Test existing file
	existingFile := CreateTempFile(t)
	err := ioutil.WriteFile(existingFile, []byte("TestMainConfig_CreateExample"), 0400)
	if err != nil {
		t.Errorf("ioutil.Writefile err != nil: %v", err)
	}

	if _, err = NewConfig(log, existingFile); err == nil {
		t.Errorf("NewConfig with existing but unparseable file err == nil")
	}

	if err = os.Remove(existingFile); err != nil {
		t.Errorf("os.Remove err != nil: %v", err)
	}

	// Test filepath.Abs (nonworking atm :( )
	curPwd := os.Getenv("PWD")
	os.Setenv("PWD", "/nonexistent/path")

	nonExistingPWD := "test"
	_, err = NewConfig(log, nonExistingPWD)
	if err != nil {
		fmt.Printf("err: %v", err)
		if !strings.Contains(err.Error(), "MainConfig.CreateExample filepath.Abs:") {
			t.Errorf("NewConfig nonExistingPWD, but err != filepath.Abs error")
		}
	}
	os.Setenv("PWD", curPwd)
	if err = os.Remove(nonExistingPWD); err != nil {
		t.Errorf("os.Remove nonExistingPWD err != nil: %v", err)
	}

	// Test os.MkdirAll
	unwritablePath := "/nonexisting/config.yaml"
	_, err = NewConfig(log, unwritablePath)
	if err != nil && !strings.Contains(err.Error(), "MainConfig.CreateExample os.MkdirAll:") {
		fmt.Printf("os.MkdirAll: %v\n", err.Error())
		t.Error("NewConfig unwritablePath, but err != os.MkdirAll error")
	}
}

func TestMainConfig_Load(t *testing.T) {
	log := logger.NewLogger(false, true)

	unreadableFile := CreateTempFile(t)
	err := ioutil.WriteFile(unreadableFile, []byte(EXAMPLE_CONFIG), 0000)
	if err != nil {
		t.Errorf("ioutil.WriteFile err != nil: %v", err)
	}

	_, err = NewConfig(log, unreadableFile)
	if err != nil && !strings.Contains(err.Error(), "config.Load ioutil.ReadFile:") {
		fmt.Printf("ioutil.Readfile: %v\n", err)
		t.Error("NewConfig unreadableFile, but err != ioutil.ReadFile error")
	}

	if err = os.Remove(unreadableFile); err != nil {
		t.Errorf("os.Remove err != nil: %v", err)
	}
}
