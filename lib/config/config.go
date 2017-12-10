package config

import (
	"github.com/r3boot/rsoc/lib/logger"
)

var (
	CfgFile string
	log     *logger.Logger
)

func NewConfig(l *logger.Logger, fname string) (*MainConfig, error) {
	log = l
	config := &MainConfig{}
	CfgFile = fname

	err := config.CreateExample()
	if err != nil {
		return nil, err
	}

	err = config.Load()

	return config, err
}
