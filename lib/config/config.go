package config

import (
	"github.com/r3boot/rsoc/lib/logger"
)

var (
	cfgFile string
	log     *logger.Logger
)

func NewConfig(l *logger.Logger, fname string) (*MainConfig, error) {
	log = l
	config := &MainConfig{}
	cfgFile = fname

	err := config.CreateExample()
	if err != nil {
		return nil, err
	}

	err = config.Load()

	return config, err
}
