package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	GCSBucketName          string
	GCSEndpoint            string
	GCSPublicBaseURL       string
	GCSSignedURLAccessID   string
	GCSSignedURLPrivateKey []byte
	GCSSignedURLHostname   string
	GCSSignedURLInsecure   bool
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

	gcsBucketName := os.Getenv("GCS_BUCKET_NAME")
	if gcsBucketName == "" {
		gcsBucketName = "tsr-pg-journey-images"
	}

	gcsEndpoint := os.Getenv("GCS_ENDPOINT")
	if gcsEndpoint == "" {
		gcsEndpoint = "http://localhost:4443"
	}

	gcsPublicBaseURL := os.Getenv("GCS_PUBLIC_BASE_URL")
	if gcsPublicBaseURL == "" {
		gcsPublicBaseURL = strings.TrimRight(gcsEndpoint, "/") + "/" + gcsBucketName
	}

	gcsSignedURLInsecure := false
	if rawValue := os.Getenv("GCS_SIGNED_URL_INSECURE"); rawValue != "" {
		parsedValue, err := strconv.ParseBool(rawValue)
		if err != nil {
			return Config{}, fmt.Errorf("invalid GCS_SIGNED_URL_INSECURE: %w", err)
		}
		gcsSignedURLInsecure = parsedValue
	}
	var gcsSignedURLPrivateKey []byte
	if rawValue := os.Getenv("GCS_SIGNED_URL_PRIVATE_KEY"); rawValue != "" {
		gcsSignedURLPrivateKey = []byte(strings.ReplaceAll(rawValue, `\n`, "\n"))
	}

	return Config{
		Port:                   port,
		DatabaseURL:            databaseURL.String(),
		GCSBucketName:          gcsBucketName,
		GCSEndpoint:            gcsEndpoint,
		GCSPublicBaseURL:       gcsPublicBaseURL,
		GCSSignedURLAccessID:   os.Getenv("GCS_SIGNED_URL_ACCESS_ID"),
		GCSSignedURLPrivateKey: gcsSignedURLPrivateKey,
		GCSSignedURLHostname:   os.Getenv("GCS_SIGNED_URL_HOSTNAME"),
		GCSSignedURLInsecure:   gcsSignedURLInsecure,
	}, nil
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
