package helpers

import (
	"os"
	"path"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/packr"
	"github.com/neighborly/go-config"
	"github.com/neighborly/go-pghelpers"
	"github.com/nrfta/go-log"
)

type appConfig struct {
	Meta             config.MetaConfig
	PostgresDatabase pghelpers.PostgresConfig `mapstructure:"postgres"`
}

func LoadConfig() appConfig {
	var c appConfig

	configPath := path.Join(FindRootPath(), "config")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Unable find config directory at: %s", configPath)
	}

	box := packr.NewBox(configPath)
	err := config.Load(box, &c)
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
		log.Info("Unable to find root of your project.", err)

		return currentPath
	}

	if current.Dir == "." {
		return currentPath
	}

	return current.Dir
}
