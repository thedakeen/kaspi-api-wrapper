package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	HTTPPort int `env:"HTTP_PORT"`
	KaspiAPI KaspiAPI
}

type KaspiAPI struct {
	BaseURL string `env:"KASPI_API_BASE_URL"`
	ApiKey  string `env:"KASPI_API_KEY"`
}

var (
	cfg *Config
)

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found or failed to load")
	}

	cfg = &Config{}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		panic("failed to load environment variables: " + err.Error())
	}

	return cfg
}
