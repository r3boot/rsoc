package cliutils

import (
	"fmt"

	"os"

	"github.com/r3boot/rsoc/lib/config"
)

var TestFd *os.File

func ShowClusters(Config *config.MainConfig) {
	outFd := os.Stdout
	if TestFd != nil {
		outFd = TestFd
	}

	fmt.Fprintf(outFd, "%-16s %-30s %-5s\n", "Name", "Description", "Hosts")
	for _, entry := range Config.GetAllClusters() {
		fmt.Fprintf(outFd, "%-16s %-30s %-5d\n", entry.Name, entry.Description, len(entry.Hosts))
	}
}

func ShowNodes(Config *config.MainConfig, name string) {
	outFd := os.Stdout
	if TestFd != nil {
		outFd = TestFd
	}

	cluster, err := Config.GetCluster(name)
	if err != nil {
		fmt.Fprintf(outFd, "ShowNodes: %v", err)
		return
	}

	fmt.Fprintf(outFd, "Nodes for %s cluster:\n", name)
	for _, node := range cluster.Hosts {
		fmt.Fprintf(outFd, "%s\n", node)
	}
}

func ShowCommands(Config *config.MainConfig) {
	outFd := os.Stdout
	if TestFd != nil {
		outFd = TestFd
	}

	fmt.Fprintf(outFd, "%-16s %-30s\n", "Name", "Description")
	for _, entry := range Config.GetAllCommands() {
		fmt.Fprintf(outFd, "%-16s %-30s\n", entry.Name, entry.Description)
	}
}
