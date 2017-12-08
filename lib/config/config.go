package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"os"

	"path/filepath"

	"strings"

	"github.com/r3boot/rsoc/lib/logger"
	yaml "gopkg.in/yaml.v2"
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

func (c *MainConfig) CreateExample() error {
	log.Infof("cfgFile: %s", cfgFile)
	_, err := os.Stat(cfgFile)
	if err == nil {
		return nil
	}

	fullPath, err := filepath.Abs(cfgFile)
	if err != nil {
		return fmt.Errorf("MainConfig.CreateExample filepath.Abs: %v", err)
	}

	dirPath := path.Dir(fullPath)
	log.Debugf("dirPath: %s", dirPath)

	tokens := strings.Split(dirPath, "/")
	mkdirPath := ""
	for _, item := range tokens {
		mkdirPath += fmt.Sprintf("/%s", item)
		log.Debugf("MainConfig.CreateExample: %s", mkdirPath)
	}

	return nil
}

func (c *MainConfig) Load() error {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("config.Load ioutil.ReadFile: %v", err)
	}

	if err = yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("config.Load yaml.Unmarshal: %v", err)
	}

	return nil
}
