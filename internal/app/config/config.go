package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

// Config of the running app
type Config struct {
	// Port on which the server should run (e.g. 8080)
	RunPort string `yaml:"run_port"`
	// Postgres settings
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	// Postgres database username (e.g. postgres)
	User string `yaml:"user"`
	// Postgres user password (e.g. 1234)
	Password string `yaml:"password"`
	// Host on which Postgres is running (e.g. localhost)
	Host string `yaml:"host"`
	// Port on which Postgres is running (e.g. 5432)
	Port string `yaml:"port"`
	// Database name (e.g. link cutter)
	DB string `yaml:"db"`
}

// Load data from the Yaml config to the specified path in the format "/path/to/file.yaml"
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
