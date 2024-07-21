package helpers

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/gobuffalo/here"
	"github.com/neighborly/go-pghelpers"
	"github.com/nrfta/go-config/v3"
	"github.com/nrfta/go-log"
)

type appConfig struct {
	Meta             config.MetaConfig
	PostgresDatabase pghelpers.PostgresConfig `mapstructure:"postgres"`
}

func LoadConfig() appConfig {
	var c appConfig

	configPath := path.Join(FindRootPath(), "internal", "config")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		internalConfigPath := configPath
		configPath = path.Join(FindRootPath(), "config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("Unable find config directory at: %s or %s", internalConfigPath, configPath)
		}
	}

	err := config.Load(newConfigFolder(configPath), &c)
	if err != nil {
		log.Panic("Unable to load config", err)
	}

	return c
}

func FindRootPath() string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Panic("Unable to find current working path")
	}

	current, err := here.Current()

	if err != nil {
		log.Infof("Unable to find root of your project. Using current working directory (%s) instead.", currentPath)

		return currentPath
	}

	if current.Dir == "." {
		return currentPath
	}

	return current.Dir
}

func newConfigFolder(configPath string) configFolder {
	configFile := path.Join(configPath, "config.json")
	configTestFile := path.Join(configPath, "/config_test.json")

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		configData = nil
	}

	configTestData, err := ioutil.ReadFile(configTestFile)
	if err != nil {
		configTestData = nil
	}

	return configFolder{
		config:      configData,
		config_test: configTestData,
	}
}

type configFolder struct {
	config      []byte
	config_test []byte
}

func (c configFolder) ReadFile(name string) ([]byte, error) {
	switch name {
	case "config.json":
		if c.config == nil {
			return nil, os.ErrNotExist
		}
		return c.config, nil
	case "config_test.json":
		if c.config_test == nil {
			return nil, os.ErrNotExist
		}
		return c.config_test, nil
	}
	return nil, os.ErrNotExist
}
