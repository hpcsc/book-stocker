package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Port        string `env:"PORT" envDefault:"3333"`
	AWSEndpoint string `env:"AWS_ENDPOINT"`
	AWSRegion   string `env:"AWS_REGION"`
	Environment string `env:"ENVIRONMENT"`
}

func Load() (*Configuration, error) {
	var c Configuration
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %v", err)
	}
	return &c, nil
}
