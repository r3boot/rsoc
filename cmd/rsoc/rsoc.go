package main

import (
	"flag"
	"os"

	"fmt"

	"github.com/r3boot/rsoc/lib/cliutils"
	"github.com/r3boot/rsoc/lib/config"
	"github.com/r3boot/rsoc/lib/jobrunner"
	"github.com/r3boot/rsoc/lib/logger"
)

const (
	D_DEBUG         = false
	D_TIMESTAMP     = false
	D_CLUSTER       = ""
	D_NUM_WORKERS   = 10
	D_SHOW_CLUSTERS = false
	D_SHOW_NODES    = ""
	D_SHOW_COMMANDS = false
	D_OUTPUT_MOD    = jobrunner.MOD_GREP

	CFG_GLOBAL = "/etc/rsoc.yaml"
	CFG_USER   = "~/.config/rsoc/config.yaml"
)

var (
	useDebug       = flag.Bool("d", D_DEBUG, "Enable debugging output")
	useTimestamp   = flag.Bool("T", D_TIMESTAMP, "Show timestamps in output")
	useCluster     = flag.String("c", D_CLUSTER, "Cluster to connect to")
	numWorkers     = flag.Int("w", D_NUM_WORKERS, "Number of jobs to run in parallel")
	showClusters   = flag.Bool("C", D_SHOW_CLUSTERS, "List available clusters")
	showNodes      = flag.String("N", D_SHOW_NODES, "List nodes part of cluster")
	showCommands   = flag.Bool("L", D_SHOW_COMMANDS, "List available commands")
	outputModifier = flag.String("o", D_OUTPUT_MOD, "Output modifier")

	Config    *config.MainConfig
	Logger    *logger.Logger
	JobRunner *jobrunner.JobRunner
)

func ConfigureApp() {
	var cfgFile string

	flag.Parse()

	Logger = logger.NewLogger(*useTimestamp, *useDebug)

	if os.Getuid() == 0 {
		cfgFile = CFG_GLOBAL
	} else {
		cfgFile = CFG_USER
	}

	cfgFile, err := cliutils.ExpandTilde(cfgFile)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	Config, err = config.NewConfig(Logger, cfgFile)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	cliutils.Setup(Logger)

	JobRunner = jobrunner.NewJobRunner(Logger, Config)
}

func RunApplication() int {
	didShowCommand := false

	if *showClusters {
		cliutils.ShowClusters(Config)
		didShowCommand = true
	}

	if *showNodes != "" {
		cliutils.ShowNodes(Config, *showNodes)
		didShowCommand = true
	}

	if *showCommands {
		cliutils.ShowCommands(Config)
		didShowCommand = true
	}

	if didShowCommand {
		return 0
	}

	if *useCluster == "" {
		fmt.Printf("ERROR: Need a cluster to connect to\n")
		return 1
	}

	if !Config.HasCluster(*useCluster) {
		fmt.Printf("ERROR: No such cluster\n")
		return 1
	}

	if len(flag.Args()) == 0 {
		fmt.Printf("ERROR: Need a command to run\n")
		return 1
	}

	for _, command := range flag.Args() {
		if !Config.HasCommand(command) {
			fmt.Printf("WARNING: No such command: %s", command)
			continue
		}

		cmd, err := Config.GetCommand(command)
		if err != nil {
			// TODO: implement onfailure behaviour
			fmt.Printf("WARNING: %v", err)
			continue
		}

		if cmd.Command != "" && cmd.Script != "" {
			fmt.Printf("WARNING: %s has both command and script", cmd.Name)
			continue
		} else if cmd.Command == "" && cmd.Script == "" {
			fmt.Printf("WARNING: %s has neither command or script set", cmd.Name)
		} else {
			job := jobrunner.Job{
				Cluster: *useCluster,
				Command: command,
			}
			cluster, err := Config.GetCluster(*useCluster)
			if err != nil {
				fmt.Printf("ERROR: %v", err)
				return 1
			}

			maxWorkers := *numWorkers
			numNodes := len(cluster.Hosts)

			if numNodes < maxWorkers {
				maxWorkers = numNodes
			}
			Logger.Debugf("main: NumWorkers set to %d", maxWorkers)

			JobRunner.StartWorkers(maxWorkers)
			JobRunner.Submit(job)
			JobRunner.SubmitKillJob()
			JobRunner.Start(*outputModifier)
		}
	}

	return 0
}

func init() {
	ConfigureApp()
}

func main() {
	os.Exit(RunApplication())
}
