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

	DBURL := os.Getenv("DB_URL")
	if DBURL == "" {
		panic("DB_URL environment variable is not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		panic("PLATFORM environment variable is not set")
	}

	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		panic("JWT_SECRET environment variable is not set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		panic("POLKA_KEY environment variable is not set")
	}

	return &Config{
		DBURL:     DBURL,
		Platform:  platform,
		JWTSecret: JWTSecret,
		PolkaKey:  polkaKey,
	}
}
