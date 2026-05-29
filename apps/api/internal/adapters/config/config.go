package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://tsr_pg:tsr_pg@localhost:5432/tsr_pg?sslmode=disable"
	}

	return Config{Port: port, DatabaseURL: databaseURL}
}

func (c Config) Address() string {
	return ":" + c.Port
}
