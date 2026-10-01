package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TransportMode string `env:"transport_mode" env-default:"http"`
	Server        Server `yaml:"server"`
	PgDb          PgDb   `yaml:"pg_db"`
	Pattern       string `yaml:"short_url_pattern"`
}

type Server struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"username" env-required:"true"`
}

type PgDb struct {
	Name     string `yaml:"name" env-default:"postgres"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"3306"`
	User     string `yaml:"username"`
	Password string `env:"pg_pass"`
}

func MustLoad() *Config {
	path := os.Getenv("config_path")
	if path == "" {
		path = "config.yaml"
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
