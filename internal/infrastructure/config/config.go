package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	App      App      `yaml:"app"`
	Log      Log      `yaml:"log"`
	Database Database `yaml:"database"`

	// ...
}

func Load(configPath string) (*Config, error) {
	var cfg Config

	return &cfg, cleanenv.ReadConfig(configPath, &cfg)
}

type App struct {
	Listen int    `yaml:"listen"`
	Auth   string `yaml:"auth"`
}

type Log struct {
	Level zapcore.Level
	File  string
}

type Database struct {
	Addr     string
	Username string
	Password string
	Database string
	// ...
}
