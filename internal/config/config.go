package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"dev"`

	HTTPPort int `env:"HTTP_PORT"`
	KaspiAPI KaspiAPI
	Database Database
}

type KaspiAPI struct {
	Scheme       string `env:"KASPI_API_SCHEME" env-default:"basic"`
	BaseURLBasic string `env:"KASPI_API_BASE_URL_BASIC"`
	BaseURLStd   string `env:"KASPI_API_BASE_URL_STANDARD"`
	BaseURLEnh   string `env:"KASPI_API_BASE_URL_ENHANCED"`
	ApiKey       string `env:"KASPI_API_KEY"`
}

type Database struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	Name     string `env:"DB_NAME" env-default:"kaspi_pay"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
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
