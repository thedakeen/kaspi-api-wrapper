package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort int `env:"HTTP_PORT"`
	KaspiAPI KaspiAPI
}

type KaspiAPI struct {
	Scheme       string `env:"KASPI_API_SCHEME" env-default:"basic"`
	BaseURLBasic string `env:"KASPI_API_BASE_URL_BASIC"`
	BaseURLStd   string `env:"KASPI_API_BASE_URL_STANDARD"`
	BaseURLEnh   string `env:"KASPI_API_BASE_URL_ENHANCED"`
	ApiKey       string `env:"KASPI_API_KEY"`
}

var (
	cfg *Config
)

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		panic("failed to load .env file" + err.Error())
	}

	cfg = &Config{}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		panic("failed to load environment variables: " + err.Error())
	}

	return cfg
}
