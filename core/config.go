package core

import (
	"log"

	"go.uber.org/config"

	"go-clean-archi/servers/api"
	"go-clean-archi/services/zap"
)

type Config struct {
	API    api.Config `yaml:"api"`
	Logger zap.Config `yaml:"logger"`
}

func NewConfigFromYAML(paths ...string) Config {
	provider, err := config.NewYAML(getPathsAsSource(paths)...)

	if err != nil {
		log.Fatal("Could not read configs", err)
	}

	var conf Config

	if err := provider.Get("").Populate(&conf); err != nil {
		log.Fatal("Could not read configs", err)
	}

	return conf
}

func getPathsAsSource(paths []string) (sources []config.YAMLOption) {
	for _, path := range paths {
		sources = append(sources, config.File(path))
	}

	return
}
