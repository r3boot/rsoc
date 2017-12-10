# Introduction
rsoc is a small utility, designed to run one-off commands on sets of
machines, grouped in 'clusters'. It accepts a command, and runs this
on a cluster. It does not accept any input.

# Features
* Use your system ssh as transport
* Job scheduler with a configurable amount of workers
* Output in greppable or JSON format

# Configuration
Configuration is done using a configuration located at
~/.config/rsoc/config.yaml if rsoc is running as a user, or
/etc/rsoc.yaml if rsoc is running as root. An example can be found
below:
~~~~
---
clusters:
  - name: webservers
    description: "Frontend facing webservers"
    hosts:
      - web01
      - web02
      - web03
  - name: databases
    description: "Dackend database servers"
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
~~~~

# Usage
### Available options
~~~~
alita [rsoc] >> ./build/rsoc -h
Usage of ./build/rsoc:
  -C    List available clusters
  -L    List available commands
  -N string
        List nodes part of cluster
  -T    Show timestamps in output
  -c string
        Cluster to connect to
  -d    Enable debugging output
  -o string
        Output modifier (default "grep")
  -w int
        Number of jobs to run in parallel (default 10)
~~~~

### Show available clusters
~~~~
alita [rsoc] >> ./build/rsoc -C
Name             Description                    Hosts
webservers       Frontend facing webservers     3
databases        Backend database servers       2
~~~~

### Show nodes for a cluster
~~~~
alita [rsoc] >> ./build/rsoc -N webservers
Nodes for webservers cluster:
web01
web02
web03
~~~~

### Show available commands
~~~~
alita [rsoc] >> ./build/rsoc -L
Name             Description
uname            Run uname -a
df               run df -h
uptime           Show the uptime of a system
~~~~

### Run command on cluster using greppable output
~~~~
[alita [rsoc] >> ./build/rsoc -c webservers uname
web01 stdout: Linux web01.example.com 4.13.12-1-ARCH #1 SMP PREEMPT Wed Nov 8 11:54:06 CET 2017 x86_64 GNU/Linux
web02 stdout: Linux web03.example.com 4.13.12-1-ARCH #1 SMP PREEMPT Wed Nov 8 11:54:06 CET 2017 x86_64 GNU/Linux
web03 stdout: Linux web03.example.com 4.13.12-1-ARCH #1 SMP PREEMPT Wed Nov 8 11:54:06 CET 2017 x86_64 GNU/Linux
~~~~

#### Run command on cluster using json output
~~~~
alita [rsoc] >> ./build/rsoc -o json -c wiekslag uptime | jq .
[
  {
    "node": "db01",
    "stdout": " 00:50:10 up 21 days,  6:34,  1 user,  load average: 0.85, 1.11, 1.42\n",
    "stderr": ""
  },
  {
    "node": "db02",
    "stdout": " 00:50:10 up 33 days, 22 min,  0 users,  load average: 0.51, 0.69, 0.66\\n",
    "stderr": ""
  },
]
~~~~