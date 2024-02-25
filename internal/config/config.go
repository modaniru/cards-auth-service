package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Port     string         `yaml:"port" env-default:"80"`
	Salt     string         `yaml:"salt" env-required:"true"`
	TokenTTL time.Duration  `yaml:"ttl" env-required:"true"`
	Env      string         `yaml:"env"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	DataSource string `yaml:"data_source" env-required:"true"`
}

func MustLoad() *Config {
	file, err := os.Open("./config.yaml")
	if err != nil {
		panic(err.Error())
	}
	cfg := Config{}

	err = cleanenv.ParseYAML(file, &cfg)
	if err != nil {
		panic(err.Error())
	}
	return &cfg
}
