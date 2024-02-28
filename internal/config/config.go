package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

var configInstance *Config

type Config struct {
	Env        string `yaml:"env" env-default:"dev"`
	HTTPServer `yaml:"http_server"`
	Storage    `yaml:"storage"`
	JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port string `yaml:"port" env-default:"8090"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" required:"true"`
	Password string `yaml:"password" required:"true"`
	DBname   string `yaml:"dbname" required:"true"`
	SSLmode  string `yaml:"sslmode" env-default:"require"`
}

type JWT struct {
	Secret string `yaml:"secret" required:"true"`
}

func mustLoad(path string) {
	if _, err := os.Stat(path); err != nil {
		log.Fatalf("error opening Config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatalf("error reading Config file: %s", err)
	}

	configInstance = &cfg
}

func Init(path string) {
	mustLoad(path)
}

func Cfg() Config {
	if configInstance == nil {
		log.Fatalf("config was not initialized")
	}
	return *configInstance
}
