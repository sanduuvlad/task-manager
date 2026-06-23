package config

import "os"

type Config struct {
	DATABASEUrl string
	JWTSecret   string
}

func Load() Config {
	return Config{
		DATABASEUrl: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}
