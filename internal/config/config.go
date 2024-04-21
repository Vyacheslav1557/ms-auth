package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"time"
)

var configInstance *Config

type Config struct {
	Env       string        `env:"ENV" env-required:"true"`
	Port      string        `env:"PORT" env-required:"true"`
	DSN       string        `env:"DSN" env-required:"true"`
	JWTSecret string        `env:"JWT_SECRET" env-required:"true"`
	JWTMaxAge time.Duration `env:"JWT_MAX_AGE" env-required:"true"`
}

func SetupCfg() {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		slog.Error("error reading Config: %s", err)
	}
	configInstance = &cfg
}

func Cfg() Config {
	if configInstance == nil {
		slog.Error("config was not initialized")
	}
	return *configInstance
}
