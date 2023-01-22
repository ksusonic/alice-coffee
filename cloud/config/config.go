package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	AutoPong bool   `yaml:"auto_pong"`
	Timeout  int64  `yaml:"timeout"`
	Debug    bool   `yaml:"debug"`
	Address  string `yaml:"address"`
}

func LoadConfig() (*Config, error) {
	var config Config
	var configPath string

	flag.StringVar(&configPath, "config", "configs/dev.yaml", "path to config file")
	flag.StringVar(&config.Address, "address", ":8080", "server address like 127.0.0.1:8080")
	flag.Int64Var(&config.Timeout, "timeout", 3000, "timeout in ms")
	flag.BoolVar(&config.Debug, "debug", false, "debug mode")
	flag.BoolVar(&config.AutoPong, "auto_pong", true, "auto pong on /ping")

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

	return &config, nil
}
