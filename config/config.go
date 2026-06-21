package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RunPort  string         `yaml:"run_port"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Db       string `yaml:"db"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
