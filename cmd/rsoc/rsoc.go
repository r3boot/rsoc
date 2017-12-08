package main

import (
	"flag"
	"os"

	"fmt"
	"os/user"
	"path/filepath"

	"github.com/r3boot/rsoc/lib/config"
	"github.com/r3boot/rsoc/lib/logger"
)

const (
	D_DEBUG     = false
	D_TIMESTAMP = false

	CFG_GLOBAL = "/etc/rsoc.yaml"
	CFG_USER   = "~/.config/rsoc/config.yaml"
)

var (
	useDebug     = flag.Bool("-d", D_DEBUG, "Enable debugging output")
	useTimestamp = flag.Bool("-T", D_TIMESTAMP, "Show timestamps in output")

	Config *config.MainConfig
	Logger *logger.Logger
)

func expandTilde(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("expandTilde user.Current: %v", err)
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func init() {
	var cfgFile string

	Logger = logger.NewLogger(*useTimestamp, *useDebug)

	if os.Getuid() == 0 {
		cfgFile = CFG_GLOBAL
	} else {
		cfgFile = CFG_USER
	}

	cfgFile, err := expandTilde(cfgFile)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	Config, err = config.NewConfig(Logger, cfgFile)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}
}

func main() {
	Logger.Debugf("%v", Config)
}
