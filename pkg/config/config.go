package config

import (
	"context"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

var (
	k            = koanf.New(".")
	loadedConfig *Config
)

type QueueConfig struct {
	Name  string `koanf:"name"`
	DB    string `koanf:"db"`
	Retry int    `koanf:"retry"`
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

func LoadConfigFile(ctx context.Context, path string) error {

	log.Ctx(ctx).Debug().Msgf("Loading config file: %s", path)
	file := file.Provider(path)
	err := k.Load(file, yaml.Parser())
	if err != nil {
		return err
	}

	loadedConfig = &Config{}
	if err := k.Unmarshal("", loadedConfig); err != nil {
		return err
	}

	return nil
}
