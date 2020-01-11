package helpers

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/gobuffalo/packr"
	"github.com/neighborly/go-config"
	"github.com/neighborly/go-errors"
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
		log.Panic("unable to load config")
	}

	return c
}

func FindRootPath() string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Panic("unable to find current working path")
	}

	found, err := findRootWithGoMod(currentPath, 0)
	if err != nil {
		log.Info("Unable to find root of your project, using working directory instead.")
		return currentPath
	}
	return *found
}

func findRootWithGoMod(dir string, currentIteration int) (*string, error) {
	// Don't allow going deep more than 5 directories
	if currentIteration == 5 {
		return nil, errors.New("Unable to find the root of your project: go.mod not found")
	}

	var foundPath *string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to read directory %s", dir)

	}

	for _, f := range files {
		if f.Name() == "go.mod" {
			foundPath = &dir
			break
		}
	}

	if foundPath == nil {
		return findRootWithGoMod(path.Join(dir, ".."), currentIteration+1)
	}

	return foundPath, nil
}
