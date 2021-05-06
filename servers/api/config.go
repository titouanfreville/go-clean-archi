package api

import (
	"fmt"

	"go-clean-archi/servers/cors"
)

type Config struct {
	Scheme  string      `yaml:"scheme"`
	Host    string      `yaml:"host"`
	Port    int         `yaml:"port"`
	NoCache bool        `yaml:"no_cache"`
	CORS    cors.Config `yaml:"CORS"`
}

func (config Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

func (config Config) GetURL() string {
	return fmt.Sprintf("%s://%s", config.Scheme, config.GetAddress())
}
