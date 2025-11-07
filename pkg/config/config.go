package config

import (
	"context"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	k            = koanf.New(".")
	loadedConfig *Config
)

type QueueConfig struct {
	Name string `koanf:"name"`
	DB   string `koanf:"db"`
}

type Config struct {
	Queues []QueueConfig `koanf:"queues"`
}

func GetConfig() *Config {
	if loadedConfig == nil {
		return &Config{} // Return empty config if none loaded
	}
	return loadedConfig
}

func loadConfigFile(ctx context.Context, path string) error {

	file := file.Provider(path)

	err := k.Load(file, yaml.Parser())
	if err != nil {
		return err
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return err
	}

	loadedConfig = &cfg
	return nil
}
