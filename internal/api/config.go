package api

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	Platform  string
	JWTSecret string
	PolkaKey  string
}

func NewConfig() *Config {
	godotenv.Load()

	return &Config{
		DBURL:     os.Getenv("DB_URL"),
		Platform:  os.Getenv("PLATFORM"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		PolkaKey:  os.Getenv("POLKA_KEY"),
	}
}
