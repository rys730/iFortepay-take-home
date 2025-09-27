package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	App AppConfig `envPrefix:"APP_"`
	DB DBConfig `envPrefix:"DB_"`
}

type AppConfig struct {
	Port int `env:"PORT" envDefault:"8080"`
}

type DBConfig struct {
	Name            string `env:"NAME"`
	Driver          string `env:"DRIVER"`
	Host            string `env:"HOST"`
	Credentials     string `env:"CREDS"`
	Port            int    `env:"PORT"`
}

func (d *DBConfig) GetConnectionString() string {
	return fmt.Sprintf("%s://%s@%s:%d/%s?sslmode=disable", d.Driver, d.Credentials, d.Host, d.Port, d.Name)
}

func NewConfig() *Config {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		log.Info().Err(err).Msg("error loading .env file, will use os env")
	}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing env")
	}
	return &cfg
}