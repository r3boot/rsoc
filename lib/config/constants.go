package config

const EXAMPLE_CONFIG = `---
clusters:
  - name: webservers
    description: "Frontend facing webservers"
    hosts:
      - web01
      - web02
      - web03
  - name: databases
    description: "backend database servers"
    hosts:
      - db01
      - db02

commands:
  - name: uname
    description: "Run uname -a"
    command: "uname -a"
  - name: df
    description: "Run df -h"
    command: "df -h"`

type ClusterConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	User        string   `yaml:"user"`
	Hosts       []string `yaml:"hosts"`
}

type CommandConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Script      string `yaml:"script"`
	Command     string `yaml:"command"`
}

type MainConfig struct {
	Clusters []ClusterConfig `yaml:"clusters"`
	Commands []CommandConfig `yaml:"commands"`
}
