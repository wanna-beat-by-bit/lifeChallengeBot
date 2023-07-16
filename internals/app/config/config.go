package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
}

func MustLoad(path string) *Config {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file '%s' does not exist", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg
}
