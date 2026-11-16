package sso_grpc

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr         string `envconfig:"ADDR" required:"true"`
	RetriesCount int    `envconfig:"RETRIES" default:"5"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("GRPC_SSO", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get grpc server config: %w", err)
		panic(err)
	}
	return config
}
