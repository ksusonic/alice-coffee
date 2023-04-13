package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug   bool          `yaml:"debug"`
	Server  ServerConfig  `json:"server"`
	Dialogs DialogsConfig `yaml:"dialogs"`
}

type ServerConfig struct {
	Debug   bool
	Address string `yaml:"address"`
}

type DialogsConfig struct {
	Timeout int64 `yaml:"timeout"`
}

func LoadConfig() (*Config, error) {
	var config Config

	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&config.Server.Address, "address", "127.0.0.1:8080", "server address like 127.0.0.1:8080")
	flag.Int64Var(&config.Dialogs.Timeout, "timeout", 3000, "timeout in ms")
	flag.BoolVar(&config.Debug, "debug", false, "debug mode")
	flag.Parse()

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("incorrect error path: %v", err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		d := yaml.NewDecoder(file)

		if err := d.Decode(&config); err != nil {
			return nil, err
		}
	}

	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	config.Server.Debug = config.Debug

	return &config, nil
}
