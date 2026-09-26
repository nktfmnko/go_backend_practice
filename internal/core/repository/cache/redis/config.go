package core_redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr        string        `envconfig:"ADDR" required:"true"`
	Password    string        `envconfig:"PASSWORD" required:"true"`
	User        string        `envconfig:"USER"`
	DB          int           `envconfig:"DB"`
	MaxRetries  int           `envconfig:"MAXRETRIES"`
	DialTimeout time.Duration `envconfig:"DIALTIMEOUT"`
	Timeout     time.Duration `envconfig:"TIMEOUT"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get redis config: %w", err)
		panic(err)
	}
	return config
}
