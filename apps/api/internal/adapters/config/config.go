package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseName, err := requiredEnv("DB_NAME")
	if err != nil {
		return Config{}, err
	}
	databaseHost, err := requiredEnv("DB_HOST")
	if err != nil {
		return Config{}, err
	}
	databaseUser, err := requiredEnv("DB_USER")
	if err != nil {
		return Config{}, err
	}
	databasePassword, err := requiredEnv("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}
	databaseSSLEnabled, err := requiredEnv("DB_SSL_ENABLED")
	if err != nil {
		return Config{}, err
	}
	databasePort, err := requiredEnv("DB_PORT")
	if err != nil {
		return Config{}, err
	}

	sslEnabled, err := strconv.ParseBool(databaseSSLEnabled)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_SSL_ENABLED: %w", err)
	}

	sslMode := "disable"
	if sslEnabled {
		sslMode = "require"
	}

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseUser, databasePassword),
		Host:   databaseHost + ":" + databasePort,
		Path:   databaseName,
	}
	query := databaseURL.Query()
	query.Set("sslmode", sslMode)
	databaseURL.RawQuery = query.Encode()

	return Config{Port: port, DatabaseURL: databaseURL.String()}, nil
}

func (c Config) Address() string {
	return ":" + c.Port
}

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}

	return value, nil
}
