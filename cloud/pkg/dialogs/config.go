package dialogs

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

// Config структура с настройками.
type Config struct {
	AutoPong bool  `yaml:"auto_pong"`
	Timeout  int64 `yaml:"timeout"`
	Debug    bool  `yaml:"debug"`
}

type CmdArgs struct {
	configPath string
}

func ParseArgs() CmdArgs {
	var path CmdArgs
	flag.StringVar(&path.configPath, "config", "configs/dev.yaml", "path to bot config file")
	flag.Parse()
	return path
}

func LoadConfig(path CmdArgs) Config {
	config := Config{}

	file, err := os.Open(path.configPath)
	if err != nil {
		panic("Incorrect config path")
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		panic(err)
	}

	return config
}
