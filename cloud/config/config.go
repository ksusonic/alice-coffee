package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Timeout  int64  `yaml:"timeout"`
	Debug    bool   `yaml:"debug"`
	Address  string `yaml:"-"`
	QueueUrl string `yaml:"queue_url"`

	SqsAccessKey string `env:"AWS_ACCESS_KEY_ID,notEmpty"`
	SqsSecretKey string `env:"AWS_SECRET_ACCESS_KEY,notEmpty"`
}

func LoadConfig() (*Config, error) {
	var config Config

	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&config.Address, "address", ":8080", "server address like 127.0.0.1:8080")
	flag.Int64Var(&config.Timeout, "timeout", 3000, "timeout in ms")
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

	return &config, nil
}
