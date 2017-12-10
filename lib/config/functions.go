package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

func (c MainConfig) HasCluster(name string) bool {
	for _, entry := range c.Clusters {
		if entry.Name == name {
			return true
		}
	}
	return false
}

func (c MainConfig) GetCluster(name string) (ClusterConfig, error) {
	for _, entry := range c.Clusters {
		if entry.Name == name {
			return entry, nil
		}
	}

	return ClusterConfig{}, fmt.Errorf("MainConfig.GetCluster: No such cluster: %v", name)
}

func (c *MainConfig) GetAllClusters() []ClusterConfig {
	return c.Clusters
}

func (c *MainConfig) HasCommand(name string) bool {
	for _, entry := range c.Commands {
		if entry.Name == name {
			return true
		}
	}

	return false
}

func (c *MainConfig) GetCommand(name string) (CommandConfig, error) {
	for _, entry := range c.Commands {
		if entry.Name == name {
			return entry, nil
		}
	}

	return CommandConfig{}, fmt.Errorf("MainConfig.GetCommand: No such command %v", name)
}

func (c *MainConfig) GetAllCommands() []CommandConfig {
	return c.Commands
}

func (c *MainConfig) CreateExample() error {
	_, err := os.Stat(CfgFile)
	if err == nil {
		return nil
	}

	fullPath, err := filepath.Abs(CfgFile)
	if err != nil {
		return fmt.Errorf("MainConfig.CreateExample filepath.Abs: %v", err)
	}

	dirPath := path.Dir(fullPath)
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("MainConfig.CreateExample os.MkdirAll: %v", err)
	}

	err = ioutil.WriteFile(CfgFile, []byte(EXAMPLE_CONFIG), 0644)
	if err != nil {
		return fmt.Errorf("MainConfig.CreateExample ioutil.WriteFile: %v", err)
	}

	log.Debugf("MainConfig.CreateExample: example created in %s", CfgFile)

	return nil
}

func (c *MainConfig) Load() error {
	data, err := ioutil.ReadFile(CfgFile)
	if err != nil {
		return fmt.Errorf("config.Load ioutil.ReadFile: %v", err)
	}

	if err = yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("config.Load yaml.Unmarshal: %v", err)
	}

	return nil
}
