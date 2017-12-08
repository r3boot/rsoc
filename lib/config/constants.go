package config

const EXAMPLE_CONFIG = `---
clusters:
  - name: webservers
    hosts:
      - 10.0.0.1
      - 10.0.0.2
      - 10.0.0.3
  - name: databases
    hosts:
      - 10.0.1.1
      - 10.0.1.2

scripts:
  - name: uname
    description: "Run uname -a"
    script: "uname -a"
  - name: temp_sshd
    description: "Spawns sshd on a high port number"
    script: |
        #!/bin/sh
        TMPFILE="$(mktemp)"
		sed -e 's,Port 22,Port 2222' /etc/ssh/sshd_config > ${TMPFILE}
        /usr/sbin/sshd -f ${TMPFILE}
  - name: df
    description "Run df -h"
    script: "df -h"`

type ClusterConfig struct {
	Name  string   `yaml:"name"`
	User  string   `yaml:"user"`
	Hosts []string `yaml:"hosts"`
}

type ScriptConfig struct {
	Name   string `yaml:"name"`
	Script string `yaml:"script"`
}

type MainConfig struct {
	Clusters []ClusterConfig `yaml:"clusters"`
	Scripts  []ScriptConfig  `yaml:"scripts"`
}
